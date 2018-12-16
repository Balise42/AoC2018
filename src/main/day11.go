package main

import (
	"fmt"
	"math"
)

const gridSize = 300


func main() {
	const gridId = 1309

	energy := make([][]int64, gridSize)
	for i := 0; i < gridSize; i++ {
		energy[i] = make([]int64, gridSize)
	}

	for i := int64(0); i < gridSize; i++ {
		for j := int64(0); j < gridSize; j++ {
			rackId := i + 1 + 10
			cellEnergy := rackId * (j + 1)
			cellEnergy = cellEnergy + gridId
			cellEnergy = cellEnergy * rackId
			cellEnergy = (cellEnergy / 100) % 10
			cellEnergy = cellEnergy - 5
			energy[i][j] = cellEnergy
		}
	}

	fmt.Println(getMaxEnergyA(energy))
	fmt.Println(getMaxEnergyB(energy))
}

func getMaxEnergyA(energy [][]int64) (int, int) {
	maxEnergy := int64(math.MinInt64)
	xMax := -1
	yMax := -1

	for x := 0; x < gridSize-3; x++ {
		for y := 0; y < gridSize-3; y++ {
			squareEnergy := int64(0)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					squareEnergy = squareEnergy + energy[x+i][y+j]
				}
			}
			if squareEnergy > maxEnergy {
				maxEnergy = squareEnergy
				xMax = x + 1
				yMax = y + 1
			}
		}
	}
	return xMax, yMax
}

func getMaxEnergyB(energy [][]int64) (int, int, int) {
	maxEnergy := int64(math.MinInt64)
	xMax := -1
	yMax := -1
	sizeMax := -1

	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			squareEnergy := energy[x][y]
			if squareEnergy > maxEnergy {
				maxEnergy = squareEnergy
				xMax = x + 1
				yMax = y + 1
				sizeMax = 1
			}
			maxCoord := x
			if y > x {
				maxCoord = y
			}
			for size := 1; size < gridSize-maxCoord; size++ {

				for offset := 0; offset < size; offset++ {
					squareEnergy = squareEnergy + energy[x+size][y+offset] + energy[x+offset][y+size]
				}
				squareEnergy = squareEnergy + energy[x+size][y+size]

				if squareEnergy > maxEnergy {
					maxEnergy = squareEnergy
					xMax = x + 1
					yMax = y + 1
					sizeMax = size + 1
				}
			}
		}
	}
	return xMax, yMax, sizeMax
}
