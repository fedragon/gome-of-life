package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/fedragon/gome-of-life/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Game struct {
	board      *core.Board
	generation int
	pixels     []byte
}

func (g *Game) Update() error {
	g.board.Evolve()
	g.generation++

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}

	g.board.TakeSnapshot(g.pixels)

	screen.WritePixels(g.pixels)

	ebitenutil.DebugPrint(screen, strconv.Itoa(g.generation))
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		board:      core.NewBoard(screenWidth, screenHeight, screenWidth*screenHeight/25),
		generation: 1,
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Game of Life")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
