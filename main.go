package main

import (
	"fmt"
	"image/color"
	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten/inpututil"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var crossHair *ebiten.Image

// Game ...
type Game struct {
	x int
	y int
}

func init() {
	crossHair, _, _ = ebitenutil.NewImageFromFile("crosshair.png", ebiten.FilterDefault)
}

// NewGame ...
func NewGame() *Game {
	g := &Game{}
	return g
}

// Update ...
func (g *Game) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return fmt.Errorf("Escape Pressed")
	}
	g.x, g.y = ebiten.CursorPosition()

	return nil
}

// Draw ...
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Press"))
	}
	op := ebiten.DrawImageOptions{}
	w, h := crossHair.Size()
	op.GeoM.Translate(float64(g.x-w/2), float64(g.y-h/2))
	screen.DrawImage(crossHair, &op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%d,%d", g.x, g.y))
}

// Layout ...
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowTitle("Shoot em up")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}
