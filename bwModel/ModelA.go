package bwModel

import (
	"bufio"
	"github.com/junpengxiao/BlockWorldServer/bwStruct"
	"github.com/srom/tokenizer"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	dictionary map[string]int
	tok        tokenizer.Tokenizer //TODO, WARNING: This may be not thread safe
	vectorLen  = 79                //the input length. If the sentence is not long enough, use 1 to fill out blanks
	unkown     = "<unk>"
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
	tok = tokenizer.NewTreebankWordTokenizer()
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
	//predict(vector)
	return input
}
