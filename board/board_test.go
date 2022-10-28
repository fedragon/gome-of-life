package board

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	T = true
	F = false
)

type args struct {
	b *Board
	c cell
}

type testCase struct {
	name string
	args args
	want *Board
}

func TestBoard_evolveCell_underpopulation_zero(t *testing.T) {
	want := empty(2, 2)

	gen := func(row, col int) args {
		board := empty(2, 2)
		board.grid[row][col] = true

		return args{b: board, c: cell{col: col, row: row}}
	}

	tests := []testCase{
		{"underpopulation_zero_top_left", gen(0, 0), want},
		{"underpopulation_zero_top_right", gen(0, 1), want},
		{"underpopulation_zero_bottom_left", gen(1, 0), want},
		{"underpopulation_zero_bottom_right", gen(1, 1), want},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.b.evolveCell(tt.args.c)
			if !assert.ElementsMatch(t, tt.args.b.grid, tt.want.grid) {
				t.Log("Got:")
				t.Log(stringify(tt.args.b))
				t.Log("Expected:")
				t.Log(stringify(tt.want))
				t.Fail()
			}
		})
	}
}

func TestBoard_evolveCell_underpopulation_one(t *testing.T) {
	want := func(alive cell) *Board {
		board := empty(2, 2)
		board.grid[alive.row][alive.col] = true

		return board
	}

	gen := func(orig cell, other cell) args {
		board := empty(2, 2)
		board.grid[orig.row][orig.col] = true
		board.grid[other.row][other.col] = true

		return args{b: board, c: cell{col: orig.col, row: orig.row}}
	}

	tests := []testCase{
		{"underpopulation_one_top_left", gen(cell{0, 0}, cell{0, 1}), want(cell{0, 1})},
		{"underpopulation_one_top_right", gen(cell{0, 1}, cell{1, 0}), want(cell{1, 0})},
		{"underpopulation_one_bottom_left", gen(cell{1, 0}, cell{0, 0}), want(cell{0, 0})},
		{"underpopulation_one_bottom_right", gen(cell{1, 1}, cell{0, 0}), want(cell{0, 0})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.b.evolveCell(tt.args.c)
			if !assert.ElementsMatch(t, tt.args.b.grid, tt.want.grid) {
				t.Log("Got:")
				t.Log(stringify(tt.args.b))
				t.Log("Expected:")
				t.Log(stringify(tt.want))
				t.Fail()
			}
		})
	}
}

func Test_evolveCell(t *testing.T) {
	tests := []testCase{
		{
			name: "any living cell with no alive neighbours dies, as if by underpopulation",
			args: args{
				b: from([][]bool{
					{F, F, F},
					{F, T, F},
					{F, F, F},
				}),
				c: cell{1, 1},
			},
			want: from([][]bool{
				{F, F, F},
				{F, F, F},
				{F, F, F},
			}),
		},
		{
			name: "any living cell with one alive neighbour dies, as if by underpopulation",
			args: args{
				b: from([][]bool{
					{F, F, T},
					{F, T, F},
					{F, F, F},
				}),
				c: cell{1, 1},
			},
			want: from([][]bool{
				{F, F, T},
				{F, F, F},
				{F, F, F},
			}),
		},
		{
			name: "any living cell with two alive neighbours lives on to the next generation",
			args: args{
				b: from([][]bool{
					{F, T, F},
					{F, T, F},
					{T, F, F},
				}),
				c: cell{1, 1},
			},
			want: from([][]bool{
				{F, T, F},
				{F, T, F},
				{T, F, F},
			}),
		},
		{
			name: "any living cell with three alive neighbours lives on to the next generation",
			args: args{
				b: from([][]bool{
					{F, T, F},
					{F, T, T},
					{T, F, F},
				}),
				c: cell{1, 1},
			},
			want: from([][]bool{
				{F, T, F},
				{F, T, T},
				{T, F, F},
			}),
		},
		{
			name: "any living cell with more than three alive neighbours dies, as if by overpopulation",
			args: args{
				b: from([][]bool{
					{F, T, T},
					{F, T, F},
					{T, F, T},
				}),
				c: cell{1, 1},
			},
			want: from([][]bool{
				{F, T, T},
				{F, F, F},
				{T, F, T},
			}),
		},
		{
			name: "any dead cell with exactly three alive neighbours becomes a living cell, as if by reproduction",
			args: args{
				b: from([][]bool{
					{T, T, F},
					{F, F, F},
					{F, T, F},
				}),
				c: cell{1, 1},
			},
			want: from([][]bool{
				{T, T, F},
				{F, T, F},
				{F, T, F},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.b.evolveCell(tt.args.c)
			if !assert.ElementsMatch(t, tt.args.b.grid, tt.want.grid) {
				t.Log(tt.name)
				t.Log("Got:")
				t.Log("\n", stringify(tt.args.b))
				t.Log("Expected:")
				t.Log("\n", stringify(tt.want))
				t.Fail()
			}
		})
	}
}

func Test_neighbours(t *testing.T) {
	type args struct {
		origin cell
		rows   int
		cols   int
	}
	tests := []struct {
		name string
		args args
		want []cell
	}{
		{
			name: "returns 3 neighbours when cell is in (0, 0) of a 3x3 board",
			args: args{cell{0, 0}, 3, 3},
			want: []cell{{1, 0}, {1, 1}, {0, 1}},
		},
		{
			name: "returns 5 neighbours when cell is in (0, 1) of a 3x3 board",
			args: args{cell{0, 1}, 3, 3},
			want: []cell{{0, 0}, {1, 0}, {1, 1}, {1, 2}, {0, 2}},
		},
		{
			name: "returns 3 neighbours when cell is in (0, 2) of a 3x3 board",
			args: args{cell{0, 2}, 3, 3},
			want: []cell{{0, 1}, {1, 1}, {1, 2}},
		},
		{
			name: "returns 5 neighbours when cell is in (1, 0) of a 3x3 board",
			args: args{cell{1, 0}, 3, 3},
			want: []cell{{2, 0}, {2, 1}, {1, 1}, {0, 1}, {0, 0}},
		},
		{
			name: "returns 8 neighbours when cell is in (1,1) of a 3x3 board",
			args: args{cell{1, 1}, 3, 3},
			want: []cell{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}, {1, 2}, {0, 2}, {0, 1}},
		},
		{
			name: "returns 5 neighbours when cell is in (1, 2) of a 3x3 board",
			args: args{cell{1, 2}, 3, 3},
			want: []cell{{0, 1}, {1, 1}, {2, 1}, {2, 2}, {0, 2}},
		},
		{
			name: "returns 3 neighbours when cell is in (2, 0) of a 3x3 board",
			args: args{cell{2, 0}, 3, 3},
			want: []cell{{2, 1}, {1, 1}, {1, 0}},
		},
		{
			name: "returns 5 neighbours when cell is in (2, 1) of a 3x3 board",
			args: args{cell{2, 1}, 3, 3},
			want: []cell{{1, 0}, {2, 0}, {2, 2}, {1, 2}, {1, 1}},
		},
		{
			name: "returns 3 neighbours when cell is in (2, 2) of a 3x3 board",
			args: args{cell{2, 2}, 3, 3},
			want: []cell{{1, 1}, {2, 1}, {1, 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.want, neighbours(tt.args.origin, tt.args.rows, tt.args.cols))
		})
	}
}

func from(grid [][]bool) *Board {
	cols := -1

	if len(grid) == 0 {
		panic("grid cannot be empty")
	}

	for _, r := range grid {
		if cols > -1 && cols != len(r) {
			panic("all rows must have the same length")
		}

		cols = len(r)
	}

	return &Board{
		grid: grid,
		rows: len(grid),
		cols: cols,
	}
}

func stringify(b *Board) (result string) {
	var symbol string

	result += "\n"

	for i := 0; i < len(b.grid); i++ {
		for j := 0; j < len(b.grid[i]); j++ {

			if b.grid[i][j] {
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
