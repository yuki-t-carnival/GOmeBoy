package joypad

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Joypad struct {
	// I/O Registers
	sel byte // P1/JOYPの4,5bit

	// Ebiten Inputs
	keys        byte // キーボードの入力状態
	prevKeys    byte
	gamepad     byte // ゲームパッドの入力状態
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

	// ButtonsとDpad両方を対象として､新たに押されたキーが1つでもあれば、
	// Joypad割り込み要求と､STOP命令解除のフラグ立てをする
	isKeysChanged := j.prevKeys&^j.keys != 0
	isGamepadChanged := j.prevGamepad&^j.gamepad != 0
	if isKeysChanged || isGamepadChanged {
		j.HasIRQ = true
		j.HasStateChanged = true
	}
}

// P1/JOYPレジスタの読み取り時､同bit4, 5のselectに応じて､下位4bitをセットして返す
func (j *Joypad) GetP1JOYP() byte {
	isSelBtn := j.sel&(1<<5) == 0
	isSelDpad := j.sel&(1<<4) == 0

	buttons := (j.keys & j.gamepad) & 0x0F
	dpad := (j.keys & j.gamepad) >> 4

	// 両方非選択なら入力なし､両方選択なら両入力AND､片方選択ならその入力
	n := byte(0)
	switch {
	case !isSelBtn && !isSelDpad:
		n = 0x0F
	case isSelBtn && isSelDpad:
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

// P1/JOYPレジスタへの書き込みはbit4, 5のみ有効
func (j *Joypad) SetP1JOYP(val byte) {
	j.sel = val & 0x30
}

// キーボードの入力状態を1Byteに保持｡(ON=0, OFF=1)
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

// ゲームパッドの入力状態を1Byteに保持｡(ON=0, OFF=1)
func (j *Joypad) updateEbitenGamepadButtons() {
	id := ebiten.GamepadID(0)
	var inputs [8]bool
	for i, v := range j.gamepadBind {
		if j.isGamepadEnabled {
			inputs[i] = ebiten.IsGamepadButtonPressed(id, ebiten.GamepadButton(v))
		} else {
			inputs[i] = false // ゲームパッド無効のときは､入力がないことにする
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
