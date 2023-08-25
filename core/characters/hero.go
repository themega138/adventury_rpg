package characters

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/themega138/adventury/resources/images/characters"
	"image"
	_ "image/png"
	"log"
	"math"
)

const (
	frameOX            = 0
	frameOY            = 0
	frameWidth         = 16
	frameHeight        = 16
	movementMultiplier = 60
	diagonalMovement   = 0.70710678118654752440084436210485 * movementMultiplier
	straightMovement   = 1.0 * movementMultiplier
	zero               = 0.0
)

var (
	heroImage       *ebiten.Image
	weaponsImage    *ebiten.Image
	shieldsImage    *ebiten.Image
	armoursImage    *ebiten.Image
	hairImage       *ebiten.Image
	weaponHandImage *ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(characters.Race_human_base_pale_png))
	if err != nil {
		log.Fatal(err)
	}
	heroImage = ebiten.NewImageFromImage(img)
	img, _, err = image.Decode(bytes.NewReader(characters.Weapon_1_near_png))
	if err != nil {
		log.Fatal(err)
	}
	weaponsImage = ebiten.NewImageFromImage(img)
	img, _, err = image.Decode(bytes.NewReader(characters.Shield_1_near_png))
	if err != nil {
		log.Fatal(err)
	}
	shieldsImage = ebiten.NewImageFromImage(img)
	img, _, err = image.Decode(bytes.NewReader(characters.Armours_png))
	if err != nil {
		log.Fatal(err)
	}
	armoursImage = ebiten.NewImageFromImage(img)
	img, _, err = image.Decode(bytes.NewReader(characters.Hair_1_png))
	if err != nil {
		log.Fatal(err)
	}
	hairImage = ebiten.NewImageFromImage(img)
	img, _, err = image.Decode(bytes.NewReader(characters.Race_human_weapon_hand_pale_png))
	if err != nil {
		log.Fatal(err)
	}
	weaponHandImage = ebiten.NewImageFromImage(img)

}

type Hero struct {
	BaseCharacter
	CharacterBody  *cp.Body
	CharacterShape *cp.Shape
	Space          *cp.Space
}

func (h *Hero) Move(direction Direction) {
	switch direction {
	case N:
		h.move(zero, -straightMovement, Up)
		break
	case NE:
		h.move(+diagonalMovement, -diagonalMovement, Right)
		break
	case E:
		h.move(+straightMovement, zero, Right)
		break
	case SE:
		h.move(+diagonalMovement, +diagonalMovement, Right)
		break
	case S:
		h.move(zero, +straightMovement, Down)
		break
	case SW:
		h.move(-diagonalMovement, +diagonalMovement, Left)
		break
	case W:
		h.move(-straightMovement, zero, Left)
		break
	case NW:
		h.move(-diagonalMovement, -diagonalMovement, Left)
		break
	default:
		log.Fatal("error...")
	}
	h.steps++
}

func (h *Hero) move(vx float64, vy float64, orientation Orientation) {
	h.moving = true
	h.steps++
	h.CharacterBody.SetVelocity(vx, vy)
	h.orientation = orientation
}

func (h *Hero) Stand() {
	h.moving = false
	h.CharacterBody.SetVelocity(0, 0)
}

func (h *Hero) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(h.CharacterBody.Position().X, h.CharacterBody.Position().Y)
	i := calculatePosition(h)
	sx, sy := frameOX+i*frameWidth, frameOY+int(h.orientation)*frameHeight
	screen.DrawImage(heroImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
	screen.DrawImage(hairImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
	screen.DrawImage(armoursImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
	screen.DrawImage(weaponsImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
	screen.DrawImage(shieldsImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
	screen.DrawImage(weaponHandImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

}

func calculatePosition(h *Hero) int {
	if h.moving {
		return int(math.Round(math.Sin(float64(h.steps/5)) + 1.0))
	} else {
		return 1
	}
}

func (h *Hero) Update() error {
	return nil
}
