package board

import (
	"math/rand"
)

type Board struct {
	grid [][]bool
	rows int
	cols int
}

func empty(rows, cols int) *Board {
	grid := make([][]bool, rows)
	for i := 0; i < rows; i++ {
		grid[i] = make([]bool, cols)
	}

	return &Board{grid: grid, rows: rows, cols: cols}
}

func NewBoard(width, height, maxAlive int) *Board {
	b := empty(height, width)

	for i := 0; i < maxAlive; i++ {
		row := rand.Intn(height)
		col := rand.Intn(width)

		b.grid[row][col] = true
	}

	return b
}

func (b *Board) Evolve() {
	for row := 0; row < b.rows; row++ {
		for col := 0; col < b.cols; col++ {
			b.evolveCell(cell{row: row, col: col})
		}
	}
}

func (b *Board) evolveCell(origin cell) {
	row, col := origin.row, origin.col
	aliveNeighbours := 0
	for _, n := range neighbours(cell{row: row, col: col}, b.rows, b.cols) {
		if b.grid[n.row][n.col] {
			aliveNeighbours++
		}
	}

	switch {
	case aliveNeighbours == 2 && b.grid[row][col]:
		b.grid[row][col] = true
	case aliveNeighbours == 3:
		b.grid[row][col] = true
	default:
		b.grid[row][col] = false
	}
}

// TakeSnapshot creates a snapshot of the board state on a pixel byte array where each pixel is represented by a 4-bytes RGBA point.
func (b *Board) TakeSnapshot(pixels []byte) {
	for row := 0; row < b.rows; row++ {
		for col := 0; col < b.cols; col++ {
			var v byte = 0
			if b.grid[row][col] {
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
	row int
	col int
}

func neighbours(origin cell, rows, cols int) []cell {
	res := make([]cell, 0)
	for row := -1; row < 2; row++ {
		for col := -1; col < 2; col++ {
			if row == 0 && col == 0 {
				continue
			}

			r := origin.row + row
			c := origin.col + col

			if r >= 0 && r < rows && c >= 0 && c < cols {
				res = append(res, cell{row: r, col: c})
			}
		}
	}

	return res
}
