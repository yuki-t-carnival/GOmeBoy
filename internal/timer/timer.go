package timer

type Timer struct {
	// I/O Registers
	tma  byte
	tima byte
	tac  byte

	// Internal Counter
	divCounter uint16 // 外部からはDIVとして上位8bitのみ見える
	prevDIV    uint16 // 立ち下がりbit監視のため､過去DIVを保持

	// Others
	HasIRQ           bool
	isPrevCPUStopped bool
	overflowDelay    int
}

func NewTimer() *Timer {
	return &Timer{
		divCounter: 0x00, // 実装上0でOK（実機は0xAB相当）
		tima:       0x00,
		tma:        0x00,
		tac:        0xF8,
	}
}

func (t *Timer) Step(cycles int, isCPUStopped bool) {
	// テストが通らないので､cpuCyclesを1ずつ分解して実行
	for i := 0; i < cycles; i++ {

		// TIMA OverflowのDelay処理
		if t.overflowDelay > 0 {
			t.overflowDelay -= 1
			if t.overflowDelay <= 0 {
				t.tima = t.tma
				t.HasIRQ = true
			}
		}
		t.prevDIV = t.divCounter

		// STOPに入った直後､diVCounterがリセットされる
		if !t.isPrevCPUStopped && isCPUStopped {
			t.divCounter = 0
		} else {
			t.divCounter += uint16(1)
		}
		t.isPrevCPUStopped = isCPUStopped

		// 以下､TIMAカウントが有効の場合の処理
		// DIVの対象ビットが立ち下がることで､TIMAが加算される
		tac := t.tac
		if tac&(1<<2) != 0 { // divCounterの立ち下がり監視対象bit
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
			if prev == 1 && now == 0 { // 対象bitが立ち下がったか?
				if t.tima == 0xFF {
					t.tima = 0          // Overflow時､即tima=0にはなる｡
					t.overflowDelay = 4 // ただし､TIMA=TMAセットと､IRQは､4サイクル遅れる
				} else {
					t.tima++ // TIMAカウントする
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
