package emulator

import (
	"gomeboy/internal/bus"
	"gomeboy/internal/cpu"
	"gomeboy/internal/memory"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	MaxCycles = 70221 // 1フレームのクロック(cycle)数
)

type Emulator struct {
	CPU *cpu.CPU

	IsPaused     bool
	isPauseMode  bool // (実際に停止中かどうかは別｡pauseMode==trueでもステップ実行はできる｡)
	isKeyP       bool
	isKeyS       bool
	isKeyEsc     bool
	isPrevKeyP   bool
	isPrevKeyS   bool
	isPrevKeyEsc bool
}

func NewEmulator(rom []byte) *Emulator {
	m := memory.NewMemory(rom)
	b := bus.NewBus(m)
	c := cpu.NewCPU(b)
	c.Tracer = cpu.NewTracer(c)

	e := &Emulator{
		CPU:         c,
		isPauseMode: false,
		IsPaused:    false,
	}
	return e
}

// Gameboyを1フレーム実行する
func (e *Emulator) RunFrame() int {
	cycles := 0
	for cycles < MaxCycles {
		e.updateEbitenKeys()
		e.updateEmuState()

		if e.IsPaused {
			return 0
		}

		c := e.CPU.Step()
		e.CPU.Tracer.Record(e.CPU)
		if e.CPU.IsPanic || e.isKeyEsc {
			e.panicDump()
			return -1
		}

		e.CPU.Bus.Timer.Step(c, e.CPU.IsStopped)
		e.CPU.Bus.PPU.Step(c)
		cycles += c
	}
	e.CPU.Bus.Joypad.Update()
	return 0
}

// ゲーム状態を更新（PキーでPAUSE切替､PAUSE中Sキーでステップ実行）
func (e *Emulator) updateEmuState() {
	if e.isKeyP {
		e.isPauseMode = !e.isPauseMode
	}
	e.IsPaused = e.isPauseMode && !e.isKeyS
}

// パニックで終了するときは､直前のCPUの状態をコンソールに出力する
func (e *Emulator) panicDump() {
	e.CPU.Tracer.Dump()
}

func (e *Emulator) updateEbitenKeys() {
	isP := ebiten.IsKeyPressed(ebiten.KeyP)
	isS := ebiten.IsKeyPressed(ebiten.KeyS)
	isEsc := ebiten.IsKeyPressed(ebiten.KeyEscape)
	e.isKeyP = !e.isPrevKeyP && isP
	e.isKeyS = !e.isPrevKeyS && isS
	e.isKeyEsc = !e.isPrevKeyEsc && isEsc
	e.isPrevKeyP = isP
	e.isPrevKeyS = isS
	e.isPrevKeyEsc = isEsc
}
