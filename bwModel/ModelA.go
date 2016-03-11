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
	delta                    = 0.5 / 3
	ErrPredictionError       = errors.New("Prediction didn't return a tuple that contains 3 elements")
	ErrBlockIndexCorssBorder = errors.New("Prediction returned index is not listed in world description")
	ErrPredictionCollision   = errors.New("Prediction collide with other block")
) //TODO use .config to configure vectorLen etc, not wire coded

func init() {
	dictionary[unkown] = 1
	file, err := os.Open("Vocab.txt")
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
	tok := tokenizer.NewTreebankWordTokenizer()
}

func ModelABuildResult(input bwStruct.BWData, result string) (ret bwStruct.BWData) {
	results := strings.Split(result, " ")
	if len(results) != 3 {
		ret.Error = ErrPredictionError.Error()
		return
	}
	nums := make([]int, len(result))
	for i, v := range results {
		nums[i] = strconv.Atoi(v)
	}
	source, err := FindBlockByID(input, nums[0])
	if err != nil {
		ret.Error = ErrBlockIndexCorssBorder.Error()
		return
	}
	target, err := FindBlockByID(input, nums[1])
	if err != nil {
		ret.Error = ErrBlockIndexCorssBorder.Error()
		return
	}
	predict := input.World[target]
	predict.Id = source
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
		predict[0] += delta
		predict[2] += delta
	case 6: //E  +*
		predict[0] += delta
	case 7: //SE +-
		predict[0] += delta
		predict[2] -= delta
	case 8: //S*-
		predict[2] -= delta
	}
	if !NoCollision(predict[0], predict[1], predict[2], input.World) {
		ret.Error = ErrPredictionCollision.Error()
	}
	ret.World = append(ret.World, predict)
	ret.Version = input.Version
}

func ModelAProcessor(input bwStruct.BWData) bwStruct.BWData {
	tokens := tok.Tokenize(strings.ToLower(input.Input)) //WARNING this may be not thread safe
	vector := make([]int, 0, vectorLen)
	for _, tk := range tokens {
		if i, ok := dictionary[tk]; ok {
			vector = append(vector, i)
		} else {
			vector = append(vector, dictionary[unkown])
		}
	}
	for len(vector) != vectorLen {
		vector = append(vector, dictionary[unkown])
	}
	// put vector into model and grab the result
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		log.Println(err)
		input.Error = err
		return input
	}
	fmt.Fprintln(conn, str)
	result, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err)
		input.Error = err
		return input
	}
	conn.Close()
	return ModelABuildResult(input, result)
}
