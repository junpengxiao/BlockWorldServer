package bwModel

import (
	"github.com/junpengxiao/BlockWorldServer/bwStruct"
	"os/exec"
)

//function signature. used for different models. To change model
//just change the processor in init()
type modelProcess func(bwStruct.BWData) bwStruct.BWData

var processor modelProcess

func init() {
	//processor = ModelSampleProcessor
	processor = ModelAProcessor
	//TODO : launch Julia server to load model into memory
	//But for now lanch it manully
	//go exec.Command("julia", "bwModel/test/testserver.jl")
	//TODO: print julia server feedback

}

func ProcessData(input bwStruct.BWData) bwStruct.BWData {
	return processor(input)
}
