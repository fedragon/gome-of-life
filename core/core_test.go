package core

import (
	"fmt"
	"reflect"
	"testing"
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
	want := &Board{state: [][]bool{
		{F, F},
		{F, F},
	}}

	gen := func(x, y int) args {
		board := deadBoard()
		board.state[x][y] = true

		return args{board, cell{x, y}}
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
			if !reflect.DeepEqual(tt.args.b, tt.want) {
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
		board := deadBoard()
		board.state[alive.col][alive.row] = true

		return board
	}

	gen := func(c cell, n cell) args {
		board := deadBoard()
		board.state[c.col][c.row] = true
		board.state[n.col][n.row] = true

		return args{board, cell{c.col, c.row}}
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
			if !reflect.DeepEqual(tt.args.b, tt.want) {
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
		{"any living cell with no alive neighbours dies, as if by underpopulation",
			args{
				b: &Board{state: [][]bool{
					{F, F, F},
					{F, T, F},
					{F, F, F},
				}},
				c: cell{1, 1},
			},
			&Board{state: [][]bool{
				{F, F, F},
				{F, F, F},
				{F, F, F},
			}},
		},
		{"any living cell with one alive neighbour dies, as if by underpopulation",
			args{
				b: &Board{state: [][]bool{
					{F, F, T},
					{F, T, F},
					{F, F, F},
				}},
				c: cell{1, 1},
			},
			&Board{state: [][]bool{
				{F, F, T},
				{F, F, F},
				{F, F, F},
			}},
		},
		{"any living cell with two alive neighbours lives on to the next generation",
			args{
				b: &Board{state: [][]bool{
					{F, T, F},
					{F, T, F},
					{T, F, F},
				}},
				c: cell{1, 1},
			},
			&Board{state: [][]bool{
				{F, T, F},
				{F, T, F},
				{T, F, F},
			}},
		},
		{"any living cell with three alive neighbours lives on to the next generation",
			args{
				b: &Board{state: [][]bool{
					{F, T, F},
					{F, T, T},
					{T, F, F},
				}},
				c: cell{1, 1},
			},
			&Board{state: [][]bool{
				{F, T, F},
				{F, T, T},
				{T, F, F},
			}},
		},
		{"any living cell with more than three alive neighbours dies, as if by overpopulation",
			args{
				b: &Board{state: [][]bool{
					{F, T, T},
					{F, T, F},
					{T, F, T},
				}},
				c: cell{1, 1},
			},
			&Board{state: [][]bool{
				{F, T, T},
				{F, F, F},
				{T, F, T},
			}},
		},
		{"any dead cell with exactly three alive neighbours becomes a living cell, as if by reproduction",
			args{
				b: &Board{state: [][]bool{
					{T, T, F},
					{F, F, F},
					{F, T, F},
				}},
				c: cell{1, 1},
			},
			&Board{state: [][]bool{
				{T, T, F},
				{F, T, F},
				{F, T, F},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.b.evolveCell(tt.args.c)
			if !reflect.DeepEqual(tt.args.b, tt.want) {
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
		{"returns 8 neighbours when cell is in the center of a 3x3 board",
			args{cell{1, 1}, 3, 3},
			[]cell{cell{0, 0}, cell{1, 0}, cell{2, 0}, cell{2, 1}, cell{2, 2}, cell{1, 2}, cell{0, 2}, cell{0, 1}}},
		{"returns 5 neighbours when cell is in the middle of the top row of a 3x3 board",
			args{cell{0, 1}, 3, 3},
			[]cell{cell{0, 0}, cell{1, 0}, cell{1, 1}, cell{1, 2}, cell{0, 2}}},
		{"returns 5 neighbours when cell is in the middle of the bottom row of a 3x3 board",
			args{cell{2, 1}, 3, 3},
			[]cell{cell{1, 0}, cell{2, 0}, cell{2, 2}, cell{1, 2}, cell{1, 1}}},
		{"returns 5 neighbours when cell is on the left of the middle row of a 3x3 board",
			args{cell{1, 0}, 3, 3},
			[]cell{cell{2, 0}, cell{2, 1}, cell{1, 1}, cell{0, 1}, cell{0, 0}}},
		{"returns 5 neighbours when cell is on the right of the middle row of a 3x3 board",
			args{cell{1, 2}, 3, 3},
			[]cell{cell{0, 1}, cell{1, 1}, cell{2, 1}, cell{2, 2}, cell{0, 2}}},
		{"returns 3 neighbours when cell is on the left of the top row of a 3x3 board",
			args{cell{0, 0}, 3, 3},
			[]cell{cell{1, 0}, cell{1, 1}, cell{0, 1}}},
		{"returns 3 neighbours when cell is on the right of the top row of a 3x3 board",
			args{cell{0, 2}, 3, 3},
			[]cell{cell{0, 1}, cell{1, 1}, cell{1, 2}}},
		{"returns 3 neighbours when cell is on the left of the bottom row of a 3x3 board",
			args{cell{2, 0}, 3, 3},
			[]cell{cell{2, 1}, cell{1, 1}, cell{1, 0}}},
		{"returns 3 neighbours when cell is on the left of the bottom row of a 3x3 board",
			args{cell{2, 2}, 3, 3},
			[]cell{cell{1, 1}, cell{2, 1}, cell{1, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := neighbours(tt.args.origin, tt.args.rows, tt.args.cols); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("neighbours() = %v, want %v", got, tt.want)
			}
		})
	}
}

func stringify(b *Board) (result string) {
	var symbol string

	result += "\n"

	for i := 0; i < len(b.state); i++ {
		for j := 0; j < len(b.state[i]); j++ {

			if b.state[i][j] {
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

func deadBoard() *Board {
	state := make([][]bool, 2)
	for i := 0; i < 2; i++ {
		state[i] = make([]bool, 2)
	}
	return &Board{state}
}
