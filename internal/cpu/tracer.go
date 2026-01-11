package cpu

import (
	"fmt"
)

const TraceSize = 256

type Tracer struct {
	buf   [TraceSize]TraceEntry
	index int // of the next buffer to be used
}

type TraceEntry struct {
	a, f, b, c, d, e, h, l byte
	sp, pc                 uint16
	op                     uint16
	opName                 string
}

func NewTracer(c *CPU) *Tracer {
	t := &Tracer{}
	t.Record(c)
	return t
}

// Save the current CPU Registers state in a ring buffer
func (t *Tracer) Record(c *CPU) {
	op := uint16(c.read(c.pc))
	var opName string
	if op == 0xCB {
		nextOp := c.read(c.pc + 1)
		opName = CBTable[nextOp].Name
		op = 0xCB00 | uint16(nextOp)
	} else {
		opName = OpTable[op].Name
	}
	t.buf[t.index] = TraceEntry{
		pc:     c.pc,
		a:      c.a,
		f:      c.f,
		b:      c.b,
		c:      c.c,
		d:      c.d,
		e:      c.e,
		h:      c.h,
		l:      c.l,
		sp:     c.sp,
		op:     op,
		opName: opName,
	}
	t.index = (t.index + 1) % TraceSize // ring buffer
}

// Output all CPU log buffers to console.
func (t *Tracer) Dump() {
	for i := range TraceSize {
		idx := (t.index + i) % TraceSize // Output from the oldest dump
		buf := t.buf[idx]
		fmt.Printf(
			"PC:%04X "+
				"A:%02X "+
				"F:%02X "+
				"BC:%02X%02X "+
				"DE:%02X%02X "+
				"HL:%02X%02X "+
				"SP:%04X "+
				"Op:%04X "+
				"OpName:%s\n",
			buf.pc,
			buf.a, buf.f,
			buf.b, buf.c,
			buf.d, buf.e,
			buf.h, buf.l,
			buf.sp,
			buf.op,
			buf.opName,
		)
	}
}

// Get CPU status strings for the debug screen
func (t *Tracer) GetCPUInfo() []string {
	var idx int
	if t.index == 0 {
		idx = (TraceSize - 1) % TraceSize
	} else {
		idx = (t.index - 1) % TraceSize
	}
	buf := t.buf[idx]
	var str []string
	str = append(str, fmt.Sprintf("PC:%04X", buf.pc))
	str = append(str, fmt.Sprintf("AF:%02X%02X", buf.a, buf.f))
	str = append(str, fmt.Sprintf("BC:%02X%02X", buf.b, buf.c))
	str = append(str, fmt.Sprintf("DE:%02X%02X", buf.d, buf.e))
	str = append(str, fmt.Sprintf("HL:%02X%02X", buf.h, buf.l))
	str = append(str, fmt.Sprintf("SP:%04X", buf.sp))
	str = append(str, fmt.Sprintf("Op:%04X", buf.op))
	str = append(str, fmt.Sprintf("OpName:%s", buf.opName))
	return str
}
