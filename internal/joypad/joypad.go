package joypad

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Joypad struct {
	// I/O Registers
	sel byte // P1/JOYP.4,5

	// Ebiten Inputs
	keys        byte
	prevKeys    byte
	gamepad     byte
	prevGamepad byte

	// Others
	HasStateChanged  bool
	HasIRQ           bool
	isGamepadEnabled bool   // From config.toml
	gamepadBind      [8]int // From config.toml
}

func NewJoypad() *Joypad {
	return &Joypad{
		keys:     0xFF,
		prevKeys: 0xFF,
		sel:      0x00,
	}
}

func (j *Joypad) Update() {
	j.updateEbitenKeys()
	j.updateEbitenGamepadButtons()

	// If any joypad input is detected,
	// sets the IRQ and STOP cacel flags
	isKeysChanged := j.prevKeys&^j.keys != 0
	isGamepadChanged := j.prevGamepad&^j.gamepad != 0
	if isKeysChanged || isGamepadChanged {
		j.HasIRQ = true
		//j.HasStateChanged = true // For STOP cancellation
	}
}

// If select buttons/d-pad bit is 0,
// then buttons/directional keys set to the lower nibble
func (j *Joypad) GetP1JOYP() byte {
	isSelBtn := j.sel&(1<<5) == 0
	isSelDpad := j.sel&(1<<4) == 0

	buttons := (j.keys & j.gamepad) & 0x0F
	dpad := (j.keys & j.gamepad) >> 4

	n := byte(0)
	switch {
	case !isSelBtn && !isSelDpad: // Neither buttons nor d-pad is selected
		n = 0x0F
	case isSelBtn && isSelDpad: // Both buttons and d-pad are selected
		n = buttons & dpad
	default:
		if isSelBtn {
			n = buttons
		} else {
			n = dpad
		}
	}
	return 0xC0 | (j.sel & 0x30) | (n & 0x0F)
}

func (j *Joypad) SetP1JOYP(val byte) {
	j.sel = val & 0x30
}

// (Pressed=0, Released=1)
func (j *Joypad) updateEbitenKeys() {
	inputs := [8]bool{
		ebiten.IsKeyPressed(ebiten.KeyZ),         // A
		ebiten.IsKeyPressed(ebiten.KeyX),         // B
		ebiten.IsKeyPressed(ebiten.KeyShiftLeft), // SELECT
		ebiten.IsKeyPressed(ebiten.KeyEnter),     // START
		ebiten.IsKeyPressed(ebiten.KeyRight),     // RIGHT
		ebiten.IsKeyPressed(ebiten.KeyLeft),      // LEFT
		ebiten.IsKeyPressed(ebiten.KeyUp),        // UP
		ebiten.IsKeyPressed(ebiten.KeyDown),      // DOWN
	}
	j.prevKeys = j.keys
	j.keys = byte(0xFF)
	for i, b := range inputs {
		if b {
			j.keys &^= (1 << i)
		}
	}
}

// (Pressed=0, Released=1)
func (j *Joypad) updateEbitenGamepadButtons() {
	id := ebiten.GamepadID(0)
	var inputs [8]bool
	for i, v := range j.gamepadBind {
		if j.isGamepadEnabled {
			inputs[i] = ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton(v))
		} else {
			inputs[i] = false // When the gamepad is disabled, it is treated as no input
		}
	}
	j.prevGamepad = j.gamepad
	j.gamepad = byte(0xFF)
	for i, b := range inputs {
		if b {
			j.gamepad &^= (1 << i)
		}
	}
}

// From config.toml
func (j *Joypad) SetIsGamepadEnabled(b bool) {
	j.isGamepadEnabled = b
}

// From config.toml
func (j *Joypad) SetIsGamepadBind(bind [8]int) {
	j.gamepadBind = bind
}
