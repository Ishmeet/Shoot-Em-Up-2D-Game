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
var emptyImage *ebiten.Image

type box struct {
	x0, y0, w, h int
}

// Game ...
type Game struct {
	x int
	y int
	b box
}

func init() {
	crossHair, _, _ = ebitenutil.NewImageFromFile("crosshair.png", ebiten.FilterDefault)
	emptyImage, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	emptyImage.Fill(color.RGBA{0xFF, 0, 0, 0x01})
}

// NewGame ...
func NewGame() *Game {
	g := &Game{}
	g.b.x0 = 100
	g.b.y0 = 100
	g.b.w = 100
	g.b.h = 100
	return g
}

func (g *Game) insideBox() bool {
	if g.x < g.b.x0+g.b.w &&
		g.x+10 > g.b.x0 &&
		g.y < g.b.y0+g.b.h &&
		g.y+10 > g.b.y0 {
		// collision detected!
		return true
	}
	return false
}

// Update ...
func (g *Game) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return fmt.Errorf("Escape Pressed")
	}
	g.x, g.y = ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Press"))
	}

	return nil
}

// Draw ...
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(g.b.w), float64(g.b.h))
	op.GeoM.Translate(float64(g.b.x0), float64(g.b.y0))
	if g.insideBox() {
		op.ColorM.Scale(0x00, 153, 0x0, 100)
	} else {
		op.ColorM.Scale(0xFF, 0x0, 0x0, 0xFF)
	}
	screen.DrawImage(emptyImage, &op)

	op.GeoM.Reset()
	op.ColorM.Reset()
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
