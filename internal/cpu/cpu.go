package cpu

import (
	"gomeboy/internal/bus"
)

type CPU struct {
	Tracer *Tracer
	Bus    *bus.Bus

	// Registers
	a, f, b, c, d, e, h, l byte
	sp, pc                 uint16

	// Others
	IsPanic      bool
	IsStopped    bool
	isHalted     bool
	isIMEEnabled bool
	imeDelay     int
	cycles       int
	intVectors   [5]uint16
	isHaltBug    bool
	prevIF       byte
}

func NewCPU(b *bus.Bus) *CPU {
	c := &CPU{
		Bus: b,
		pc:  0x100,
		sp:  0xFFFE,

		// IFでの割り込みで使用
		intVectors: [5]uint16{
			0x40,
			0x48,
			0x50,
			0x58,
			0x60,
		},
	}
	return c
}

func (c *CPU) Step() int {
	c.cycles = 0

	// STOP状態は、何かしらの入力があったとき解除される。
	if c.IsStopped {
		if c.Bus.Joypad.HasStateChanged {
			c.IsStopped = false
		} else {
			return 4 // ebitenのUpdate()を回すため
		}
	}

	// 各ハードのIRQに応じて､IFレジスタのビットを立てる
	c.updateIF()

	// IFに立ち上がったbitがあればHALT解除
	if c.isHalted {
		curIF := c.read(bus.IF) & 0x1F
		c.isHalted = (^c.prevIF & curIF & 0x1F) == 0
		c.prevIF = curIF
		return 4
	}

	// フェッチ・実行
	op := c.fetchOpcode()
	OpTable[op].fn(c)

	// EIでIME=trueにするときの､反映までの遅延処理
	if c.imeDelay > 0 {
		c.imeDelay--
		if c.imeDelay == 0 {
			c.isIMEEnabled = true
		}
	}

	// 割り込み処理
	c.handleInterrupt()

	// 費やしたcycle数を返す
	return c.cycles
}

func (c *CPU) GetBC() uint16 {
	return (uint16(c.b) << 8) | uint16(c.c)
}
func (c *CPU) SetBC(val uint16) {
	c.b = byte(val >> 8)
	c.c = byte(val & 0x00FF)
}
func (c *CPU) GetDE() uint16 {
	return (uint16(c.d) << 8) | uint16(c.e)
}
func (c *CPU) SetDE(val uint16) {
	c.d = byte(val >> 8)
	c.e = byte(val & 0x00FF)
}
func (c *CPU) GetHL() uint16 {
	return (uint16(c.h) << 8) | uint16(c.l)
}
func (c *CPU) SetHL(val uint16) {
	c.h = byte(val >> 8)
	c.l = byte(val & 0x00FF)
}
func (c *CPU) GetAF() uint16 {
	return (uint16(c.a) << 8) | uint16(c.f&0xF0)
}
func (c *CPU) SetAF(val uint16) {
	c.a = byte(val >> 8)
	c.f = byte(val & 0x00F0)
}

func (c *CPU) GetFlagZ() bool {
	return (c.f & 0x80) == 0x80
}

func (c *CPU) SetFlagZ(b bool) {
	if b {
		c.f |= 0x80
	} else {
		c.f &= (^byte(0x80))
	}
}

func (c *CPU) GetFlagN() bool {
	return (c.f & 0x40) == 0x40
}

func (c *CPU) SetFlagN(b bool) {
	if b {
		c.f |= 0x40
	} else {
		c.f &= (^byte(0x40))
	}
}

func (c *CPU) GetFlagH() bool {
	return (c.f & 0x20) == 0x20
}

func (c *CPU) SetFlagH(b bool) {
	if b {
		c.f |= 0x20
	} else {
		c.f &= (^byte(0x20))
	}
}

func (c *CPU) GetFlagC() bool {
	return (c.f & 0x10) == 0x10
}

func (c *CPU) SetFlagC(b bool) {
	if b {
		c.f |= 0x10
	} else {
		c.f &= (^byte(0x10))
	}
}

func (c *CPU) fetchOpcode() byte {
	op := c.read(c.pc)
	if c.isHaltBug {
		c.isHaltBug = false
	} else {
		c.pc++
	}
	return op
}

func (c *CPU) read(addr uint16) byte {
	return c.Bus.Read(addr)
}

func (c *CPU) write(addr uint16, val byte) {
	c.Bus.Write(addr, val)
}

func (c *CPU) fetch() byte {
	v := c.read(c.pc)
	c.pc++
	return v
}
func (c *CPU) fetch16() uint16 {
	lo := c.fetch()
	hi := c.fetch()
	return uint16(hi)<<8 | uint16(lo)
}

func (c *CPU) handleInterrupt() {
	if !c.isIMEEnabled {
		return
	}

	// 有効な割込があるかチェック
	curIE := c.read(bus.IE) & 0x1F
	curIF := c.read(bus.IF) & 0x1F
	pending := curIE & curIF
	if pending == 0 {
		return
	}

	// PCを該当する割り込みベクタにセット
	for i, v := range c.intVectors {
		if (pending & (1 << i)) != 0 {
			c.write(bus.IF, curIF & ^(1<<i)) // 割り込みに入る前に､IFを解除
			c.isIMEEnabled = false           // 他の割り込みを禁止する
			c.push(c.pc)
			c.pc = v
			c.cycles += 20
			break
		}
	}
}

func (c *CPU) updateIF() {
	if c.Bus.PPU.HasVBlankIRQ {
		newIF := c.read(bus.IF) | (1 << 0)
		c.write(bus.IF, newIF)
		c.Bus.PPU.HasVBlankIRQ = false
	}
	if c.Bus.PPU.HasSTATIRQ {
		newIF := c.read(bus.IF) | (1 << 1)
		c.write(bus.IF, newIF)
		c.Bus.PPU.HasSTATIRQ = false
	}
	if c.Bus.Timer.HasIRQ {
		newIF := c.read(bus.IF) | (1 << 2)
		c.write(bus.IF, newIF)
		c.Bus.Timer.HasIRQ = false
	}
	if c.Bus.Joypad.HasIRQ {
		newIF := c.read(bus.IF) | (1 << 4)
		c.write(bus.IF, newIF)
		c.Bus.Joypad.HasIRQ = false
	}
}
