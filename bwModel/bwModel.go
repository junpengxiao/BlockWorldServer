package bwModel

import (
	"bwStruct"
)

//function signature. used for different models. To change model
//just change the processor in init()
type modelProcess func(bwStruct.BWData) bwStruct.BWData

var processor modelProcess

func init() {
	processor = ModelSampleProcessor
}

func ProcessData(input bwStruct.BWData) bwStruct.BWData {
	return processor(input)
}
