package bwModel

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/junpengxiao/BlockWorldServer/bwStruct"
	"github.com/srom/tokenizer"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	//varibles used in yonatan's python codes
	dictionary map[string]int
	tok        tokenizer.Tokenizer //TODO, WARNING: This may be not thread safe
	vectorLen  = 79                //the input length. If the sentence is not long enough, use 1 to fill out blanks
	unkown     = "<unk>"
	//variables used in model A
	delta                  = 0.5 / 3
	threshold              = 0.01 //if loc is below threshold, then treat it as 0
	ErrPredictionError     = errors.New("Prediction didn't return a tuple that contains 3 elements")
	ErrBlockIndexNotFound  = errors.New("Prediction returned index is not listed in world description")
	ErrPredictionCollision = errors.New("Prediction collide with other block")
	port                   = 8081
) //TODO use .config to configure vectorLen etc, not wire coded

func init() {
	dictionary = make(map[string]int)
	dictionary[unkown] = 1
	file, err := os.Open("bwModel/VocabA.txt")
	if err != nil {
		log.Fatal("In ModelA Init, ", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		if !scanner.Scan() {
			break
		}
		freq, err := strconv.Atoi(scanner.Text())
		if err != nil {
			break
		}
		if freq >= 5 {
			dictionary[word] = len(dictionary) + 1
		}
	}
	tok = tokenizer.NewTreebankWordTokenizer()
}

func FindBlockByID(world []bwStruct.BWBlock, id int) (int, error) {
	for i, v := range world {
		if v.Id == id {
			return i, nil
		}
	}
	return -1, ErrBlockIndexNotFound
}

func ModelABuildResult(input bwStruct.BWData, result string) (ret bwStruct.BWData) {
	results := strings.Split(result, " ")
	if len(results) != 3 {
		ret.Error = ErrPredictionError.Error()
		return ret
	}
	results[2] = results[2][0 : len(results[2])-1] //removr /n which will cause Atoi error
	nums := make([]int, len(results))
	for i, v := range results {
		nums[i], _ = strconv.Atoi(v)
		log.Println("results ", v, " nums: ", nums)
	}
	source, err := FindBlockByID(input.World, nums[0])
	if err != nil {
		ret.Error = ErrBlockIndexNotFound.Error()
		return ret
	}
	target, err := FindBlockByID(input.World, nums[1])
	if err != nil {
		ret.Error = ErrBlockIndexNotFound.Error()
		return ret
	}
	predict := input.World[target]
	predict.Id = nums[0]
	switch nums[2] {
	case 1: //SW --
		predict.Loc[0] -= delta
		predict.Loc[2] -= delta
	case 2: //W  -*
		predict.Loc[0] -= delta
	case 3: //NW -+
		predict.Loc[0] -= delta
		predict.Loc[2] += delta
	case 4: //N  *+
		predict.Loc[2] += delta
	case 5: //NE ++
		predict.Loc[0] += delta
		predict.Loc[2] += delta
	case 6: //E  +*
		predict.Loc[0] += delta
	case 7: //SE +-
		predict.Loc[0] += delta
		predict.Loc[2] -= delta
	case 8: //S*-
		predict.Loc[2] -= delta
	}
	if predict.Loc[0]<threshold {
		predict.Loc[0] = 0
	}
	if predict.Loc[2]<threshold {
		predict.Loc[2] = 0
	}
	if !NoCollision(predict.Loc[0], predict.Loc[1], predict.Loc[2], input.World) {
		ret.Error = ErrPredictionCollision.Error()
	}
	ret.World = append(ret.World, predict)
	ret.Version = input.Version
	return ret
}

func debug(str string, object interface{}) {
	log.Println(str, ' ', object)
}

func ModelAProcessor(input bwStruct.BWData) bwStruct.BWData {
	tokens := tok.Tokenize(strings.ToLower(input.Input)) //WARNING this may be not thread safe
	debug("Tokens ", tokens)
	vector := make([]int, 0, vectorLen)
	for _, tk := range tokens {
		if i, ok := dictionary[tk]; ok {
			vector = append(vector, i)
		} else {
			vector = append(vector, dictionary[unkown])
		}
	}
	for len(vector) < vectorLen {
		vector = append(vector, dictionary[unkown])
	}
	debug("vector", vector)
	// put vector into model and grab the result
	conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Println(err)
		input.Error = err.Error()
		return input
	}
	str := strconv.Itoa(vector[0])
	for i := 1; i != len(vector); i++ {
		str = str + " " + strconv.Itoa(vector[i])
	}
	debug("send str: ", str)
	fmt.Fprintln(conn, str)
	result, err := bufio.NewReader(conn).ReadString('\n')
	debug("result", result)
	if err != nil {
		log.Println(err)
		input.Error = err.Error()
		return input
	}
	conn.Close()
	return ModelABuildResult(input, result)
}
