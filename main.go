package main

import (
	"image/color"
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
	shakeHeadPositionX     = float64(screenWidth) / 2
	shakeHeadPositionY     = float64(screenHeight) / 2
	shakeMovementPositionX = float64(0)
	shakeMovementPositionY = float64(0)
	prevUpdateTime         = time.Now()
	direction              = ""

	simpleShader *ebiten.Shader
	snake        = new(models.Snake)
	foods        []models.Food
)

type Game struct {
	pressedKeys []ebiten.Key
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

	timeDelta := float64(time.Since(prevUpdateTime))

	g.pressedKeys = inpututil.AppendPressedKeys(g.pressedKeys[:0])

	for _, key := range g.pressedKeys {
		switch key.String() {
		case "ArrowDown":
			if direction != "up" {
				shakeMovementPositionX = 0
				shakeMovementPositionY = 0.0000001
				direction = "down"
			}
		case "ArrowUp":
			if direction != "down" {
				shakeMovementPositionX = 0
				shakeMovementPositionY = -0.0000001
				direction = "up"
			}
		case "ArrowRight":
			if direction != "left" {
				shakeMovementPositionX = 0.0000001
				shakeMovementPositionY = 0
				direction = "right"
			}
		case "ArrowLeft":
			if direction != "right" {
				shakeMovementPositionX = -0.0000001
				shakeMovementPositionY = 0
				direction = "left"

			}
		}
	}

	shakeHeadPositionX += shakeMovementPositionX * timeDelta
	shakeHeadPositionY += shakeMovementPositionY * timeDelta

	prevUpdateTime = time.Now()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{60, 179, 113, 255})
	purpleCol := color.RGBA{255, 0, 255, 255}
	x := int(math.Round(shakeHeadPositionX))
	y := int(math.Round(shakeHeadPositionY))
	snake.DrawHead(screen, x, y, headSize, purpleCol)
	for _, food := range foods {
		food.DrawFood(screen, purpleCol)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	snake.SimpleShader = simpleShader
	// food.Radius = foodRadius
	// food.SimpleShader = simpleShader
	// food.SetCoordinate(screenWidth, screenHeight)
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Snake Game")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}