package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/artem-shestakov/Snake-Go.git/models"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 600
	screenHeight = 600
	headSize     = 20
	foodRadius   = 10
)

var (
	shakeMovementPositionX = float64(0)
	shakeMovementPositionY = float64(0)
	direction              = ""

	simpleShader *ebiten.Shader
	snake        = new(models.Snake)
	foods        []models.Food
)

type Game struct {
	pressedKeys []ebiten.Key
	score       int
}

func init() {
	var err error

	simpleShader, err = ebiten.NewShader([]byte(`
		package main

		func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
			return color
		}
	`))
	if err != nil {
		panic(err)
	}

}

func (g *Game) Update() error {
	if len(foods) < 1 {
		food := models.NewFood(simpleShader, foodRadius, screenWidth, screenHeight)
		foods = append(foods, *food)
	}

	for foodIndex, food := range foods {
		if snake.IsHitFood(&food) {
			foods = append(foods[:foodIndex], foods[foodIndex+1:]...)
			g.score += 1
			snake.Grow()
			fmt.Print(snake.Bodies)
		}
	}

	g.pressedKeys = inpututil.AppendPressedKeys(g.pressedKeys[:0])

	for _, key := range g.pressedKeys {
		switch key.String() {
		case "S":
			if direction != "up" {
				shakeMovementPositionX = 0
				shakeMovementPositionY = headSize
				direction = "down"
			}
		case "W":
			if direction != "down" {
				shakeMovementPositionX = 0
				shakeMovementPositionY = -headSize
				direction = "up"
			}
		case "D":
			if direction != "left" {
				shakeMovementPositionX = headSize
				shakeMovementPositionY = 0
				direction = "right"
			}
		case "A":
			if direction != "right" {
				shakeMovementPositionX = -headSize
				shakeMovementPositionY = 0
				direction = "left"

			}
		}
	}
	snake.MoveBody()
	snake.X += int(math.Round(shakeMovementPositionX))
	snake.Y += int(math.Round(shakeMovementPositionY))

	time.Sleep(150 * time.Millisecond)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	snake.Draw(screen, snake.X, snake.Y,
		[]float32{0xf0 / float32(0xff), 0xc8 / float32(0xff), 0x00 / float32(0xff)})
	for _, body := range snake.Bodies {
		snake.Draw(screen, body.X, body.Y,
			[]float32{0x43 / float32(0xff), 0xff / float32(0xff), 0x64 / float32(0xff)})
	}

	for _, food := range foods {
		food.DrawFood(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	snake.SimpleShader = simpleShader
	snake.X = screenWidth / 2
	snake.Y = screenHeight / 2
	snake.Size = headSize
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Snake Game")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
