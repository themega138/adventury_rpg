package scenes

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jakecoffman/cp"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	"github.com/themega138/adventury/core/characters"
	"github.com/themega138/adventury/core/commons"
	"image/color"
	_ "image/png"
	"os"
)

const (
	tileSize               = 16
	mapPath                = "resources/maps/map1.tmx"
	INITIAL_LOCATION       = "INITIAL_LOCATION"
	LOCATIONS              = "LOCATIONS"
	TERRAIN                = "TERRAIN"
	TERRAIN_BORDER         = "TERRAIN_BORDER"
	COLLISION_TYPE_ONE_WAY = 1
)

var (
	tilesImage []*ebiten.Image
	dot        = ebiten.NewImage(1, 1)
)

var GRABBABLE_MASK_BIT uint = 1 << 31

var NotGrabbableFilter cp.ShapeFilter = cp.ShapeFilter{
	cp.NO_GROUP, ^GRABBABLE_MASK_BIT, ^GRABBABLE_MASK_BIT,
}

var space *cp.Space
var playerBody *cp.Body
var playerShape *cp.Shape
var hero *characters.Hero

type GameScene struct {
	Hero *characters.Hero
}

func NewGameScene() *GameScene {
	// Parse .tmx file.
	gameMap, err := tiled.LoadFile(mapPath)
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}

	tilesImage = []*ebiten.Image{}

	// You can also render the map to an in-memory image for direct
	// use with the default Renderer, or by making your own.
	renderer, err := render.NewRenderer(gameMap)
	if err != nil {
		fmt.Printf("map unsupported for rendering: %s", err.Error())
		os.Exit(2)
	}
	space = cp.NewSpace()
	space.Iterations = 1
	space.SetCollisionSlop(0.5)

	for _, group := range gameMap.ObjectGroups {
		println(group.ID)
		println(group.Name)
		println(group.Class)
		for _, object := range group.Objects {
			if (object.Type == TERRAIN_BORDER || group.Class == TERRAIN) && object.Visible {
				body := cp.NewStaticBody()
				body.SetPosition(cp.Vector{X: object.X, Y: object.Y})
				body.SetType(cp.BODY_STATIC)
				body.UserData = object

				bb := cp.BB{float64(0), float64(0), float64(object.Width), float64(object.Height)}

				shape := cp.NewBox2(space.StaticBody, bb, 3)
				shape.SetElasticity(1)
				shape.SetFriction(1)

				space.AddBody(body)
				space.AddShape(shape)
			} else if object.Type == INITIAL_LOCATION {
				if playerBody == nil {
					playerBody = cp.NewBody(1.0, cp.INFINITY)
					playerBody.SetPosition(cp.Vector{X: object.X, Y: object.Y})
					playerBody.UserData = object

					bb := cp.BB{float64(0), float64(0), float64(16), float64(16)}

					playerShape = cp.NewBox2(playerBody, bb, 0)
					playerShape.SetElasticity(1)
					playerShape.SetFriction(1)
				}
			}
		}
	}

	for _, layer := range gameMap.Layers {
		// Render just item 0 to the Renderer.
		err = renderer.RenderLayer(int(layer.ID - 1))
		if err != nil {
			fmt.Printf("item unsupported for rendering: %s", err.Error())
			os.Exit(2)
		}

		// Get a reference to the Renderer's output, an image.NRGBA struct.
		img := renderer.Result

		// Clear the render result after copying the output if separation of
		// layers is desired.
		renderer.Clear()

		// And so on. You can also export the image to a file by using the
		// Renderer's Save functions.
		tilesImage = append(tilesImage, ebiten.NewImageFromImage(img))
	}

	space.AddBody(playerBody)
	space.AddShape(playerShape)
	hero = &characters.Hero{
		CharacterBody: playerBody,
		Space:         space,
	}

	return &GameScene{
		Hero: hero,
	}
}

func (g GameScene) Update(state *commons.GameState) error {
	space.Step(1.0 / float64(ebiten.TPS()))
	if state.Input.StateForLeft() >= 1 && state.Input.StateForUp() >= 1 {
		g.Hero.Move(characters.NW)
	} else if state.Input.StateForRight() >= 1 && state.Input.StateForUp() >= 1 {
		g.Hero.Move(characters.NE)
	} else if state.Input.StateForLeft() >= 1 && state.Input.StateForDown() >= 1 {
		g.Hero.Move(characters.SW)
	} else if state.Input.StateForRight() >= 1 && state.Input.StateForDown() >= 1 {
		g.Hero.Move(characters.SE)
	} else if state.Input.StateForLeft() >= 1 {
		g.Hero.Move(characters.W)
	} else if state.Input.StateForRight() >= 1 {
		g.Hero.Move(characters.E)
	} else if state.Input.StateForDown() >= 1 {
		g.Hero.Move(characters.S)
	} else if state.Input.StateForUp() >= 1 {
		g.Hero.Move(characters.N)
	} else {
		g.Hero.Stand()
	}
	g.Hero.Update()
	return nil
}

func (g GameScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Hx: %0.2f Hy: %0.2f", g.Hero.CharacterBody.Position().X, g.Hero.CharacterBody.Position().Y))
	drawBG(screen)
	g.Hero.Draw(screen)
	//drawCollisionLines(screen, g)
}

func drawCollisionLines(screen *ebiten.Image, g GameScene) {
	g.Hero.Space.EachBody(func(body *cp.Body) {
		body.EachShape(func(shape *cp.Shape) {
			// return fmt.Sprintf("%v %v %v %v", bb.L, bb.T, bb.R, bb.B)
			//log.Print(shape.BB())
			bb := shape.BB()
			vector.StrokeLine(screen, float32(bb.L), float32(bb.B), float32(bb.L), float32(bb.T), 1, color.RGBA{0xff, 0xff, 0xff, 0xff}, true)
			vector.StrokeLine(screen, float32(bb.R), float32(bb.B), float32(bb.R), float32(bb.T), 1, color.RGBA{0xff, 0xff, 0xff, 0xff}, true)
			vector.StrokeLine(screen, float32(bb.R), float32(bb.B), float32(bb.L), float32(bb.B), 1, color.RGBA{0xff, 0xff, 0xff, 0xff}, true)
			vector.StrokeLine(screen, float32(bb.R), float32(bb.T), float32(bb.L), float32(bb.T), 1, color.RGBA{0xff, 0xff, 0xff, 0xff}, true)
		})
		if object := body.UserData; object != nil {
			vector.DrawFilledRect(screen, float32(body.Position().X), float32(body.Position().Y), float32(object.(*tiled.Object).Width), float32(object.(*tiled.Object).Height), color.RGBA{0x80, 0x80, 0x80, 0xc0}, true)
		}
	})
}

func drawBG(screen *ebiten.Image) {
	for _, image := range tilesImage {
		screen.DrawImage(image, nil)
	}
}
