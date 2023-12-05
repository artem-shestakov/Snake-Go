package models

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Snake struct {
	SimpleShader *ebiten.Shader
	X            int
	Y            int
	Size         int
	Bodies       []Body
}

type Body struct {
	X int
	Y int
}

func (s *Snake) SetHeadPosition(x, y int) {
	s.X = x
	s.Y = y
}

func (s *Snake) Draw(screen *ebiten.Image, x, y int, clr []float32) {
	var path vector.Path
	// Draw square
	path.MoveTo(float32(x-s.Size/2), float32(y-s.Size/2))
	path.LineTo(float32(x+s.Size/2), float32(y-s.Size/2))
	path.LineTo(float32(x+s.Size/2), float32(y+s.Size/2))
	path.LineTo(float32(x-s.Size/2), float32(y+s.Size/2))
	path.Close()

	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)

	redScaled := clr[0]
	greenScaled := clr[1]
	blueScaled := clr[2]
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

func (s *Snake) Grow() {
	body := Body{
		X: s.X,
		Y: s.Y,
	}
	s.Bodies = append(s.Bodies, body)
}

func (s *Snake) MoveBody() {
	for i := len(s.Bodies) - 1; i >= 1; i-- {
		x := s.Bodies[i-1].X
		y := s.Bodies[i-1].Y
		s.Bodies[i].X = x
		s.Bodies[i].Y = y
	}
	if len(s.Bodies) > 0 {
		s.Bodies[0].X = s.X
		s.Bodies[0].Y = s.Y
	}
	// if len(s.Bodies) > 0 {
	// 	tmp := s.Bodies
	// 	s.Bodies = []Body{{X: s.X, Y: s.Y}}
	// 	s.Bodies = append(s.Bodies, tmp[:1]...)
	// 	// s.Bodies = s.Bodies[:len(s.Bodies)]
	// 	// s.Bodies = append(s.Bodies, Body{X: s.X, Y})
	// }
}

func (s *Snake) IsHitFood(food *Food) bool {
	if (s.X-food.x <= s.Size/2+food.Radius && s.X-food.x >= 0 || food.x-s.X <= s.Size/2+food.Radius && food.x-s.X >= 0) && (s.Y-food.y <= s.Size/2+food.Radius && s.Y-food.y >= 0 || food.y-s.Y <= s.Size/2+food.Radius && food.y-s.Y >= 0) {
		return true
	}
	return false
}
