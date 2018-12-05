package main

import (
	"fmt"
	"github.com/golang-collections/collections/stack"
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
	molecule := new(stack.Stack)
	for _, b := range polymer {
		if molecule.Peek() == b + 32 || molecule.Peek() == b - 32 {
			molecule.Pop()
		} else {
			molecule.Push(b)
		}
	}

	result := make([]byte, molecule.Len())
	for i := 0; i<molecule.Len(); i++ {
		result[i] = molecule.Pop().(byte)
	}
	return result
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