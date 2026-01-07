package cpu

import (
	"fmt"
)

const TraceSize = 256

type Tracer struct {
	buf   [TraceSize]TraceEntry
	index int // 次に使うべきindex
}

type TraceEntry struct {
	a, f, b, c, d, e, h, l byte
	sp, pc                 uint16
	op                     byte
}

func NewTracer(c *CPU) *Tracer {
	t := &Tracer{}
	t.Record(c)
	return t
}

// 現在のCPUレジスタの状態をリングバッファに保存する
func (t *Tracer) Record(c *CPU) {
	t.buf[t.index] = TraceEntry{
		pc: c.pc,
		a:  c.a,
		f:  c.f,
		b:  c.b,
		c:  c.c,
		d:  c.d,
		e:  c.e,
		h:  c.h,
		l:  c.l,
		sp: c.sp,
		op: c.read(c.pc),
	}
	t.index = (t.index + 1) % TraceSize // リングバッファ
}

// コンソールにCPUバッファをすべて出力する
func (t *Tracer) Dump() {
	for i := range TraceSize {
		idx := (t.index + i) % TraceSize // +1 は一番古いダンプから出力するため
		buf := t.buf[idx]
		fmt.Printf(
			"PC:%04X "+"OP:%02X "+
				"A:%02X "+"F:%02X "+
				"BC:%02X%02X "+"DE:%02X%02X "+
				"HL:%02X%02X "+"SP:%04X\n",
			buf.pc, buf.op,
			buf.a, buf.f,
			buf.b, buf.c,
			buf.d, buf.e,
			buf.h, buf.l,
			buf.sp,
		)
	}
}

// 画面表示用のCPUステータスを[]stringで取得
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
	str = append(str, fmt.Sprintf("OP:%04X", buf.op))
	str = append(str, fmt.Sprintf("A :%02X", buf.a))
	str = append(str, fmt.Sprintf("F :%02X", buf.f))
	str = append(str, fmt.Sprintf("BC:%02X%02X", buf.b, buf.c))
	str = append(str, fmt.Sprintf("DE:%02X%02X", buf.d, buf.e))
	str = append(str, fmt.Sprintf("HL:%02X%02X", buf.h, buf.l))
	str = append(str, fmt.Sprintf("SP:%04X", buf.sp))
	return str
}
