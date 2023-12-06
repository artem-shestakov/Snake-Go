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
	Direction    string
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
	var body Body
	// body := Body{
	// 	X: s.X,
	// 	Y: s.Y,
	// }
	for i := 1; i <= 3; i++ {
		switch s.Direction {
		case "up":
			body = Body{
				X: s.X,
				Y: s.Y + i,
			}
		case "down":
			body = Body{
				X: s.X,
				Y: s.Y - i,
			}
		case "left":
			body = Body{
				X: s.X + i,
				Y: s.Y,
			}
		case "right":
			body = Body{
				X: s.X - i,
				Y: s.Y,
			}
		}
		s.Bodies = append(s.Bodies, body)
	}
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
}

func (s *Snake) IsHitFood(food *Food) bool {
	if (s.X-food.x <= s.Size/2+food.Radius && s.X-food.x >= 0 || food.x-s.X <= s.Size/2+food.Radius && food.x-s.X >= 0) && (s.Y-food.y <= s.Size/2+food.Radius && s.Y-food.y >= 0 || food.y-s.Y <= s.Size/2+food.Radius && food.y-s.Y >= 0) {
		return true
	}
	return false
}

func (s *Snake) BoardCollision(screenWidth, screenHeight int) bool {
	if s.X-s.Size/2 <= 0 || s.X+s.Size/2 >= screenWidth || s.Y-s.Size/2 <= 0 || s.Y+s.Size/2 >= screenHeight {
		return true
	}
	return false
}

func (s *Snake) BodyCollision() bool {
	// var headX1, headX2, headY1, headY2 int
	// switch s.Direction {
	// case "up":
	// 	headX1 = s.X - s.Size/2
	// 	headX2 = s.X + s.Size/2
	// 	headY1 = s.Y - s.Size/2
	// 	headY2 = headY1
	// case "down":
	// 	headX1 = s.X + s.Size/2
	// 	headX2 = s.X - s.Size/2
	// 	headY1 = s.Y + s.Size/2
	// 	headY2 = headY1
	// case "left":
	// 	headX1 = s.X - s.Size/2
	// 	headX2 = headX1
	// 	headY1 = s.Y + s.Size/2
	// 	headY2 = s.Y - s.Size/2
	// case "right":
	// 	headX1 = s.X + s.Size/2
	// 	headX2 = headX1
	// 	headY1 = s.Y - s.Size/2
	// 	headY2 = s.Y + s.Size/2
	// }
	for _, body := range s.Bodies {
		// if s.Direction == "up" && (body.X+s.Size/2 > headX1 || body.X-s.Size/2 < headX2) && (body.Y+s.Size/2 > headY1 || body.Y-s.Size/2 < headY2) {
		// 	return true
		// }
		if s.Direction == "up" && ((s.X-body.X <= s.Size && s.X-body.X >= 0) || (body.X-s.X <= s.Size && body.X-s.X >= 0)) && (s.Y-body.Y <= s.Size && s.Y-body.Y > 0) {
			return true
		}
	}
	return false
}
