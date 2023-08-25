package controls

import "github.com/hajimehoshi/ebiten/v2"

type virtualGamepadButton int

const axisThreshold = 0.75

type axis struct {
	id       int
	positive bool
}

type gamepadConfig struct {
	gamepadID            ebiten.GamepadID
	gamepadIDInitialized bool

	current         virtualGamepadButton
	buttons         map[virtualGamepadButton]ebiten.GamepadButton
	axes            map[virtualGamepadButton]axis
	assignedButtons map[ebiten.GamepadButton]struct{}
	assignedAxes    map[axis]struct{}

	defaultAxesValues map[int]float64
}

func (c *gamepadConfig) IsGamepadIDInitialized() bool {
	return c.gamepadIDInitialized
}

func (c *gamepadConfig) initializeIfNeeded() {
	if !c.gamepadIDInitialized {
		panic("not reached")
	}

	if ebiten.IsStandardGamepadLayoutAvailable(c.gamepadID) {
		return
	}

	if c.buttons == nil {
		c.buttons = map[virtualGamepadButton]ebiten.GamepadButton{}
	}
	if c.axes == nil {
		c.axes = map[virtualGamepadButton]axis{}
	}
	if c.assignedButtons == nil {
		c.assignedButtons = map[ebiten.GamepadButton]struct{}{}
	}
	if c.assignedAxes == nil {
		c.assignedAxes = map[axis]struct{}{}
	}

	// Set default values.
	// It is assumed that all axes are not pressed here.
	//
	// These default values are used to detect if an axis is actually pressed.
	// For example, on PS4 controllers, L2/R2's axes valuse can be -1.0.
	if c.defaultAxesValues == nil {
		c.defaultAxesValues = map[int]float64{}
		na := ebiten.GamepadAxisCount(c.gamepadID)
		for a := 0; a < na; a++ {
			c.defaultAxesValues[a] = ebiten.GamepadAxisValue(c.gamepadID, a)
		}
	}
}

// IsButtonPressed reports whether the given virtual button b is pressed.
func (c *gamepadConfig) IsButtonPressed(b virtualGamepadButton) bool {
	if !c.gamepadIDInitialized {
		panic("not reached")
	}

	if ebiten.IsStandardGamepadLayoutAvailable(c.gamepadID) {
		if ebiten.IsStandardGamepadButtonPressed(c.gamepadID, b.StandardGamepadButton()) {
			return true
		}

		const threshold = 0
		switch b {
		case virtualGamepadButtonLeft:
			return ebiten.StandardGamepadAxisValue(c.gamepadID, ebiten.StandardGamepadAxisLeftStickHorizontal) < -threshold
		case virtualGamepadButtonRight:
			return ebiten.StandardGamepadAxisValue(c.gamepadID, ebiten.StandardGamepadAxisLeftStickHorizontal) > threshold
		case virtualGamepadButtonDown:
			return ebiten.StandardGamepadAxisValue(c.gamepadID, ebiten.StandardGamepadAxisLeftStickVertical) < -threshold
		case virtualGamepadButtonUp:
			return ebiten.StandardGamepadAxisValue(c.gamepadID, ebiten.StandardGamepadAxisLeftStickVertical) > threshold
		}
		return false
	}

	c.initializeIfNeeded()

	bb, ok := c.buttons[b]
	if ok {
		return ebiten.IsGamepadButtonPressed(c.gamepadID, bb)
	}

	a, ok := c.axes[b]
	if ok {
		v := ebiten.GamepadAxisValue(c.gamepadID, a.id)
		if a.positive {
			return axisThreshold <= v && v <= 1.0
		}
		return -1.0 <= v && v <= -axisThreshold
	}
	return false
}

func (v virtualGamepadButton) StandardGamepadButton() ebiten.StandardGamepadButton {
	switch v {
	case virtualGamepadButtonLeft:
		return ebiten.StandardGamepadButtonLeftLeft
	case virtualGamepadButtonRight:
		return ebiten.StandardGamepadButtonLeftRight
	case virtualGamepadButtonDown:
		return ebiten.StandardGamepadButtonLeftBottom
	case virtualGamepadButtonButtonA:
		return ebiten.StandardGamepadButtonRightBottom
	case virtualGamepadButtonButtonB:
		return ebiten.StandardGamepadButtonRightRight
	default:
		panic("not reached")
	}
}
