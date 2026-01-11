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
	isHaltBug    bool
	prevIF       byte
}

func NewCPU(b *bus.Bus) *CPU {
	c := &CPU{
		Bus: b,
		pc:  0x100,
		sp:  0xFFFE,
		a:   0x01,
		f:   0xB0,
		b:   0x00,
		c:   0x13,
		d:   0x00,
		e:   0xD8,
		h:   0x01,
		l:   0x4D,
	}
	return c
}

func (c *CPU) Step() int {
	c.cycles = 0

	c.checkIRQ()

	// During DMA transfer, only check for interrupts
	if c.Bus.IsDMATransferInProgress {
		c.Bus.DMATransfer()
		return 4
	}

	// STOP is not implemented
	// Stop mode ends when any input is received
	/* if c.IsStopped {
		if c.Bus.Joypad.HasStateChanged {
			c.IsStopped = false
		} else {
			return 4 // When set to 0, g.Update() does not finish
		}
	} */

	// HALT ends when any bit is set in IF
	if c.isHalted {
		if (c.read(bus.IF) & 0x1F) != 0 {
			c.isHalted = false
		} else {
			c.cycles += 4
			return c.cycles
		}
	}

	if c.handleInterrupt() {
		return c.cycles
	}

	op := c.fetchOpcode()
	OpTable[op].fn(c)

	// For EI
	if c.imeDelay > 0 {
		c.imeDelay--
		if c.imeDelay == 0 {
			c.isIMEEnabled = true
		}
	}

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
	c.pc++
	return op
	// TODO: Implement the HALT bug
	/* if c.isHaltBug {
		c.isHaltBug = false
	} else {
		c.pc++
	} */
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

func (c *CPU) handleInterrupt() bool {
	if !c.isIMEEnabled {
		return false
	}

	// allowed interrupts == IME & IEbit & IFbit
	curIE := c.read(bus.IE) & 0x1F
	curIF := c.read(bus.IF) & 0x1F
	pending := curIE & curIF
	if pending == 0 {
		return false
	}

	// PC set to $40, $48, $50, $58, $60, if interrupts are enabled
	for i := 0; i < 5; i++ {
		if (pending & (1 << i)) != 0 {
			c.write(bus.IF, curIF & ^(1<<i)) // Clear IF bit before the interrupt
			c.isIMEEnabled = false           // Disable IME before the interrupt
			c.push(c.pc)
			c.pc = 0x40 + 0x08*uint16(i)
			c.cycles += 20
			return true
		}
	}
	return false
}

func (c *CPU) checkIRQ() {
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
