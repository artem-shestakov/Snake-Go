package models

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Food struct {
	SimpleShader *ebiten.Shader
	x            float32
	y            float32
	Radius       int
}

func NewFood(simpleShader *ebiten.Shader, radius, screenWidth, screenHeight int) *Food {
	food := Food{
		SimpleShader: simpleShader,
		Radius:       radius,
	}
	food.setCoordinate(screenWidth, screenHeight)
	return &food
}

func (f *Food) setCoordinate(screenWidth, screenHeight int) {
	f.x = float32(rand.Intn(screenWidth - f.Radius))
	f.y = float32(rand.Intn(screenHeight - f.Radius))
}

func (f *Food) DrawFood(screen *ebiten.Image, clr color.Color) {
	var path vector.Path

	// Draw square
	path.MoveTo(float32(f.x), float32(f.y))
	path.Arc(float32(f.x), float32(f.y), float32(f.Radius), 0, 2*math.Pi, vector.Clockwise)

	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)

	redScaled := 0xff / float32(0xff)
	greenScaled := 0x00 / float32(0xff)
	blueScaled := 0x00 / float32(0xff)
	alphaScaled := 0.85

	for i := range vertices {
		v := &vertices[i]
		v.ColorR = redScaled
		v.ColorG = greenScaled
		v.ColorB = blueScaled
		v.ColorA = float32(alphaScaled)
	}

	screen.DrawTrianglesShader(vertices, indices, f.SimpleShader, &ebiten.DrawTrianglesShaderOptions{
		FillRule: ebiten.EvenOdd,
	})
}
