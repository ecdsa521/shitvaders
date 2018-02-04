package main

import (
	"fmt"
	"image"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func loadPicture(path string) pixel.Picture {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	return pixel.PictureDataFromImage(img)
}

type element struct {
	x         float64
	y         float64
	p         pixel.Picture
	s         *pixel.Sprite
	active    bool
	direction float64
}

func (e *element) draw(win *pixelgl.Window) {
	if e.active {
		e.s.Draw(win, pixel.IM.Moved(pixel.Vec{X: e.x, Y: e.y}))
	}
}
func (e *element) loadSprite(path string) {
	e.p = loadPicture(path)
	e.s = pixel.NewSprite(e.p, e.p.Bounds())
}

type gameData struct {
	started   bool
	speed     float64
	player    *element
	over      *element
	enemies   []*element
	bullet    *element
	bomb      *element
	direction float64
	row       float64
}

func (game *gameData) makeEnemies(enemyCount int) {
	game.enemies = nil
	row := 0
	col := 0
	game.direction = 1
	for i := 0; i < enemyCount; i++ {
		col++
		if i%10 == 0 {
			row++
			col = 0
		}
		enemy := &element{
			x:         100 + float64(col*60),
			y:         700 - float64(row*50),
			direction: 1,
			active:    true,
		}
		enemy.loadSprite("invader.png")
		game.enemies = append(game.enemies, enemy)

		fmt.Printf("%v %v", row, col)
	}

}
func (game *gameData) shoot() {
	if !game.bullet.active {
		game.bullet.active = true
		game.bullet.x = game.player.x
		game.bullet.y = game.player.y + 50
	}
}

func (game *gameData) dropBomb() {
	rand.Seed(time.Now().Unix())
	if game.started && game.bomb.active == false {

		e := game.randInvader()

		game.bomb.active = true
		game.bomb.x = game.enemies[e].x
		game.bomb.y = game.enemies[e].y - 5

	}
}
func (game *gameData) checkCollision() {

	for _, e := range game.enemies {
		minX := e.x - 25
		maxX := e.x + 25
		minY := e.y - 25
		maxY := e.y + 25
		if game.bullet.active && game.bullet.x > minX && game.bullet.x < maxX &&
			game.bullet.y >= minY && game.bullet.y <= maxY {
			if e.active {
				e.active = false
				game.bullet.active = false
			}
		}

	}

	if game.bullet.y >= 1000 {
		game.bullet.active = false
	}

}
func (game *gameData) checkBomb() {
	minX := game.player.x - 65
	maxX := game.player.x + 65

	if game.bomb.active && game.bomb.x >= minX && game.bomb.x <= maxX && game.bomb.y <= 100 && game.bomb.y >= 50 {
		game.bomb.active = false
		game.started = false
	}
}
func (game *gameData) checkWin() bool {
	for _, e := range game.enemies {
		if e.active {
			return false
		}
	}
	return true
}
func (game *gameData) checkLoss() bool {
	for _, e := range game.enemies {
		if e.y <= 200 && e.active {
			return true
		}
	}

	return false
}

//find left-most and right-most invader
func (game *gameData) findBorder() (float64, float64) {
	var maxE, minE float64
	minE = 1000
	for _, e := range game.enemies {
		if e.x < minE {
			minE = e.x
		}
		if e.x > maxE {
			maxE = e.x
		}
	}
	return minE, maxE
}

//find random active invader
func (game *gameData) randInvader() int {
	var active []int
	for id, e := range game.enemies {
		if e.active {
			active = append(active, id)
		}
	}
	rand.Seed(time.Now().Unix())
	i := active[rand.Intn(len(active))]
	return i
}

func (game *gameData) update(win *pixelgl.Window) {
	game.player.draw(win)
	game.bullet.draw(win)
	game.bomb.draw(win)
	minX, maxX := game.findBorder()
	if maxX >= 900 && game.direction == 1 {
		game.direction = -1
		game.row++
		for _, e := range game.enemies {
			e.y -= 50
		}
	} else if minX <= 100 && game.direction == -1 {
		game.direction = 1
		game.row++
		for _, e := range game.enemies {
			e.y -= 50
		}
	}

	for _, e := range game.enemies {
		e.x += game.direction * game.speed
		e.draw(win)
	}

	if game.bomb.y <= 10 {
		game.bomb.active = false
	}

}

func (game *gameData) move(key pixelgl.Button) {
	player := game.player
	switch key {

	case pixelgl.KeyRight:
		player.x += 10
		if player.x >= 900 {
			player.x = 900
		}

	case pixelgl.KeyLeft:
		player.x -= 10
		if player.x <= 100 {
			player.x = 100
		}

	case pixelgl.KeySpace:
		game.shoot()

	case pixelgl.KeyUnknown:
		//nothing

	default:

		fmt.Printf("%v\n", key)
	}
}

func loadSprites() {

}

func loadTitle() {

}

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "Shitvaders",
		Bounds: pixel.R(0, 0, 1000, 800),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.Clear(colornames.Black)

	game := &gameData{
		started: false,
		player: &element{
			x:      500,
			y:      100,
			active: true,
		},
		over: &element{
			x: 100,
			y: 700,
		},
		speed:  1,
		bullet: &element{x: 0, y: 0},
		bomb:   &element{x: 0, y: 0},
	}
	game.over.loadSprite("gameover.png")
	game.player.loadSprite("player.png")
	game.bullet.loadSprite("bullet.png")
	game.bomb.loadSprite("bomb.png")

	//last := time.Now()

	for !win.Closed() {

		//dt := time.Since(last).Seconds()
		//last = time.Now()
		if game.started {
			if game.checkWin() {
				game.speed++
				game.makeEnemies(40)
			}
			game.move(pollKeys(win))

			win.Clear(colornames.Black)
			if game.bullet.active {
				game.bullet.y += 10
				game.checkCollision()
			}
			if game.bomb.active {
				game.bomb.y -= 5
				game.checkBomb()
			}

			game.update(win)
			game.dropBomb()
			win.Update()
		} else {
			if pollKeys(win) == pixelgl.KeyEscape {
				game.makeEnemies(40)
				game.started = true
			}
			win.Clear(colornames.Black)
			game.over.s.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

			win.Update()
			//start screen
		}
	}
}

func pollKeys(win *pixelgl.Window) pixelgl.Button {
	if win.Pressed(pixelgl.KeyLeft) {
		return pixelgl.KeyLeft
	}
	if win.Pressed(pixelgl.KeyRight) {
		return pixelgl.KeyRight
	}
	if win.Pressed(pixelgl.KeySpace) {
		return pixelgl.KeySpace
	}
	if win.Pressed(pixelgl.KeyEscape) {
		return pixelgl.KeyEscape
	}
	return pixelgl.KeyUnknown
}
func main() {
	pixelgl.Run(run)
}
