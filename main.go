package main

import (
	"fmt"
	"image/color"
	_ "image/jpeg"
	"time"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten/inpututil"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 640
	screenHeight = 480
	// MAXBOXES no. of boxes...
	MAXBOXES = 6
)

var crossHair *ebiten.Image
var emptyImage [MAXBOXES]*ebiten.Image
var healthBar *ebiten.Image
var enemy1 [MAXBOXES]*ebiten.Image

type box struct {
	x0, y0, w, h  int
	cw, ch        int
	alive         bool
	staydeadTimer int
	aliveTimer    int
	shoot         bool
}

type click struct {
	x, y int
}

// Game ...
type Game struct {
	x      int
	y      int
	b      [MAXBOXES]box
	c      click
	health int
}

func init() {
	crossHair, _, _ = ebitenutil.NewImageFromFile("crosshair.png", ebiten.FilterDefault)
	for i := 0; i < len(enemy1); i++ {
		enemy1[i], _, _ = ebitenutil.NewImageFromFile("enemy1_50_50.png", ebiten.FilterDefault)
	}
	for i := 0; i < len(emptyImage); i++ {
		emptyImage[i], _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
		emptyImage[i].Fill(color.RGBA{0xFF, 0, 0, 0x01})
	}
	healthBar, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	healthBar.Fill(color.RGBA{0xFF, 0, 0, 0x01})
}

// NewGame ...
func NewGame() *Game {
	g := &Game{}
	g.health = 100

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

	g.b[4].x0 = 350
	g.b[4].y0 = 250
	g.b[4].w = 50
	g.b[4].h = 50

	g.b[5].x0 = 450
	g.b[5].y0 = 100
	g.b[5].w = 50
	g.b[5].h = 50

	for i := 0; i < len(g.b); i++ {
		g.b[i].alive = true
		g.b[i].shoot = false
	}
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
		if g.b[i].alive {
			op.GeoM.Reset()
			op.ColorM.Reset()
			op.GeoM.Translate(float64(g.b[i].x0), float64(g.b[i].y0))
			if i == g.clickedInsideBox() {
				msg = msg + " " + fmt.Sprintf("%d", i)
				g.enemyKilled(i)
				op.ColorM.Scale(0xFF, 0x0, 0x0, 0xFF)
				// op.ColorM.Scale(0x00, 153, 0x0, 100)
			} else {
				// op.ColorM.Scale(0xFF, 0x0, 0x0, 0xFF)
			}
			if g.b[i].shoot {
				op.ColorM.Scale(0xFF, 0xFF, 0x0, 100)
				g.b[i].shoot = false
			}
			screen.DrawImage(enemy1[i], &op)
		}
	}

	op.GeoM.Reset()
	op.ColorM.Reset()
	w, h := crossHair.Size()
	op.GeoM.Translate(float64(g.x-w/2), float64(g.y-h/2))
	screen.DrawImage(crossHair, &op)

	op.GeoM.Reset()
	op.ColorM.Reset()
	op.GeoM.Scale(float64(g.health), 20)
	op.GeoM.Translate((screenWidth/2)-50, 20)
	op.ColorM.Scale(0x00, 0xFF, 0x00, 0xAA)
	screen.DrawImage(healthBar, &op)

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

func (g *Game) youGotShot(i int) {
	g.health -= 20
	if g.health <= 0 {
		g.health = 0
	}
	g.b[i].shoot = true
}

func (g *Game) enemyRespawn(i int) {
	g.b[i].h = 50
	g.b[i].w = 50
	g.b[i].alive = true
}

func (g *Game) enemyKilled(i int) {
	g.b[i].h = 0
	g.b[i].w = 0
	g.b[i].alive = false
	g.c.x = 0
	g.c.y = 0
}

func (g *Game) timer() {
	for range time.Tick(time.Second) {
		for i := 0; i < len(g.b); i++ {
			if !g.b[i].alive {
				if g.b[i].staydeadTimer >= 3 {
					g.enemyRespawn(i)
					g.b[i].staydeadTimer = 0
				}
				g.b[i].staydeadTimer++
				// Reset alive timer here
				g.b[i].aliveTimer = 0
			} else {
				if g.b[i].aliveTimer >= 3 {
					g.youGotShot(i)
					g.b[i].aliveTimer = 0
				}
				g.b[i].aliveTimer++
			}
		}
	}
}

// Layout ...
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowTitle("Shoot em up")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	g := NewGame()
	go g.timer()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
