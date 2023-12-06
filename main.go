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
	ModeTitle Mode = iota
	ModeGame
	ModeGameOver

	screenWidth   = 600
	screenHeight  = 600
	headSize      = 40
	foodRadius    = 17
	fontSize      = 25
	titleFontSize = fontSize * 1.5
)

var (
	shakeMovementPositionX = 0
	shakeMovementPositionY = 0

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

func (g *Game) init() {
	// Reset snake
	snake.SimpleShader = simpleShader
	snake.X = screenWidth / 2
	snake.Y = screenHeight / 2
	snake.Size = headSize
	snake.Bodies = []models.Body{}
	snake.Direction = ""
	shakeMovementPositionX = 0
	shakeMovementPositionY = 0

	// Reset foods
	foods = []models.Food{}

	// Reset score
	g.score = 0
}

// Check if keys pressed and return what keys
func (g *Game) isKeyPressed() (bool, []ebiten.Key) {
	if len(inpututil.AppendPressedKeys(g.pressedKeys[:0])) > 0 {
		return true, inpututil.AppendPressedKeys(g.pressedKeys[:0])
	}
	return false, []ebiten.Key{}
}

func (g *Game) Update() error {
	switch g.mode {
	case ModeTitle:
		if ok, keys := g.isKeyPressed(); ok {
			for _, key := range keys {
				switch key.String() {
				case "ArrowDown",
					"ArrowUp",
					"ArrowRight",
					"ArrowLeft":
					g.mode = ModeGame
				}
			}
		}
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

		_, g.pressedKeys = g.isKeyPressed()
		for _, key := range g.pressedKeys {
			switch key.String() {
			case "ArrowDown":
				if snake.Direction != "up" {
					shakeMovementPositionX = 0
					shakeMovementPositionY = g.speed
					snake.Direction = "down"
				}
			case "ArrowUp":
				if snake.Direction != "down" {
					shakeMovementPositionX = 0
					shakeMovementPositionY = -g.speed
					snake.Direction = "up"
				}
			case "ArrowRight":
				if snake.Direction != "left" {
					shakeMovementPositionX = g.speed
					shakeMovementPositionY = 0
					snake.Direction = "right"
				}
			case "ArrowLeft":
				if snake.Direction != "right" {
					shakeMovementPositionX = -g.speed
					shakeMovementPositionY = 0
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
		if snake.BodyCollision() {
			g.mode = ModeGameOver
		}
	case ModeGameOver:
		if ok, _ := g.isKeyPressed(); ok {
			g.init()
			g.mode = ModeGame
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.mode != ModeTitle {
		snake.Draw(screen, snake.X, snake.Y,
			[]float32{0x43 / float32(0xff), 0xff / float32(0xff), 0x64 / float32(0xff)})
		for _, body := range snake.Bodies {
			snake.Draw(screen, body.X, body.Y,
				[]float32{0x43 / float32(0xff), 0xff / float32(0xff), 0x64 / float32(0xff)})
		}

		for _, food := range foods {
			food.DrawFood(screen)
		}
	}
	var titleTexts string
	var texts string
	switch g.mode {
	case ModeTitle:
		titleTexts = "SNAKE"
		texts = "\n\nUse arrow Up, Down, Left, Right\nTo start press any arrow"
	case ModeGameOver:
		titleTexts = "GAME OVER!"
		texts = "\n\nPress any key..."
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(screenWidth-20, 20)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacingInPixels = fontSize
	op.PrimaryAlign = text.AlignEnd
	text.Draw(screen, fmt.Sprintf("score %03d", g.score), &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   fontSize / 2,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(screenWidth/2, screenHeight/2-titleFontSize/2)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacingInPixels = titleFontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, titleTexts, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   titleFontSize,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(screenWidth/2, screenHeight/2-titleFontSize/2)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacingInPixels = fontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, texts, &text.GoTextFace{
		Source: arcadeFaceSource,
		Size:   fontSize / 2,
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
