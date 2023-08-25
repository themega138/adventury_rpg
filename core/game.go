package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/themega138/adventury/core/commons"
	"github.com/themega138/adventury/core/controls"
	"github.com/themega138/adventury/core/scenes"
)

type Game struct {
	sceneManager *scenes.SceneManager
	input        controls.Input
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return commons.ScreenWidth, commons.ScreenHeight
}

func (g *Game) Update() error {
	if g.sceneManager == nil {
		g.sceneManager = &scenes.SceneManager{}
		g.sceneManager.GoTo(scenes.NewGameScene())
	}

	g.input.Update()
	if err := g.sceneManager.Update(&g.input); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneManager.Draw(screen)
}
