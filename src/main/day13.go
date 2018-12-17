package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Track struct {
	Type  rune
	Paths map[rune]*Track
	X     int
	Y     int
}

func NewTrack(r rune, x int, y int) Track {
	t := Track{}
	t.Type = r
	t.X = x
	t.Y = y
	t.Paths = make(map[rune]*Track)
	return t
}

func (t *Track) fillTrack(network [][]Track) {
	if t.Type == '-' {
		t.fillEast(network)
		t.fillWest(network)
	} else if t.Type == '|' {
		t.fillNorth(network)
		t.fillSouth(network)
	} else if t.Type == '+' {
		t.fillEast(network)
		t.fillWest(network)
		t.fillNorth(network)
		t.fillSouth(network)
	} else if t.Type == '\\' {
		if t.X > 1 && (network[t.X-1][t.Y].Type == '-' || network[t.X-1][t.Y].Type == '+') {
			t.fillWest(network)
			t.fillSouth(network)
		} else {
			t.fillEast(network)
			t.fillNorth(network)
		}
	} else if t.Type == '/' {
		if t.Y < len(network)-1 && (network[t.X][t.Y+1].Type == '|' || network[t.X][t.Y+1].Type == '+') {
			t.fillEast(network)
			t.fillSouth(network)
		} else {
			t.fillWest(network)
			t.fillNorth(network)
		}
	}
}

func (t *Track) fillWest(network [][]Track) {
	t.Paths['W'] = &network[t.X-1][t.Y]
}

func (t *Track) fillEast(network [][]Track) {
	t.Paths['E'] = &network[t.X+1][t.Y]
}

func (t *Track) fillNorth(network [][]Track) {
	t.Paths['N'] = &network[t.X][t.Y-1]
}

func (t *Track) fillSouth(network [][]Track) {
	t.Paths['S'] = &network[t.X][t.Y+1]
}

type Cart struct {
	Pos    *Track
	Dir    rune
	Option string
	Id     int
	Moved  bool
}

func NewCart(pos *Track, dir rune, id int) Cart {
	cart := Cart{}
	cart.Pos = pos
	cart.Option = "left"
	cart.Id = id
	if dir == '>' {
		cart.Dir = 'E'
	} else if dir == '<' {
		cart.Dir = 'W'
	} else if dir == '^' {
		cart.Dir = 'N'
	} else if dir == 'v' {
		cart.Dir = 'S'
	}
	cart.Moved = false
	return cart
}

func (c Cart) compare(other Cart) int {
	if c.Pos.X == other.Pos.X && c.Pos.Y == other.Pos.Y {
		return 0
	}
	if c.Pos.Y < other.Pos.Y || (c.Pos.Y == other.Pos.Y && c.Pos.X < other.Pos.X) {
		return -1
	}
	return 1
}

func (c Cart) Collides(carts []Cart) int {
	for j, cart := range carts {
		if cart.Id == c.Id {
			continue
		}
		if cart.Pos == c.Pos {
			return j
		}
	}
	return -1
}

var dirs = "NESW"

func (c *Cart) Move() {
	if c.Moved {
		return
	}
	if c.Pos.Type == '|' || c.Pos.Type == '-' {
		c.Pos = c.Pos.Paths[c.Dir]
	} else if c.Pos.Type == '/' {
		if c.Dir == 'N' {
			c.Pos = c.Pos.Paths['E']
			c.Dir = 'E'
		} else if c.Dir == 'E' {
			c.Pos = c.Pos.Paths['N']
			c.Dir = 'N'
		} else if c.Dir == 'W' {
			c.Pos = c.Pos.Paths['S']
			c.Dir = 'S'
		} else {
			c.Pos = c.Pos.Paths['W']
			c.Dir = 'W'
		}
	} else if c.Pos.Type == '\\' {
		if c.Dir == 'S' {
			c.Pos = c.Pos.Paths['E']
			c.Dir = 'E'
		} else if c.Dir == 'W' {
			c.Pos = c.Pos.Paths['N']
			c.Dir = 'N'
		} else if c.Dir == 'E' {
			c.Pos = c.Pos.Paths['S']
			c.Dir = 'S'
		} else {
			c.Pos = c.Pos.Paths['W']
			c.Dir = 'W'
		}
	} else if c.Pos.Type == '+' {
		if c.Option == "straight" {
			c.Option = "right"
			c.Pos = c.Pos.Paths[c.Dir]
		} else if c.Option == "left" {
			i := strings.IndexRune(dirs, c.Dir)
			newDirIndex := i - 1
			if newDirIndex < 0 {
				newDirIndex = newDirIndex + 4
			}
			c.Dir = (rune)(dirs[newDirIndex])
			c.Pos = c.Pos.Paths[c.Dir]
			c.Option = "straight"
		} else if c.Option == "right" {
			i := strings.IndexRune(dirs, c.Dir)
			newDirIndex := (i + 1) % 4
			c.Dir = (rune)(dirs[newDirIndex])
			c.Pos = c.Pos.Paths[c.Dir]
			c.Option = "left"
		}
	}
	c.Moved = true
}

func fillNetwork(network [][]Track) {
	for i, l := range network {
		for j, t := range l {
			t.fillTrack(network)
			network[i][j] = t
		}
	}
}

type byPos []Cart

func (s byPos) Len() int {
	return len(s)
}

func (s byPos) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byPos) Less(i, j int) bool {
	return s[i].compare(s[j]) < 0
}

func main() {
	filename := os.Args[1]

	f, _ := os.Open(filename)
	defer f.Close()

	network := make([][]Track, 250)
	for i := 0; i < 250; i++ {
		network[i] = make([]Track, 250)
	}
	carts := make([]Cart, 0)

	scanner := bufio.NewScanner(f)
	y := 0
	id := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x, r := range line {
			if r == '>' || r == '<' {
				id++
				track := NewTrack('-', x, y)
				carts = append(carts, NewCart(&track, r, id))
				network[x][y] = track
			} else if r == '^' || r == 'v' {
				id++
				track := NewTrack('|', x, y)
				carts = append(carts, NewCart(&track, r, id))
				network[x][y] = track
			} else {
				network[x][y] = NewTrack(r, x, y)
			}
		}
		y++
	}

	fillNetwork(network)

	for {
		sort.Sort(byPos(carts))
		collided := false
		for j, c := range carts {
			c.Move()
			carts[j] = c
			collision := c.Collides(carts)
			if collision != -1 {
				if collision > j {
					carts = append(carts[:collision], carts[collision+1:]...)
					carts = append(carts[:j], carts[j+1:]...)
				} else {
					carts = append(carts[:j], carts[j+1:]...)
					carts = append(carts[:collision], carts[collision+1:]...)
				}
				collided = true
				break
			}
		}
		if len(carts) == 1 {
			if !carts[0].Moved {
				carts[0].Move()
			}
			fmt.Println(carts[0].Pos.X, carts[0].Pos.Y)
			break
		}
		if !collided {
			for i, _ := range carts {
				carts[i].Moved = false
			}
		}
	}
}
