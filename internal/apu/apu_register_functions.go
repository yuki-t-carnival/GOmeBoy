package apu

// =================================== Global control registers ====================================

func (a *APU) GetNR52() byte {
	return a.nr52
}

func (a *APU) SetNR52(val byte) {
	a.nr52 = val
}
func (a *APU) GetNR51() byte {
	return a.nr51
}

func (a *APU) SetNR51(val byte) {
	a.nr51 = val
}
func (a *APU) GetNR50() byte {
	return a.nr50
}

func (a *APU) SetNR50(val byte) {
	a.nr50 = val
}

// ====================================== Sound channel 1 registers ================================

func (a *APU) GetNR10() byte {
	return a.nr10
}
func (a *APU) SetNR10(val byte) {
	a.nr10 = val
}
func (a *APU) GetNR11() byte {
	return a.nr11
}
func (a *APU) SetNR11(val byte) {
	a.ch1LengthTimer = int(val & 0x3F)
	a.nr11 = val
}
func (a *APU) GetNR12() byte {
	return a.nr12
}
func (a *APU) SetNR12(val byte) {
	if val&0xF8 == 0 {
		a.nr52 &^= 1 << 0
	}
	a.nr12 = val
}
func (a *APU) GetNR13() byte {
	return a.nr13
}
func (a *APU) SetNR13(val byte) {
	a.ch1Period = (uint16(a.nr14&0x07) << 8) | uint16(val)
	a.nr13 = val
}
func (a *APU) GetNR14() byte {
	return a.nr14
}
func (a *APU) SetNR14(val byte) {
	a.ch1Period = (uint16(val&0x07) << 8) | uint16(a.nr13)
	if val&0x80 != 0 {
		a.nr52 |= (1 << 0)
		a.ch1Vol = a.nr12 >> 4
		a.ch1LengthTimer = int(a.nr11 & 0x3F)
		a.ch1SampleCountForLengthTimer = 0
		a.ch1SampleCountForEnvelope = 0
		a.ch1Phase = 0
	}
	a.nr14 = val
}

// ================================== Sound channel 2 registers ====================================

func (a *APU) GetNR21() byte {
	return a.nr21
}
func (a *APU) SetNR21(val byte) {
	a.ch2LengthTimer = int(val & 0x3F)
	a.nr21 = val
}
func (a *APU) GetNR22() byte {
	return a.nr22
}
func (a *APU) SetNR22(val byte) {
	if val&0xF8 == 0 {
		a.nr52 &^= 1 << 1
	}
	a.nr22 = val
}
func (a *APU) GetNR23() byte {
	return a.nr23
}
func (a *APU) SetNR23(val byte) {
	a.ch2Period = (uint16(a.nr24&0x07) << 8) | uint16(val)
	a.nr23 = val
}
func (a *APU) GetNR24() byte {
	return a.nr24
}
func (a *APU) SetNR24(val byte) {
	a.ch2Period = (uint16(val&0x07) << 8) | uint16(a.nr23)
	if val&0x80 != 0 {
		a.nr52 |= (1 << 1)
		a.ch2Vol = a.nr22 >> 4
		a.ch2LengthTimer = int(a.nr21 & 0x3F)
		a.ch2SampleCountForLengthTimer = 0
		a.ch2SampleCountForEnvelope = 0
		a.ch2Phase = 0
	}
	a.nr24 = val
}

// ================================== Sound channel 3 registers ====================================

func (a *APU) GetNR30() byte {
	return a.nr30
}
func (a *APU) SetNR30(val byte) {
	a.nr30 = val
}
func (a *APU) GetNR31() byte {
	return a.nr31
}
func (a *APU) SetNR31(val byte) {
	a.ch3LengthTimer = int(val)
	a.nr31 = val
}
func (a *APU) GetNR32() byte {
	return a.nr32
}
func (a *APU) SetNR32(val byte) {
	a.nr32 = val
}
func (a *APU) GetNR33() byte {
	return a.nr33
}
func (a *APU) SetNR33(val byte) {
	a.nr33 = val
}
func (a *APU) GetNR34() byte {
	return a.nr34
}
func (a *APU) SetNR34(val byte) {
	if val&0x80 != 0 {
		a.nr52 |= (1 << 2)
		a.ch3LengthTimer = int(a.nr31)
		a.ch3SampleCountForLengthTimer = 0
		a.indexWaveRAM = 0
		a.ch3Phase = 0
	}
	a.nr34 = val
}

// ================================== Sound channel 4 registers ====================================

func (a *APU) GetNR41() byte {
	return a.nr41
}
func (a *APU) SetNR41(val byte) {
	a.ch4LengthTimer = int(val & 0x3F)
	a.nr41 = val
}
func (a *APU) GetNR42() byte {
	return a.nr42
}
func (a *APU) SetNR42(val byte) {
	if val&0xF8 == 0 {
		a.nr52 &^= 1 << 3
	}
	a.nr42 = val
}
func (a *APU) GetNR43() byte {
	return a.nr43
}
func (a *APU) SetNR43(val byte) {
	a.nr43 = val
}
func (a *APU) GetNR44() byte {
	return a.nr44
}
func (a *APU) SetNR44(val byte) {
	if val&0x80 != 0 {
		a.nr52 |= (1 << 3)
		a.ch4Vol = a.nr42 >> 4
		a.ch4LengthTimer = int(a.nr41 & 0x3F)
		a.ch4SampleCountForLengthTimer = 0
		a.ch4SampleCountForEnvelope = 0

		a.lfsr = 0x7FFF
		a.ch4SampleCountForLFSR = 0
	}
	a.nr44 = val
}
