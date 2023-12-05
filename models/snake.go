package models

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Snake struct {
	SimpleShader *ebiten.Shader
	X            int
	Y            int
	HeadSize     int
}

func (s *Snake) SetHeadPosition(x, y int) {
	s.X = x
	s.Y = y
}

func (s *Snake) DrawHead(screen *ebiten.Image, clr color.Color) {
	var path vector.Path

	// Draw square
	path.MoveTo(float32(s.X-s.HeadSize/2), float32(s.Y-s.HeadSize/2))
	path.LineTo(float32(s.X+s.HeadSize/2), float32(s.Y-s.HeadSize/2))
	path.LineTo(float32(s.X+s.HeadSize/2), float32(s.Y+s.HeadSize/2))
	path.LineTo(float32(s.X-s.HeadSize/2), float32(s.Y+s.HeadSize/2))
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

func (s *Snake) IsHitFood(food *Food) bool {
	if (s.X-food.x <= s.HeadSize/2+food.Radius && s.X-food.x >= 0 || food.x-s.X <= s.HeadSize/2+food.Radius && food.x-s.X >= 0) && (s.Y-food.y <= s.HeadSize/2+food.Radius && s.Y-food.y >= 0 || food.y-s.Y <= s.HeadSize/2+food.Radius && food.y-s.Y >= 0) {
		return true
	}
	return false
}
