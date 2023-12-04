package models

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Snake struct {
	SimpleShader *ebiten.Shader
}

func (s *Snake) DrawHead(screen *ebiten.Image, x, y, size int, clr color.Color) {
	var path vector.Path

	// Draw square
	path.MoveTo(float32(x), float32(y))
	path.LineTo(float32(x+size), float32(y))
	path.LineTo(float32(x+size), float32(y+size))
	path.LineTo(float32(x), float32(y+size))
	path.Close()

	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)

	redScaled := 0x43 / float32(0xff)
	greenScaled := 0xff / float32(0xff)
	blueScaled := 0x64 / float32(0xff)
	alphaScaled := 0.85

	for i := range vertices {
		v := &vertices[i]
		v.ColorR = redScaled
		v.ColorG = greenScaled
		v.ColorB = blueScaled
		v.ColorA = float32(alphaScaled)
	}

	screen.DrawTrianglesShader(vertices, indices, s.SimpleShader, &ebiten.DrawTrianglesShaderOptions{
		FillRule: ebiten.EvenOdd,
	})
}
