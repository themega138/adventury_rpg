package commons

import (
	"github.com/themega138/adventury/core/controls"
)

type Position struct {
	x int
	y int
}

type GameState struct {
	//SceneManager *scenes.SceneManager
	Input *controls.Input
}
