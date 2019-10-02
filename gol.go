package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

type board [][]bool

func (b board) String() (result string) {
	var symbol string

	for i := 0; i < len(b); i++ {
		for j := 0; j < len(b[i]); j++ {

			if b[i][j] {
				symbol = "■"
			} else {
				symbol = "□"
			}

			result += fmt.Sprintf("%v ", symbol)
		}
		result += "\n"
	}
	result += "\n"

	return result
}

type coord struct {
	x int
	y int
}

func newBoard(cardinality int) board {
	b := make(board, cardinality)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < cardinality; i++ {
		b[i] = make([]bool, cardinality)

		for j := 0; j < cardinality; j++ {
			b[i][j] = rand.Intn(5) < 3
		}
	}

	return b
}

func up(a coord) coord {
	return coord{a.x, a.y - 1}
}

func down(a coord) coord {
	return coord{a.x, a.y + 1}
}

func left(a coord) coord {
	return coord{a.x - 1, a.y}
}

func right(a coord) coord {
	return coord{a.x + 1, a.y}
}

func neighbours(origin coord, rows int, cols int) []coord {
	all := []coord{
		up(left(origin)),
		up(origin),
		up(right(origin)),
		right(origin),
		right(down(origin)),
		down(origin),
		left(down(origin)),
		left(origin),
	}

	res := make([]coord, 0)

	for _, r := range all {
		if r.x >= 0 && r.x < rows && r.y >= 0 && r.y < cols {
			res = append(res, r)
		}
	}

	return res
}

func evolveCell(b board, c coord) {
	i, j := c.x, c.y
	aliveNeighbours := 0
	for _, n := range neighbours(coord{i, j}, len(b), len(b[i])) {
		if b[n.x][n.y] {
			aliveNeighbours++
		}
	}

	if b[i][j] {
		if aliveNeighbours < 2 || aliveNeighbours > 3 {
			b[i][j] = false
		} else {
			b[i][j] = true
		}
	} else {
		if aliveNeighbours == 3 {
			b[i][j] = true
		}
	}
}

func evolve(b board) {
	for i := 0; i < len(b); i++ {
		for j := 0; j < len(b[i]); j++ {
			evolveCell(b, coord{i, j})
		}
	}
}

func growUp(b board) board {
	var nb board

	rows := len(b)
	cols := len(b[0])
	anyAlive := false

	for j := range b[0] {
		anyAlive = anyAlive || b[0][j]
	}

	if anyAlive {
		nb = make(board, rows+1)
		nb[0] = make([]bool, cols)
		for j := 0; j < cols; j++ {
			nb[0][j] = false
		}
		for i := 0; i < rows; i++ {
			nb[i+1] = b[i]
		}
		return nb
	}

	return b
}

func growDown(b board) board {
	var nb board

	anyAlive := false
	rows := len(b)
	cols := len(b[rows-1])

	for j := range b[rows-1] {
		anyAlive = anyAlive || b[rows-1][j]
	}

	if anyAlive {
		if nb == nil {
			nb = make(board, rows+1)
		}
		for i := 0; i < rows; i++ {
			nb[i] = b[i]
		}
		nb[rows] = make([]bool, cols)
		for j := 0; j < cols; j++ {
			nb[rows][j] = false
		}
		return nb
	}
	return b
}

func growLeft(b board) board {
	var nb board

	anyAlive := false
	rows := len(b)

	for i := 0; i < rows; i++ {
		anyAlive = anyAlive || b[i][0]
	}

	if anyAlive {
		if nb == nil {
			nb = make(board, rows)
		}
		for k, v := range b {
			nb[k] = append([]bool{false}, v[:]...)
		}
		return nb
	}
	return b
}

func growRight(b board) board {
	var nb board
	var anyAlive bool

	rows := len(b)
	cols := len(b[0])

	for i := 0; i < rows; i++ {
		anyAlive = anyAlive || b[i][cols-1]
	}

	if anyAlive {
		if nb == nil {
			nb = make(board, rows)
		}
		for k, v := range b {
			nb[k] = append(v[:], false)
		}
		return nb
	}
	return b
}

func expand(b board) board {
	return growRight(growLeft(growDown(growUp(b))))
}

func print(b board) {
	fmt.Println(b.String())
}

func main() {
	current := newBoard(5)

	generation := 1
	ticker := time.Tick(time.Second / 2)

	for {
		<-ticker
		fmt.Println("Generation", generation, "=", len(current), "x", len(current[0]))
		fmt.Println()
		print(current)
		evolve(current)

		if next := expand(current); next != nil {
			if reflect.DeepEqual(current, next) {
				fmt.Println("Reached a stalemate. Quitting.")
				fmt.Println()
				print(next)
				return
			}

			current = next
		}
		generation++
	}

}
