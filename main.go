package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/fedragon/gome-of-life/board"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Game struct {
	board      *board.Board
	generation int
	paused     bool
	pixels     []byte
}

func (g *Game) Update() error {
	if !g.paused {
		g.board.Evolve()
		g.generation++

	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.paused = !g.paused
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}

	g.board.TakeSnapshot(g.pixels)

	screen.WritePixels(g.pixels)

	ebitenutil.DebugPrintAt(screen, strconv.Itoa(g.generation), 10, 10)
	ebitenutil.DebugPrintAt(screen, "Press SPACE to pause or resume", 60, 200)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		board:      board.NewBoard(screenWidth, screenHeight, screenWidth*screenHeight/20),
		generation: 1,
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Game of Life")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
