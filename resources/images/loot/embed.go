package loot

import (
	_ "embed"
)

var (
	//go:embed shadowed_loot.png
	shadowed_loot []byte

	//go:embed transparent_loot.png
	transparent_loot []byte

	//go:embed shadowed_loot_v2.png
	shadowed_loot_v2 []byte

	//go:embed transparent_loot_v2.png
	transparent_loot_v2 []byte
)
