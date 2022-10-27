package core

import (
	"math/rand"
)

type Board struct {
	state [][]bool
}

func NewBoard(rows, cols, maxAlive int) *Board {
	b := make([][]bool, rows)

	for x := 0; x < rows; x++ {
		b[x] = make([]bool, cols)
	}

	for k := 0; k < maxAlive; k++ {
		x := rand.Intn(rows)
		y := rand.Intn(cols)

		b[x][y] = true
	}

	return &Board{state: b}
}

func (b *Board) Evolve() {
	for x := 0; x < len(b.state); x++ {
		for y := 0; y < len(b.state[x]); y++ {
			b.evolveCell(coord{x, y})
		}
	}
}

func (b *Board) evolveCell(c coord) {
	x, y := c.x, c.y
	aliveNeighbours := 0
	for _, n := range neighbours(c, len(b.state), len(b.state[x])) {
		if b.state[n.x][n.y] {
			aliveNeighbours++
		}
	}

	switch {
	case b.state[x][y] && aliveNeighbours == 2:
		b.state[x][y] = true
	case aliveNeighbours == 3:
		b.state[x][y] = true
	default:
		b.state[x][y] = false
	}
}

// TakeSnapshot creates a snapshot of the board state on a pixel byte array where each pixel is represented by a 4-bytes RGBA point.
func (b *Board) TakeSnapshot(pixels []byte) {
	for x := 0; x < len(b.state); x++ {
		for y := 0; y < len(b.state[x]); y++ {
			var v byte = 0
			if b.state[x][y] {
				v = 0xff
			}

			for k := 0; k < 3; k++ {
				pixels[4*x*y+k] = v
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
