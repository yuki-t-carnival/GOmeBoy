package timer

type Timer struct {
	// I/O Registers
	tma  byte
	tima byte
	tac  byte

	// Internal Counter
	divCounter uint16 // The upper 8 bits of divCounter == DIV register
	prevDIV    uint16

	// TIMA overflow delay
	isOverflow bool
	cycleCount int

	// Others
	HasIRQ           bool
	isPrevCPUStopped bool
}

func NewTimer() *Timer {
	return &Timer{
		divCounter: 0x00,
		tima:       0x00,
		tma:        0x00,
		tac:        0xF8,
	}
}

func (t *Timer) Step(cycles int, isCPUStopped bool) {
	// process every cycle to pass the test ROM
	for i := 0; i < cycles; i++ {

		// Processing after TIMA overflow is delayed by 4 cycles
		// and is executed here
		if t.isOverflow {
			t.cycleCount++
			if t.cycleCount >= 4 {
				t.tima = t.tma
				t.HasIRQ = true
				t.isOverflow = false
				t.cycleCount = 0
			}
		}

		// When STOP occurs, the divCounter is reset
		t.prevDIV = t.divCounter
		if !t.isPrevCPUStopped && isCPUStopped {
			t.divCounter = 0
		} else {
			t.divCounter += uint16(1)
		}
		t.isPrevCPUStopped = isCPUStopped

		// TIMA is incremented when the specified bit of DIV falls
		tac := t.tac
		if tac&(1<<2) != 0 { // Gets the target bit of divCounter
			var checkBit int
			switch tac & 0x03 {
			case 0b00:
				checkBit = 9
			case 0b01:
				checkBit = 3
			case 0b10:
				checkBit = 5
			case 0b11:
				checkBit = 7
			}
			prev := (t.prevDIV >> checkBit) & 1
			now := (t.divCounter >> checkBit) & 1
			if prev == 1 && now == 0 { // Did the target bit fall?
				if t.tima == 0xFF {
					t.tima = 0
					t.isOverflow = true // tima=tma and IRQ are executed after a delay
					t.cycleCount = 0
				} else {
					t.tima++
				}
			}
		}
	}
}

func (t *Timer) GetDIV() byte {
	return byte(t.divCounter >> 8)
}

func (t *Timer) ResetDiv() {
	t.prevDIV = t.divCounter
	t.divCounter = 0
}

func (t *Timer) GetTIMA() byte {
	return t.tima
}

func (t *Timer) SetTIMA(val byte) {
	t.tima = val
}

func (t *Timer) GetTMA() byte {
	return t.tma
}

func (t *Timer) SetTMA(val byte) {
	t.tma = val
}

func (t *Timer) GetTAC() byte {
	return 0xF8 | (t.tac & 0x07)
}

func (t *Timer) SetTAC(val byte) {
	t.tac = 0xF8 | (val & 0x07)
	t.prevDIV = t.divCounter
}
