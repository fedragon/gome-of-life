package core

import (
	"math/rand"
)

type Board struct {
	state [][]bool
	cols  int
	rows  int
}

func NewBoard(width, height, maxAlive int) *Board {
	b := make([][]bool, height)

	for i := 0; i < height; i++ {
		b[i] = make([]bool, width)
	}

	for i := 0; i < maxAlive; i++ {
		row := rand.Intn(height)
		col := rand.Intn(width)

		b[row][col] = true
	}

	return &Board{state: b, cols: width, rows: height}
}

func (b *Board) Evolve() {
	for row := 0; row < b.rows; row++ {
		for col := 0; col < b.cols; col++ {
			b.evolveCell(row, col)
		}
	}
}

func (b *Board) evolveCell(row, col int) {
	aliveNeighbours := 0
	for _, n := range neighbours(cell{row: row, col: col}, b.rows, b.cols) {
		if b.state[n.row][n.col] {
			aliveNeighbours++
		}
	}

	switch {
	case aliveNeighbours < 2:
		b.state[row][col] = false
	case (aliveNeighbours == 2 || aliveNeighbours == 3) && b.state[row][col]:
		b.state[row][col] = true
	case aliveNeighbours > 3:
		b.state[row][col] = false
	case aliveNeighbours == 3:
		b.state[row][col] = true
	default:
		b.state[row][col] = false
	}
}

// TakeSnapshot creates a snapshot of the board state on a pixel byte array where each pixel is represented by a 4-bytes RGBA point.
func (b *Board) TakeSnapshot(pixels []byte) {
	for row := 0; row < b.rows; row++ {
		for col := 0; col < b.cols; col++ {
			var v byte = 0
			if b.state[row][col] {
				v = 0xff
			}

			offset := row*b.cols + col
			for i := 0; i < 3; i++ {
				pixels[offset*4+i] = v
			}
		}
	}
}

type cell struct {
	col int
	row int
}

func up(a cell) cell {
	return cell{col: a.col, row: a.row - 1}
}

func down(a cell) cell {
	return cell{col: a.col, row: a.row + 1}
}

func left(a cell) cell {
	return cell{col: a.col - 1, row: a.row}
}

func right(a cell) cell {
	return cell{col: a.col + 1, row: a.row}
}

func neighbours(origin cell, rows, cols int) []cell {
	all := []cell{
		up(left(origin)),
		up(origin),
		up(right(origin)),
		right(origin),
		right(down(origin)),
		down(origin),
		left(down(origin)),
		left(origin),
	}

	res := make([]cell, 0)
	for _, r := range all {
		if r.col >= 0 && r.col < cols && r.row >= 0 && r.row < rows {
			res = append(res, r)
		}
	}

	return res
}
