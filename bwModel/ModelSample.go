package bwModel

import (
	"bwStruct"
)

var (
	delta = 0.5 / 3
	dx    = []float64{-delta, 0, delta}
	dz    = []float64{-delta, 0, delta}
)

func check(x1, x2, d float64) bool {
	if (x1 < x2 && x1+d > x2-d) || (x1-d < x2+d && x1 > x2) {
		return false
	}
	return true
}

func NoCollision(loc [3]float64, x, y, z float64, world []bwStruct.BWBlock) bool {
	if len(loc) != 3 {
		return false
	}
	x += loc[0]
	//y += loc[1]
	z += loc[2]
	if x < delta/2-1 || x > 1-delta/2 {
		return false
	}
	if z < delta/2-1 || z > 1-delta/2 {
		return false
	}
	for _, block := range world {
		if !check(block.Loc[0], x, delta/2) && !check(block.Loc[2], z, delta/2) {
			return false
		}
	}
	return true
}

func ModelSampleProcessor(input bwStruct.BWData) (output bwStruct.BWData) {
	if len(input.World) == 0 {
		output = input
		return
	}
	output.Version = input.Version
	output.Error = "Null"
	block := input.World[0]
	for _, tx := range dx {
		for _, tz := range dz {
			if NoCollision(block.Loc, tx, 0, tz, input.World) {
				block.Loc[0] += tx
				block.Loc[2] += tz
				output.World = append(output.World, block)
				return
			}
		}
	}
	block.Loc = [3]float64{0, 0, 0}
	output.World = append(output.World, block)
	return
}
