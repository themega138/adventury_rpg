package characters

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Orientation int

const (
	Down Orientation = iota
	Left
	Right
	Up
)

type Direction int

const (
	N Direction = iota
	NE
	E
	SE
	S
	SW
	W
	NW
)

type ICharacter interface {
	Draw(screen *ebiten.Image)
	Update() error
	Move(direction Direction)
}

type BaseCharacter struct {
	x           float64
	y           float64
	vx          float64
	vy          float64
	moving      bool
	steps       int
	orientation Orientation
}
