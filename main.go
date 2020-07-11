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
	// MAXBOXES no. of boxes...
	MAXBOXES = 4
)

var crossHair *ebiten.Image
var emptyImage [MAXBOXES]*ebiten.Image

type box struct {
	x0, y0, w, h int
	cw, ch       int
}

type click struct {
	x, y int
}

// Game ...
type Game struct {
	x int
	y int
	b [MAXBOXES]box
	c click
}

func init() {
	crossHair, _, _ = ebitenutil.NewImageFromFile("crosshair.png", ebiten.FilterDefault)
	for i := 0; i < len(emptyImage); i++ {
		emptyImage[i], _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
		emptyImage[i].Fill(color.RGBA{0xFF, 0, 0, 0x01})
	}
}

// NewGame ...
func NewGame() *Game {
	g := &Game{}
	g.b[0].x0 = 100
	g.b[0].y0 = 100
	g.b[0].w = 50
	g.b[0].h = 50

	g.b[1].x0 = 300
	g.b[1].y0 = 100
	g.b[1].w = 50
	g.b[1].h = 50

	g.b[2].x0 = 200
	g.b[2].y0 = 200
	g.b[2].w = 50
	g.b[2].h = 50

	g.b[3].x0 = 250
	g.b[3].y0 = 300
	g.b[3].w = 50
	g.b[3].h = 50
	return g
}

func (g *Game) insideBox() int {
	for i := 0; i < len(g.b); i++ {
		if g.x < g.b[i].x0+g.b[i].w &&
			g.x+10 > g.b[i].x0 &&
			g.y < g.b[i].y0+g.b[i].h &&
			g.y+10 > g.b[i].y0 {
			// collision detected!
			return i
		}
	}
	return -1
}

func (g *Game) clickedInsideBox() int {
	for i := 0; i < len(g.b); i++ {
		if g.c.x < g.b[i].x0+g.b[i].w &&
			g.c.x+10 > g.b[i].x0 &&
			g.c.y < g.b[i].y0+g.b[i].h &&
			g.c.y+10 > g.b[i].y0 {
			// collision detected!
			return i
		}
	}
	return -1
}

// Update ...
func (g *Game) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return fmt.Errorf("Escape Pressed")
	}
	g.x, g.y = ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.c.x, g.c.y = ebiten.CursorPosition()
	}

	return nil
}

// Draw ...
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	msg := fmt.Sprintf("%d,%d", g.x, g.y)
	op := ebiten.DrawImageOptions{}
	for i := 0; i < len(g.b); i++ {
		op.GeoM.Reset()
		op.ColorM.Reset()
		g.enemyMovingDown(i)
		op.GeoM.Scale(float64(g.b[i].w), float64(g.b[i].ch))
		op.GeoM.Translate(float64(g.b[i].x0), float64(g.b[i].y0))
		if i == g.clickedInsideBox() {
			msg = msg + " " + fmt.Sprintf("%d", i)
			g.enemyReset(i)
			op.ColorM.Scale(0x00, 153, 0x0, 100)
		} else {
			op.ColorM.Scale(0xFF, 0x0, 0x0, 0xFF)
		}
		screen.DrawImage(emptyImage[i], &op)
	}

	op.GeoM.Reset()
	op.ColorM.Reset()
	w, h := crossHair.Size()
	op.GeoM.Translate(float64(g.x-w/2), float64(g.y-h/2))
	screen.DrawImage(crossHair, &op)

	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) enemyMovingDown(i int) {
	if g.b[i].ch <= g.b[i].h {
		g.b[i].ch++
	}
}

func (g *Game) enemyMovingUp(i int) {
	g.b[i].h = -g.b[i].h
	if g.b[i].ch > g.b[i].h {
		g.b[i].ch--
	}
}

func (g *Game) enemyReset(i int) {
	g.b[i].ch = 0
}

// Layout ...
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowTitle("Shoot em up")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}
