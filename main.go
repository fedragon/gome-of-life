package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/fedragon/gome-of-life/core"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 160
	screenHeight = 120
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
	g.generation++
	g.board.Evolve()

	ebiten.SetWindowTitle(fmt.Sprintf("Game of Life ~ gen %d", g.generation))

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}

	g.board.TakeSnapshot(g.pixels)

	screen.WritePixels(g.pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{board: core.NewBoard(screenWidth, screenHeight, screenWidth*screenHeight/10)}

	ebiten.SetTPS(60)
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Game of Life")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
