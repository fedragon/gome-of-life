package core

import (
	"math/rand"
)

type Board struct {
	state [][]bool
}

func NewBoard(rows, cols, maxAlive int) *Board {
	b := make([][]bool, rows)

	for i := 0; i < rows; i++ {
		b[i] = make([]bool, cols)
	}

	for k := 0; k < maxAlive; k++ {
		i := rand.Intn(rows)
		j := rand.Intn(cols)

		b[i][j] = true
	}

	return &Board{state: b}
}

func (b *Board) Evolve() {
	for i := 0; i < len(b.state); i++ {
		for j := 0; j < len(b.state[i]); j++ {
			b.evolveCell(coord{i, j})
		}
	}
}

func (b *Board) evolveCell(c coord) {
	i, j := c.x, c.y
	aliveNeighbours := 0
	for _, n := range neighbours(c, len(b.state), len(b.state[i])) {
		if b.state[n.x][n.y] {
			aliveNeighbours++
		}
	}

	if b.state[i][j] {
		if aliveNeighbours < 2 || aliveNeighbours > 3 {
			b.state[i][j] = false
		} else {
			b.state[i][j] = true
		}
	} else if aliveNeighbours == 3 {
		b.state[i][j] = true
	}
}

// TakeSnapshot creates a snapshot of the board state on a pixel byte array where each pixel is represented by a 4-bytes RGBA point.
func (b *Board) TakeSnapshot(pixels []byte) {
	for i := 0; i < len(b.state); i++ {
		for j := 0; j < len(b.state[i]); j++ {
			var v byte = 0
			if b.state[i][j] {
				v = 0xff
			}

			for k := 0; k < 3; k++ {
				pixels[4*i*j+k] = v
			}
		}
	}
}

type coord struct {
	x int
	y int
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
