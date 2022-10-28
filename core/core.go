package core

import (
	"math/rand"
)

type Board struct {
	state [][]bool
}

func NewBoard(width, height, maxAlive int) *Board {
	b := make([][]bool, height)

	for x := 0; x < height; x++ {
		b[x] = make([]bool, width)
	}

	for i := 0; i < maxAlive; i++ {
		y := rand.Intn(height)
		x := rand.Intn(width)

		b[y][x] = true
	}

	return &Board{state: b}
}

func (b *Board) Evolve() {
	for y := 0; y < len(b.state); y++ {
		for x := 0; x < len(b.state[y]); x++ {
			b.evolveCell(coord{x, y})
		}
	}
}

func (b *Board) evolveCell(c coord) {
	x, y := c.x, c.y
	aliveNeighbours := 0
	for _, n := range neighbours(c, len(b.state), len(b.state[y])) {
		if b.state[n.y][n.x] {
			aliveNeighbours++
		}
	}

	switch {
	case aliveNeighbours < 2:
		b.state[y][x] = false
	case (aliveNeighbours == 2 || aliveNeighbours == 3) && b.state[y][x]:
		b.state[y][x] = true
	case aliveNeighbours > 3:
		b.state[y][x] = false
	case aliveNeighbours == 3:
		b.state[y][x] = true
	default:
		b.state[y][x] = false
	}
}

// TakeSnapshot creates a snapshot of the board state on a pixel byte array where each pixel is represented by a 4-bytes RGBA point.
func (b *Board) TakeSnapshot(pixels []byte) {
	for y := 0; y < len(b.state); y++ {
		width := len(b.state[y])

		for x := 0; x < width; x++ {
			var v byte = 0
			if b.state[y][x] {
				v = 0xff
			}

			offset := y*width + x
			for i := 0; i < 3; i++ {
				pixels[offset*4+i] = v
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

func neighbours(origin coord, height, width int) []coord {
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
		if r.x >= 0 && r.x < width && r.y >= 0 && r.y < height {
			res = append(res, r)
		}
	}

	return res
}
