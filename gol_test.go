package main

import (
	"reflect"
	"testing"
)

func Test_growUp(t *testing.T) {
	type args struct {
		b board
	}
	tests := []struct {
		name string
		args args
		want board
	}{
		{name: "grows up if any cell in it is alive",
			args: args{board{[]bool{false, true, false}, []bool{false, false, false}}},
			want: board{[]bool{false, false, false}, []bool{false, true, false}, []bool{false, false, false}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := growUp(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("growUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_growDown(t *testing.T) {
	type args struct {
		b board
	}
	tests := []struct {
		name string
		args args
		want board
	}{
		{name: "grows down if any cell in it is alive",
			args: args{board{[]bool{false, false, false}, []bool{false, true, false}}},
			want: board{[]bool{false, false, false}, []bool{false, true, false}, []bool{false, false, false}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := growDown(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("growDown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_growLeft(t *testing.T) {
	type args struct {
		b board
	}
	tests := []struct {
		name string
		args args
		want board
	}{
		{name: "grows to the left if any cell in it is alive",
			args: args{board{[]bool{false, false}, []bool{true, false}}},
			want: board{[]bool{false, false, false}, []bool{false, true, false}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := growLeft(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("growLeft() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_growRight(t *testing.T) {
	type args struct {
		b board
	}
	tests := []struct {
		name string
		args args
		want board
	}{
		{name: "grows to the right if any cell in it is alive",
			args: args{board{[]bool{false, true}, []bool{false, false}}},
			want: board{[]bool{false, true, false}, []bool{false, false, false}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := growRight(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("growRight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_evolveCell(t *testing.T) {
	type args struct {
		b board
		c coord
	}
	tests := []struct {
		name string
		args args
		want board
	}{
		{"any live cell with fewer than two live neighbours dies, as if by underpopulation",
			args{b: [][]bool{[]bool{false, false, false}, []bool{false, true, false}, []bool{false, false, false}}, c: coord{1, 1}},
			[][]bool{[]bool{false, false, false}, []bool{false, false, false}, []bool{false, false, false}},
		},
		{"any live cell with two or three live neighbours lives on to the next generation",
			args{b: [][]bool{[]bool{false, true, false}, []bool{false, true, false}, []bool{true, false, false}}, c: coord{1, 1}},
			[][]bool{[]bool{false, true, false}, []bool{false, true, false}, []bool{true, false, false}},
		},
		{"any live cell with more than three live neighbours dies, as if by overpopulation",
			args{b: [][]bool{[]bool{false, true, true}, []bool{false, true, false}, []bool{true, false, true}}, c: coord{1, 1}},
			[][]bool{[]bool{false, true, true}, []bool{false, false, false}, []bool{true, false, true}},
		},
		{"any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction",
			args{b: [][]bool{[]bool{true, true, false}, []bool{false, false, false}, []bool{false, true, false}}, c: coord{1, 1}},
			[][]bool{[]bool{true, true, false}, []bool{false, true, false}, []bool{false, true, false}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evolveCell(tt.args.b, tt.args.c)
			if !reflect.DeepEqual(tt.args.b, tt.want) {
				t.Log("Got:")
				t.Log(tt.args.b.String())
				t.Log("Expected:")
				t.Log(tt.want.String())
				t.Fail()
			}
		})
	}
}

func Test_neighbours(t *testing.T) {
	type args struct {
		origin coord
		rows   int
		cols   int
	}
	tests := []struct {
		name string
		args args
		want []coord
	}{
		{"returns 8 neighbours when cell is in the center of a 3x3 board",
			args{coord{1, 1}, 3, 3},
			[]coord{coord{0, 0}, coord{1, 0}, coord{2, 0}, coord{2, 1}, coord{2, 2}, coord{1, 2}, coord{0, 2}, coord{0, 1}}},
		{"returns 5 neighbours when cell is in the middle of the top row of a 3x3 board",
			args{coord{0, 1}, 3, 3},
			[]coord{coord{0, 0}, coord{1, 0}, coord{1, 1}, coord{1, 2}, coord{0, 2}}},
		{"returns 5 neighbours when cell is in the middle of the bottom row of a 3x3 board",
			args{coord{2, 1}, 3, 3},
			[]coord{coord{1, 0}, coord{2, 0}, coord{2, 2}, coord{1, 2}, coord{1, 1}}},
		{"returns 5 neighbours when cell is on the left of the middle row of a 3x3 board",
			args{coord{1, 0}, 3, 3},
			[]coord{coord{2, 0}, coord{2, 1}, coord{1, 1}, coord{0, 1}, coord{0, 0}}},
		{"returns 5 neighbours when cell is on the right of the middle row of a 3x3 board",
			args{coord{1, 2}, 3, 3},
			[]coord{coord{0, 1}, coord{1, 1}, coord{2, 1}, coord{2, 2}, coord{0, 2}}},
		{"returns 3 neighbours when cell is on the left of the top row of a 3x3 board",
			args{coord{0, 0}, 3, 3},
			[]coord{coord{1, 0}, coord{1, 1}, coord{0, 1}}},
		{"returns 3 neighbours when cell is on the right of the top row of a 3x3 board",
			args{coord{0, 2}, 3, 3},
			[]coord{coord{0, 1}, coord{1, 1}, coord{1, 2}}},
		{"returns 3 neighbours when cell is on the left of the bottom row of a 3x3 board",
			args{coord{2, 0}, 3, 3},
			[]coord{coord{2, 1}, coord{1, 1}, coord{1, 0}}},
		{"returns 3 neighbours when cell is on the left of the bottom row of a 3x3 board",
			args{coord{2, 2}, 3, 3},
			[]coord{coord{1, 1}, coord{2, 1}, coord{1, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := neighbours(tt.args.origin, tt.args.rows, tt.args.cols); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("neighbours() = %v, want %v", got, tt.want)
			}
		})
	}
}
