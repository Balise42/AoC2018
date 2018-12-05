package main

import (
	"fmt"
	"io/ioutil"

	"os"
)

func main() {
	filename := os.Args[1]
	polymer, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	polymer = reactPolymerWithStack(polymer)

	currLen := len(polymer)
	fmt.Println(currLen)

	for b:= byte(65); b <= byte(90); b++ {
		candidate := len(reactPolymerWithStack(clearPolymer(polymer, b)))
		if candidate < currLen {
			currLen = candidate
		}
	}
	fmt.Println(currLen)

}

func reactPolymerWithStack(polymer []byte) []byte {
	molecule := make([]byte, 0)
	for _, b := range polymer {
		if len(molecule) > 0 && (molecule[len(molecule) - 1] == b + 32 || molecule[len(molecule) - 1] == b - 32){
			molecule = molecule[:len(molecule) - 1]
		} else {
			molecule = append(molecule, b)
		}
	}
	return molecule
}


func reactPolymer(polymer []byte) []byte {
	molecule := make([]byte, len(polymer))
	copy(molecule, polymer)
	for true {
		removed := false
		for i := 0; i<len(molecule) - 1; i++ {
			if molecule[i] == molecule[i+1] + 32 || molecule[i+1] == molecule[i] + 32 {
				molecule = append(molecule[:i], molecule[i+1:]...)
				molecule = append(molecule[:i], molecule[i+1:]...)
				removed = true
			}
		}
		if !removed {
			break
		}
	}
	return molecule
}

func clearPolymer(polymer []byte, toClear byte) []byte {
	molecule := make([]byte, 0)

	for i := 0; i<len(polymer); i++ {
		if polymer[i] != toClear && polymer[i] != toClear + 32 {
			molecule = append(molecule, polymer[i])
		}
	}
	return molecule
}