package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/artem-shestakov/Snake-Go/models"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Mode int

const (
	ModeGame Mode = iota
	ModeGameOver

	screenWidth   = 400
	screenHeight  = 400
	headSize      = 20
	foodRadius    = 10
	fontSize      = 24
	titleFontSize = fontSize * 1.5
)

var (
	shakeMovementPositionX = 0
	shakeMovementPositionY = 0
	direction              = ""

	simpleShader *ebiten.Shader
	snake        = new(models.Snake)
	foods        []models.Food

	arcadeFaceSource *text.GoTextFaceSource
)

type Game struct {
	mode        Mode
	pressedKeys []ebiten.Key
	score       int
	speed       int
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

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	arcadeFaceSource = s

}

func (g *Game) Update() error {
	switch g.mode {
	case ModeGame:
		if len(foods) < 1 {
			food := models.NewFood(simpleShader, foodRadius, screenWidth, screenHeight)
			foods = append(foods, *food)
		}

		for foodIndex, food := range foods {
			if snake.IsHitFood(&food) {
				foods = append(foods[:foodIndex], foods[foodIndex+1:]...)
				g.score += 1
				snake.Grow()
			}
		}

		g.pressedKeys = inpututil.AppendPressedKeys(g.pressedKeys[:0])
		for _, key := range g.pressedKeys {
			switch key.String() {
			case "ArrowDown":
				if direction != "up" {
					shakeMovementPositionX = 0
					shakeMovementPositionY = g.speed
					direction = "down"
					snake.Direction = "down"
				}
			case "ArrowUp":
				if direction != "down" {
					shakeMovementPositionX = 0
					shakeMovementPositionY = -g.speed
					direction = "up"
					snake.Direction = "up"
				}
			case "ArrowRight":
				if direction != "left" {
					shakeMovementPositionX = g.speed
					shakeMovementPositionY = 0
					direction = "right"
					snake.Direction = "right"
				}
			case "ArrowLeft":
				if direction != "right" {
					shakeMovementPositionX = -g.speed
					shakeMovementPositionY = 0
					direction = "left"
					snake.Direction = "left"
				}
			}
		}
		snake.MoveBody()
		snake.X += shakeMovementPositionX
		snake.Y += shakeMovementPositionY
		if snake.BoardCollision(screenWidth, screenHeight) {
			g.mode = ModeGameOver
		}
	case ModeGameOver:

	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	snake.Draw(screen, snake.X, snake.Y,
		[]float32{0x43 / float32(0xff), 0xff / float32(0xff), 0x64 / float32(0xff)})
	for _, body := range snake.Bodies {
		snake.Draw(screen, body.X, body.Y,
			[]float32{0x43 / float32(0xff), 0xff / float32(0xff), 0x64 / float32(0xff)})
	}

	for _, food := range foods {
		food.DrawFood(screen)
	}
	// var titleTexts string
	var texts string
	switch g.mode {
	case ModeGameOver:
		texts = "\nGAME OVER!"
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(screenWidth, 0)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacingInPixels = fontSize
	op.PrimaryAlign = text.AlignEnd
	text.Draw(screen, fmt.Sprintf("score %03d", g.score), &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   fontSize / 2,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(screenWidth/2, 3*titleFontSize)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacingInPixels = fontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, texts, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   fontSize,
	}, op)

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
	ebiten.SetWindowTitle("Snake Game in Golang")
	if err := ebiten.RunGame(&Game{
		speed: 5,
	}); err != nil {
		log.Fatal(err)
	}
}
