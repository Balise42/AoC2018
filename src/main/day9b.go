package main

import "fmt"

type Node struct {
	Value int
	Prev *Node
	Next *Node
}

func main() {
	numPlayers := 428
	numMarbles := 7206100

	scores := make([]int, numPlayers)

	player := 0
	maxScore := 0

	pos := &Node{Value : 0}
	pos.Next = pos
	pos.Prev = pos

	for i := 1; i<=numMarbles; i++ {
		if i % 23 != 0 {
			curr := pos
			pos = curr.Next
			newNode := Node {Value : i, Prev:pos, Next:pos.Next}
			pos.Next.Prev = &newNode
			pos.Next = &newNode
			pos = &newNode
		} else {
			scores[player] += i
			for j := 0; j<7; j++ {
				pos = pos.Prev
			}
			scores[player] += pos.Value
			pos.Next.Prev = pos.Prev
			pos.Prev.Next = pos.Next
			pos = pos.Next
		}

		if scores[player] > maxScore {
			maxScore = scores[player]
		}
		player = (player + 1)%numPlayers
	}
	fmt.Println(maxScore)
}