package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/themega138/adventury/core"
	"log"
)

func main() {
	ebiten.SetWindowSize(1024, 800)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&core.Game{}); err != nil {
		log.Fatal(err)
	}
}
