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

	currLen := len(reactPolymer(polymer))
	fmt.Println(currLen)

	for b:= byte(65); b <= byte(90); b++ {
		candidate := len(reactPolymer(clearPolymer(polymer, b)))
		if candidate < currLen {
			currLen = candidate
		}
	}
	fmt.Println(currLen)

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
	molecule := make([]byte, len(polymer))
	copy(molecule, polymer)
	for i := 0; i<len(molecule); i++ {
		if molecule[i] == toClear || molecule[i] == toClear + 32 {
			molecule = append(molecule[:i], molecule[i+1:]...)
			i--
		}
	}
	return molecule
}