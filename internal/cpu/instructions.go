package cpu

import (
	"gomeboy/internal/bus"
)

// -------------------------------- NOP ----------------------------------------
func (c *CPU) opNOP() { // 00
	c.cycles += 4
}

// --------------------------------- DI ----------------------------------------
func (c *CPU) opDI() { // F3
	c.isIMEEnabled = false
	c.cycles += 4
}

// --------------------------------- EI ----------------------------------------
func (c *CPU) opEI() { // FB
	c.imeDelay = 2
	c.cycles += 4
}

// --------------------------------- STOP --------------------------------------
func (c *CPU) opSTOP_n8() { // 10
	_ = c.fetch()
	//c.IsStopped = true // TODO: Implement the STOP
	c.cycles += 4
}

// --------------------------------- HALT --------------------------------------
func (c *CPU) opHALT() { // 76
	ieReg := c.read(bus.IE) & 0x1F
	ifReg := c.read(bus.IF) & 0x1F
	pending := ieReg & ifReg

	if !c.isIMEEnabled && (pending != 0) {
		//c.isHaltBug = true // TODO: Implement the HALT bug
		c.isHalted = false
	} else {
		c.isHalted = true
	}
	c.cycles += 4
}

// --------------------------------- PREFIX ------------------------------------
func (c *CPU) opPREFIX() { // CB
	op := c.fetch()
	CBTable[op].fn(c)
	// c.cycles += 4 // Added by each CB prefixed instrs
}

// ----------------------------------- LD r16, n16 -----------------------------
func (c *CPU) opLD_BC_n16() { // 01
	c.SetBC(c.fetch16())
	c.cycles += 12
}
func (c *CPU) opLD_DE_n16() { // 11
	c.SetDE(c.fetch16())
	c.cycles += 12
}
func (c *CPU) opLD_HL_n16() { // 21
	c.SetHL(c.fetch16())
	c.cycles += 12
}
func (c *CPU) opLD_SP_n16() { // 31
	c.sp = c.fetch16()
	c.cycles += 12
}

// ------------------------------------ LD [r16], A ----------------------------
func (c *CPU) opLD_aBC_A() { // 02
	c.write(c.GetBC(), c.a)
	c.cycles += 8
}
func (c *CPU) opLD_aDE_A() { // 12
	c.write(c.GetDE(), c.a)
	c.cycles += 8
}
func (c *CPU) opLD_aHLp_A() { // 22
	c.write(c.GetHL(), c.a)
	c.SetHL(c.GetHL() + 1)
	c.cycles += 8
}
func (c *CPU) opLD_aHLm_A() { // 32
	c.write(c.GetHL(), c.a)
	c.SetHL(c.GetHL() - 1)
	c.cycles += 8
}

// -------------------------------- LD r8, n8 ----------------------------------
func (c *CPU) ld_r_n8(dst *byte) {
	val := c.fetch()
	*dst = val
	c.cycles += 8
}
func (c *CPU) opLD_B_n8() { c.ld_r_n8(&c.b) } // 06
func (c *CPU) opLD_D_n8() { c.ld_r_n8(&c.d) } // 16
func (c *CPU) opLD_H_n8() { c.ld_r_n8(&c.h) } // 26
func (c *CPU) opLD_C_n8() { c.ld_r_n8(&c.c) } // 0E
func (c *CPU) opLD_E_n8() { c.ld_r_n8(&c.e) } // 1E
func (c *CPU) opLD_L_n8() { c.ld_r_n8(&c.l) } // 2E
func (c *CPU) opLD_A_n8() { c.ld_r_n8(&c.a) } // 3E

// --------------------------------- LD [HL], n8 -------------------------------
func (c *CPU) opLD_aHL_n8() { // 36
	val := c.fetch()
	c.write(c.GetHL(), val)
	c.cycles += 12
}

// ---------------------------------- LD [a16], SP -----------------------------
func (c *CPU) opLD_aa16_SP() { // 08
	addr := c.fetch16()
	lo := byte(c.sp & 0x00FF)
	hi := byte((c.sp & 0xFF00) >> 8)
	c.write(addr, lo)
	c.write(addr+1, hi)
	c.cycles += 20
}

// ---------------------------------- LD A, [r16] ------------------------------
func (c *CPU) opLD_A_aBC() { // 0A
	c.a = c.read(c.GetBC())
	c.cycles += 8
}
func (c *CPU) opLD_A_aDE() { // 1A
	c.a = c.read(c.GetDE())
	c.cycles += 8
}
func (c *CPU) opLD_A_aHLp() { // 2A
	c.a = c.read(c.GetHL())
	c.SetHL(c.GetHL() + 1)
	c.cycles += 8
}
func (c *CPU) opLD_A_aHLm() { // 3A
	c.a = c.read(c.GetHL())
	c.SetHL(c.GetHL() - 1)
	c.cycles += 8
}

// ------------------------------ LD r8, r8 ------------------------------------
func (c *CPU) ld_r_r(dst *byte, src byte) {
	*dst = src
	c.cycles += 4
}
func (c *CPU) opLD_B_B() { c.ld_r_r(&c.b, c.b) } // 40
func (c *CPU) opLD_D_B() { c.ld_r_r(&c.d, c.b) } // 50
func (c *CPU) opLD_H_B() { c.ld_r_r(&c.h, c.b) } // 60
func (c *CPU) opLD_B_C() { c.ld_r_r(&c.b, c.c) } // 41
func (c *CPU) opLD_D_C() { c.ld_r_r(&c.d, c.c) } // 51
func (c *CPU) opLD_H_C() { c.ld_r_r(&c.h, c.c) } // 61
func (c *CPU) opLD_B_D() { c.ld_r_r(&c.b, c.d) } // 42
func (c *CPU) opLD_D_D() { c.ld_r_r(&c.d, c.d) } // 52
func (c *CPU) opLD_H_D() { c.ld_r_r(&c.h, c.d) } // 62
func (c *CPU) opLD_B_E() { c.ld_r_r(&c.b, c.e) } // 43
func (c *CPU) opLD_D_E() { c.ld_r_r(&c.d, c.e) } // 53
func (c *CPU) opLD_H_E() { c.ld_r_r(&c.h, c.e) } // 63
func (c *CPU) opLD_B_H() { c.ld_r_r(&c.b, c.h) } // 44
func (c *CPU) opLD_D_H() { c.ld_r_r(&c.d, c.h) } // 54
func (c *CPU) opLD_H_H() { c.ld_r_r(&c.h, c.h) } // 64
func (c *CPU) opLD_B_L() { c.ld_r_r(&c.b, c.l) } // 45
func (c *CPU) opLD_D_L() { c.ld_r_r(&c.d, c.l) } // 55
func (c *CPU) opLD_H_L() { c.ld_r_r(&c.h, c.l) } // 65
func (c *CPU) opLD_B_A() { c.ld_r_r(&c.b, c.a) } // 47
func (c *CPU) opLD_D_A() { c.ld_r_r(&c.d, c.a) } // 57
func (c *CPU) opLD_H_A() { c.ld_r_r(&c.h, c.a) } // 67
func (c *CPU) opLD_C_B() { c.ld_r_r(&c.c, c.b) } // 48
func (c *CPU) opLD_E_B() { c.ld_r_r(&c.e, c.b) } // 58
func (c *CPU) opLD_L_B() { c.ld_r_r(&c.l, c.b) } // 68
func (c *CPU) opLD_A_B() { c.ld_r_r(&c.a, c.b) } // 78
func (c *CPU) opLD_C_C() { c.ld_r_r(&c.c, c.c) } // 49
func (c *CPU) opLD_E_C() { c.ld_r_r(&c.e, c.c) } // 59
func (c *CPU) opLD_L_C() { c.ld_r_r(&c.l, c.c) } // 69
func (c *CPU) opLD_A_C() { c.ld_r_r(&c.a, c.c) } // 79
func (c *CPU) opLD_C_D() { c.ld_r_r(&c.c, c.d) } // 4A
func (c *CPU) opLD_E_D() { c.ld_r_r(&c.e, c.d) } // 5A
func (c *CPU) opLD_L_D() { c.ld_r_r(&c.l, c.d) } // 6A
func (c *CPU) opLD_A_D() { c.ld_r_r(&c.a, c.d) } // 7A
func (c *CPU) opLD_C_E() { c.ld_r_r(&c.c, c.e) } // 4B
func (c *CPU) opLD_E_E() { c.ld_r_r(&c.e, c.e) } // 5B
func (c *CPU) opLD_L_E() { c.ld_r_r(&c.l, c.e) } // 6B
func (c *CPU) opLD_A_E() { c.ld_r_r(&c.a, c.e) } // 7B
func (c *CPU) opLD_C_H() { c.ld_r_r(&c.c, c.h) } // 4C
func (c *CPU) opLD_E_H() { c.ld_r_r(&c.e, c.h) } // 5C
func (c *CPU) opLD_L_H() { c.ld_r_r(&c.l, c.h) } // 6C
func (c *CPU) opLD_A_H() { c.ld_r_r(&c.a, c.h) } // 7C
func (c *CPU) opLD_C_L() { c.ld_r_r(&c.c, c.l) } // 4D
func (c *CPU) opLD_E_L() { c.ld_r_r(&c.e, c.l) } // 5D
func (c *CPU) opLD_L_L() { c.ld_r_r(&c.l, c.l) } // 6D
func (c *CPU) opLD_A_L() { c.ld_r_r(&c.a, c.l) } // 7D
func (c *CPU) opLD_C_A() { c.ld_r_r(&c.c, c.a) } // 4F
func (c *CPU) opLD_E_A() { c.ld_r_r(&c.e, c.a) } // 5F
func (c *CPU) opLD_L_A() { c.ld_r_r(&c.l, c.a) } // 6F
func (c *CPU) opLD_A_A() { c.ld_r_r(&c.a, c.a) } // 7F

// -------------------------- LD [HL], r8 --------------------------------------
func (c *CPU) ld_aHL_r(src byte) {
	c.write(c.GetHL(), src)
	c.cycles += 8
}

func (c *CPU) opLD_aHL_B() { c.ld_aHL_r(c.b) } // 70
func (c *CPU) opLD_aHL_C() { c.ld_aHL_r(c.c) } // 71
func (c *CPU) opLD_aHL_D() { c.ld_aHL_r(c.d) } // 72
func (c *CPU) opLD_aHL_E() { c.ld_aHL_r(c.e) } // 73
func (c *CPU) opLD_aHL_H() { c.ld_aHL_r(c.h) } // 74
func (c *CPU) opLD_aHL_L() { c.ld_aHL_r(c.l) } // 75
func (c *CPU) opLD_aHL_A() { c.ld_aHL_r(c.a) } // 77

// -------------------------- LD r8, [HL] --------------------------------------
func (c *CPU) ld_r_aHL(dst *byte) {
	*dst = c.read(c.GetHL())
	c.cycles += 8
}
func (c *CPU) opLD_B_aHL() { c.ld_r_aHL(&c.b) } // 46
func (c *CPU) opLD_D_aHL() { c.ld_r_aHL(&c.d) } // 56
func (c *CPU) opLD_H_aHL() { c.ld_r_aHL(&c.h) } // 66
func (c *CPU) opLD_C_aHL() { c.ld_r_aHL(&c.c) } // 4E
func (c *CPU) opLD_E_aHL() { c.ld_r_aHL(&c.e) } // 5E
func (c *CPU) opLD_L_aHL() { c.ld_r_aHL(&c.l) } // 6E
func (c *CPU) opLD_A_aHL() { c.ld_r_aHL(&c.a) } // 7E

// -------------------------- LD HL, SP + e8 -----------------------------------
func (c *CPU) opLD_HL_SP_p_e8() { // F8
	offset := int8(c.fetch())

	sp := c.sp
	result := uint16(int32(sp) + int32(offset))

	// Good: uint8(int8(-1)) == 0xFF
	// Bad:  uint16(-1) == 0xFFFF
	uOffset := uint16(uint8(offset))

	c.SetFlagZ(false)
	c.SetFlagN(false)
	cResult := ((sp & 0xFF) + (uOffset & 0xFF)) > 0xFF
	c.SetFlagC(cResult)

	hResult := ((sp & 0xF) + (uOffset & 0xF)) > 0xF
	c.SetFlagH(hResult)

	c.SetHL(result)
	c.cycles += 12
}

// ------------------------------ LD SP, HL ------------------------------------
func (c *CPU) opLD_SP_HL() { // F9
	c.sp = c.GetHL()
	c.cycles += 8
}

// ------------------------------ LD [a16], A ----------------------------------
func (c *CPU) opLD_aa16_A() { // EA
	addr := c.fetch16()
	c.write(addr, c.a)
	c.cycles += 16
}

// --------------------------------- LD A, [a16] -------------------------------
func (c *CPU) opLD_A_aa16() { // FA
	addr := c.fetch16()
	val := c.read(addr)
	c.a = val
	c.cycles += 16
}

// --------------------------------- LDH [a8], A -------------------------------
func (c *CPU) opLDH_aa8_A() { // E0
	addr := 0xFF00 + uint16(c.fetch())
	c.write(addr, c.a)
	c.cycles += 12
}

// --------------------------------- LDH A, [a8] -------------------------------
func (c *CPU) opLDH_A_aa8() { // F0
	addr := 0xFF00 + uint16(c.fetch())
	val := c.read(addr)
	c.a = val
	c.cycles += 12
}

// --------------------------------- LDH [C], A --------------------------------
func (c *CPU) opLDH_aC_A() { // E2
	addr := 0xFF00 + uint16(c.c)
	c.write(addr, c.a)
	c.cycles += 8
}

// --------------------------------- LDH A, [C] --------------------------------
func (c *CPU) opLDH_A_aC() { // F2
	addr := 0xFF00 + uint16(c.c)
	val := c.read(addr)
	c.a = val
	c.cycles += 8
}

// --------------------------------- JR ----------------------------------------
func (c *CPU) opJR_NZ_e8() { // 20
	e8 := int8(c.fetch())
	if !c.GetFlagZ() {
		c.pc = uint16(int32(c.pc) + int32(e8))
		c.cycles += 12
	} else {
		c.cycles += 8
	}
}
func (c *CPU) opJR_NC_e8() { // 30
	e8 := int8(c.fetch())
	if !c.GetFlagC() {
		c.pc = uint16(int32(c.pc) + int32(e8))
		c.cycles += 12
	} else {
		c.cycles += 8
	}
}
func (c *CPU) opJR_e8() { // 18
	e8 := int8(c.fetch())
	c.pc = uint16(int32(c.pc) + int32(e8))
	c.cycles += 12
}
func (c *CPU) opJR_Z_e8() { // 28
	e8 := int8(c.fetch())
	if c.GetFlagZ() {
		c.pc = uint16(int32(c.pc) + int32(e8))
		c.cycles += 12
	} else {
		c.cycles += 8
	}
}
func (c *CPU) opJR_C_e8() { // 38
	e8 := int8(c.fetch())
	if c.GetFlagC() {
		c.pc = uint16(int32(c.pc) + int32(e8))
		c.cycles += 12
	} else {
		c.cycles += 8
	}
}

// --------------------------------- JP ----------------------------------------
func (c *CPU) opJP_NZ_a16() { // C2
	addr := c.fetch16()
	if !c.GetFlagZ() {
		c.pc = addr
		c.cycles += 16
	} else {
		c.cycles += 12
	}
}
func (c *CPU) opJP_NC_a16() { // D2
	addr := c.fetch16()
	if !c.GetFlagC() {
		c.pc = addr
		c.cycles += 16
	} else {
		c.cycles += 12
	}
}
func (c *CPU) opJP_a16() { // C3
	addr := c.fetch16()
	c.pc = addr
	c.cycles += 16
}
func (c *CPU) opJP_HL() { // E9
	c.pc = c.GetHL()
	c.cycles += 4
}
func (c *CPU) opJP_Z_a16() { // CA
	addr := c.fetch16()
	if c.GetFlagZ() {
		c.pc = addr
		c.cycles += 16
	} else {
		c.cycles += 12
	}
}
func (c *CPU) opJP_C_a16() { // DA
	addr := c.fetch16()
	if c.GetFlagC() {
		c.pc = addr
		c.cycles += 16
	} else {
		c.cycles += 12
	}
}

// ----------------------------- ADD A, r8 -----------------------------------
func (c *CPU) add_r(src byte) {
	a := c.a
	sum := uint16(a) + uint16(src)
	half := (a & 0xF) + (src & 0xF)
	result := byte(sum)

	c.SetFlagZ(result == 0)
	c.SetFlagN(false)
	c.SetFlagH(half > 0xF)
	c.SetFlagC(sum > 0xFF)

	c.a = result
	c.cycles += 4
}
func (c *CPU) opADD_B() { c.add_r(c.b) } // 80
func (c *CPU) opADD_C() { c.add_r(c.c) } // 81
func (c *CPU) opADD_D() { c.add_r(c.d) } // 82
func (c *CPU) opADD_E() { c.add_r(c.e) } // 83
func (c *CPU) opADD_H() { c.add_r(c.h) } // 84
func (c *CPU) opADD_L() { c.add_r(c.l) } // 85
func (c *CPU) opADD_A() { c.add_r(c.a) } // 87

// ----------------------------------- ADD r16, r16 ----------------------------
func (c *CPU) opADD_HL_BC() { // 09
	hl := uint32(c.GetHL())
	bc := uint32(c.GetBC())

	sum := hl + bc
	sumH := (hl & 0xFFF) + (bc & 0xFFF)

	c.SetFlagN(false)
	c.SetFlagH(sumH > 0xFFF)
	c.SetFlagC(sum > 0xFFFF)

	c.SetHL(uint16(sum))
	c.cycles += 8
}
func (c *CPU) opADD_HL_DE() { // 19
	hl := uint32(c.GetHL())
	de := uint32(c.GetDE())

	sum := hl + de
	sumH := (hl & 0xFFF) + (de & 0xFFF)

	c.SetFlagN(false)
	c.SetFlagH(sumH > 0xFFF)
	c.SetFlagC(sum > 0xFFFF)

	c.SetHL(uint16(sum))
	c.cycles += 8
}
func (c *CPU) opADD_HL_HL() { // 29
	hl := uint32(c.GetHL())

	sum := hl + hl
	sumH := (hl & 0xFFF) + (hl & 0xFFF)

	c.SetFlagN(false)
	c.SetFlagH(sumH > 0xFFF)
	c.SetFlagC(sum > 0xFFFF)

	c.SetHL(uint16(sum))
	c.cycles += 8
}
func (c *CPU) opADD_HL_SP() { // 39
	hl := uint32(c.GetHL())
	sp := uint32(c.sp)

	sum := hl + sp
	sumH := (hl & 0xFFF) + (sp & 0xFFF)

	c.SetFlagN(false)
	c.SetFlagH(sumH > 0xFFF)
	c.SetFlagC(sum > 0xFFFF)

	c.SetHL(uint16(sum))
	c.cycles += 8
}

// ----------------------------- ADD others ------------------------------------
func (c *CPU) opADD_A_aHL() { // 86
	a := uint16(c.a)
	b := uint16(c.read(c.GetHL()))

	sum := a + b
	half := (a & 0xF) + (b & 0xF)
	result := byte(sum)

	c.SetFlagZ(result == 0)
	c.SetFlagN(false)
	c.SetFlagH(half > 0xF)
	c.SetFlagC(sum > 0xFF)

	c.a = result
	c.cycles += 8
}
func (c *CPU) opADD_A_n8() { // C6
	a := uint16(c.a)
	b := uint16(c.fetch())

	sum := a + b
	half := (a & 0xF) + (b & 0xF)
	result := byte(sum)

	c.SetFlagZ(result == 0)
	c.SetFlagN(false)
	c.SetFlagH(half > 0xF)
	c.SetFlagC(sum > 0xFF)

	c.a = result
	c.cycles += 8
}
func (c *CPU) opADD_SP_e8() { // E8
	offset := int8(c.fetch())

	sp := c.sp
	result := uint16(int32(sp) + int32(offset))

	// Good: uint8(int8(-1)) == 0xFF
	// Bad:  uint16(-1)      == 0xFFFF
	uOffset := uint16(uint8(offset))

	c.SetFlagZ(false)
	c.SetFlagN(false)
	cResult := ((sp & 0xFF) + (uOffset & 0xFF)) > 0xFF
	c.SetFlagC(cResult)

	hResult := ((sp & 0xF) + (uOffset & 0xF)) > 0xF
	c.SetFlagH(hResult)

	c.sp = result
	c.cycles += 16
}

// --------------------------- ADC A, r8 --------------------------------------
func (c *CPU) adc_r(src byte) {
	carry := 0
	if c.GetFlagC() {
		carry = 1
	}

	a := c.a
	sum := uint16(a) + uint16(src) + uint16(carry)
	half := (a & 0xF) + (src & 0xF) + byte(carry)
	result := byte(sum)

	c.SetFlagZ(result == 0)
	c.SetFlagN(false)
	c.SetFlagH(half > 0xF)
	c.SetFlagC(sum > 0xFF)

	c.a = result
	c.cycles += 4
}
func (c *CPU) opADC_B() { c.adc_r(c.b) } // 88
func (c *CPU) opADC_C() { c.adc_r(c.c) } // 89
func (c *CPU) opADC_D() { c.adc_r(c.d) } // 8A
func (c *CPU) opADC_E() { c.adc_r(c.e) } // 8B
func (c *CPU) opADC_H() { c.adc_r(c.h) } // 8C
func (c *CPU) opADC_L() { c.adc_r(c.l) } // 8D
func (c *CPU) opADC_A() { c.adc_r(c.a) } // 8F

// --------------------------- ADC A, n8 ---------------------------------------
func (c *CPU) opADC_A_n8() { // CE
	carry := 0
	if c.GetFlagC() {
		carry = 1
	}

	a := c.a
	b := c.fetch()

	sum := uint16(a) + uint16(b) + uint16(carry)
	half := (a & 0xF) + (b & 0xF) + byte(carry)
	result := byte(sum)

	c.SetFlagZ(result == 0)
	c.SetFlagN(false)
	c.SetFlagH(half > 0xF)
	c.SetFlagC(sum > 0xFF)

	c.a = result
	c.cycles += 8
}

// --------------------------- ADC A, [HL] -------------------------------------
func (c *CPU) opADC_A_aHL() { // 8E
	carry := 0
	if c.GetFlagC() {
		carry = 1
	}

	a := c.a
	b := c.read(c.GetHL())

	sum := uint16(a) + uint16(b) + uint16(carry)
	half := (a & 0xF) + (b & 0xF) + byte(carry)
	result := byte(sum)

	c.SetFlagZ(result == 0)
	c.SetFlagN(false)
	c.SetFlagH(half > 0xF)
	c.SetFlagC(sum > 0xFF)

	c.a = result
	c.cycles += 8
}

// ----------------------------- SUB A, r8 ------------------------------------
func (c *CPU) sub_r(src byte) {
	a := c.a
	diff := byte(a - src)

	c.SetFlagZ(diff == 0)
	c.SetFlagN(true)
	c.SetFlagH((a & 0xF) < (src & 0xF))
	c.SetFlagC(a < src)

	c.a = diff
	c.cycles += 4
}
func (c *CPU) opSUB_B() { c.sub_r(c.b) } // 90
func (c *CPU) opSUB_C() { c.sub_r(c.c) } // 91
func (c *CPU) opSUB_D() { c.sub_r(c.d) } // 92
func (c *CPU) opSUB_E() { c.sub_r(c.e) } // 93
func (c *CPU) opSUB_H() { c.sub_r(c.h) } // 94
func (c *CPU) opSUB_L() { c.sub_r(c.l) } // 95
func (c *CPU) opSUB_A() { c.sub_r(c.a) } // 97

// ----------------------------- SUB A, [HL] -----------------------------------
func (c *CPU) opSUB_A_aHL() { // 96
	a := c.a
	b := c.read(c.GetHL())

	diff := byte(a - b)

	c.SetFlagZ(diff == 0)
	c.SetFlagN(true)
	c.SetFlagH((a & 0xF) < (b & 0xF))
	c.SetFlagC(a < b)

	c.a = diff
	c.cycles += 8
}

// ------------------------------ SUB A, n8 ------------------------------------
func (c *CPU) opSUB_A_n8() { // D6
	a := c.a
	b := c.fetch()

	diff := byte(a - b)

	c.SetFlagZ(diff == 0)
	c.SetFlagN(true)
	c.SetFlagH((a & 0xF) < (b & 0xF))
	c.SetFlagC(a < b)

	c.a = diff
	c.cycles += 8
}

// ----------------------------- SBC A, r8 ------------------------------------
func (c *CPU) sbc_r(src byte) {
	carry := 0
	if c.GetFlagC() {
		carry = 1
	}

	a16 := uint16(c.a)
	src16 := uint16(src)
	carry16 := uint16(carry)

	diff := byte(a16 - src16 - carry16)

	c.SetFlagZ(diff == 0)
	c.SetFlagN(true)
	c.SetFlagH((a16 & 0xF) < ((src16 & 0xF) + carry16))
	c.SetFlagC(a16 < (src16 + carry16))

	c.a = diff
	c.cycles += 4
}
func (c *CPU) opSBC_B() { c.sbc_r(c.b) } // 98
func (c *CPU) opSBC_C() { c.sbc_r(c.c) } // 99
func (c *CPU) opSBC_D() { c.sbc_r(c.d) } // 9A
func (c *CPU) opSBC_E() { c.sbc_r(c.e) } // 9B
func (c *CPU) opSBC_H() { c.sbc_r(c.h) } // 9C
func (c *CPU) opSBC_L() { c.sbc_r(c.l) } // 9D
func (c *CPU) opSBC_A() { c.sbc_r(c.a) } // 9F

// ------------------------------ SBC A, [HL] ----------------------------------
func (c *CPU) opSBC_A_aHL() { // 9E
	carry := 0
	if c.GetFlagC() {
		carry = 1
	}

	a16 := uint16(c.a)
	b16 := uint16(c.read(c.GetHL()))
	carry16 := uint16(carry)

	diff := byte(a16 - b16 - carry16)

	c.SetFlagZ(diff == 0)
	c.SetFlagN(true)
	c.SetFlagH((a16 & 0xF) < ((b16 & 0xF) + carry16))
	c.SetFlagC(a16 < (b16 + carry16))

	c.a = diff
	c.cycles += 8
}

// ---------------------------- SBC A, n8 --------------------------------------
func (c *CPU) opSBC_A_n8() { // DE
	carry := 0
	if c.GetFlagC() {
		carry = 1
	}

	a16 := uint16(c.a)
	b16 := uint16(c.fetch())
	carry16 := uint16(carry)

	diff := byte(a16 - b16 - carry16)

	c.SetFlagZ(diff == 0)
	c.SetFlagN(true)
	c.SetFlagH((a16 & 0xF) < ((b16 & 0xF) + carry16))
	c.SetFlagC(a16 < (b16 + carry16))

	c.a = diff
	c.cycles += 8
}

// ----------------------------- AND A, r8 ------------------------------------
func (c *CPU) and_r(b byte) {
	result := c.a & b

	if result == 0 {
		c.SetFlagZ(true)
	} else {
		c.SetFlagZ(false)
	}

	c.SetFlagN(false)
	c.SetFlagH(true) // AND is always half carry
	c.SetFlagC(false)

	c.a = result
	c.cycles += 4
}
func (c *CPU) opAND_B() { c.and_r(c.b) } // A0
func (c *CPU) opAND_C() { c.and_r(c.c) } // A1
func (c *CPU) opAND_D() { c.and_r(c.d) } // A2
func (c *CPU) opAND_E() { c.and_r(c.e) } // A3
func (c *CPU) opAND_H() { c.and_r(c.h) } // A4
func (c *CPU) opAND_L() { c.and_r(c.l) } // A5
func (c *CPU) opAND_A() { c.and_r(c.a) } // A7

// ----------------------------- AND A, [HL] -----------------------------------
func (c *CPU) opAND_A_aHL() { // A6
	result := c.a & c.read(c.GetHL())

	if result == 0 {
		c.SetFlagZ(true)
	} else {
		c.SetFlagZ(false)
	}

	c.SetFlagN(false)
	c.SetFlagH(true) // AND is always half carry
	c.SetFlagC(false)

	c.a = result
	c.cycles += 8
}

// ----------------------------- AND A, n8 -------------------------------------
func (c *CPU) opAND_A_n8() { // E6
	result := c.a & c.fetch()

	if result == 0 {
		c.SetFlagZ(true)
	} else {
		c.SetFlagZ(false)
	}

	c.SetFlagN(false)
	c.SetFlagH(true) // AND is always half carry
	c.SetFlagC(false)

	c.a = result
	c.cycles += 8
}

// ----------------------------- OR A, r8 ------------------------------------
func (c *CPU) or_r(b byte) {
	result := c.a | b

	if result == 0 {
		c.SetFlagZ(true)
	} else {
		c.SetFlagZ(false)
	}

	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(false)

	c.a = result
	c.cycles += 4
}
func (c *CPU) opOR_B() { c.or_r(c.b) } // B0
func (c *CPU) opOR_C() { c.or_r(c.c) } // B1
func (c *CPU) opOR_D() { c.or_r(c.d) } // B2
func (c *CPU) opOR_E() { c.or_r(c.e) } // B3
func (c *CPU) opOR_H() { c.or_r(c.h) } // B4
func (c *CPU) opOR_L() { c.or_r(c.l) } // B5
func (c *CPU) opOR_A() { c.or_r(c.a) } // B7

// ----------------------------- OR A, [HL] -----------------------------------
func (c *CPU) opOR_A_aHL() { // B6
	result := c.a | c.read(c.GetHL())

	if result == 0 {
		c.SetFlagZ(true)
	} else {
		c.SetFlagZ(false)
	}

	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(false)

	c.a = result
	c.cycles += 8
}

// ----------------------------- OR A, n8 -------------------------------------
func (c *CPU) opOR_A_n8() { // F6
	result := c.a | c.fetch()

	if result == 0 {
		c.SetFlagZ(true)
	} else {
		c.SetFlagZ(false)
	}

	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(false)

	c.a = result
	c.cycles += 8
}

// ----------------------------- XOR A, r8 ------------------------------------
func (c *CPU) xor_r(b byte) {
	result := c.a ^ b

	if result == 0 {
		c.SetFlagZ(true)
	} else {
		c.SetFlagZ(false)
	}

	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(false)

	c.a = result
	c.cycles += 4
}
func (c *CPU) opXOR_B() { c.xor_r(c.b) } // A8
func (c *CPU) opXOR_C() { c.xor_r(c.c) } // A9
func (c *CPU) opXOR_D() { c.xor_r(c.d) } // AA
func (c *CPU) opXOR_E() { c.xor_r(c.e) } // AB
func (c *CPU) opXOR_H() { c.xor_r(c.h) } // AC
func (c *CPU) opXOR_L() { c.xor_r(c.l) } // AD
func (c *CPU) opXOR_A() { c.xor_r(c.a) } // AF

// ----------------------------- XOR A, [HL] -----------------------------------
func (c *CPU) opXOR_A_aHL() { // AE
	result := c.a ^ c.read(c.GetHL())

	if result == 0 {
		c.SetFlagZ(true)
	} else {
		c.SetFlagZ(false)
	}

	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(false)

	c.a = result
	c.cycles += 8
}

// ----------------------------- XOR A, n8 -------------------------------------
func (c *CPU) opXOR_A_n8() { // EE
	result := c.a ^ c.fetch()

	if result == 0 {
		c.SetFlagZ(true)
	} else {
		c.SetFlagZ(false)
	}

	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(false)

	c.a = result
	c.cycles += 8
}

// ----------------------------- CP A, r8 ------------------------------------
func (c *CPU) cp_r(b byte) {
	a := c.a
	c.SetFlagZ(a == b)
	c.SetFlagN(true)
	c.SetFlagH((a & 0xF) < (b & 0xF))
	c.SetFlagC(a < b)
	c.cycles += 4
}
func (c *CPU) opCP_B() { c.cp_r(c.b) } // B8
func (c *CPU) opCP_C() { c.cp_r(c.c) } // B9
func (c *CPU) opCP_D() { c.cp_r(c.d) } // BA
func (c *CPU) opCP_E() { c.cp_r(c.e) } // BB
func (c *CPU) opCP_H() { c.cp_r(c.h) } // BC
func (c *CPU) opCP_L() { c.cp_r(c.l) } // BD
func (c *CPU) opCP_A() { c.cp_r(c.a) } // BF

// ----------------------------- CP A, [HL] -----------------------------------
func (c *CPU) opCP_A_aHL() { // BE
	a := c.a
	b := c.read(c.GetHL())

	c.SetFlagZ(a == b)
	c.SetFlagN(true)
	c.SetFlagH((a & 0xF) < (b & 0xF))
	c.SetFlagC(a < b)
	c.cycles += 8
}

// ------------------------------ CP A, n8 ------------------------------------
func (c *CPU) opCP_A_n8() { // FE
	a := c.a
	b := c.fetch()

	c.SetFlagZ(a == b)
	c.SetFlagN(true)
	c.SetFlagH((a & 0xF) < (b & 0xF))
	c.SetFlagC(a < b)
	c.cycles += 8
}

// ------------------------------ INC r8 ---------------------------------------
func (c *CPU) inc_r(r *byte) {
	result := *r + 1

	c.SetFlagZ(result == 0)
	c.SetFlagN(false)
	c.SetFlagH(((*r & 0xF) + 1) > 0xF)

	*r = result
	c.cycles += 4
}
func (c *CPU) opINC_B() { c.inc_r(&c.b) } // 04
func (c *CPU) opINC_D() { c.inc_r(&c.d) } // 14
func (c *CPU) opINC_H() { c.inc_r(&c.h) } // 24
func (c *CPU) opINC_C() { c.inc_r(&c.c) } // 0C
func (c *CPU) opINC_E() { c.inc_r(&c.e) } // 1C
func (c *CPU) opINC_L() { c.inc_r(&c.l) } // 2C
func (c *CPU) opINC_A() { c.inc_r(&c.a) } // 3C

// ------------------------------ INC r16 --------------------------------------
func (c *CPU) opINC_BC() { // 03
	c.SetBC(c.GetBC() + 1)
	c.cycles += 8
}
func (c *CPU) opINC_DE() { // 13
	c.SetDE(c.GetDE() + 1)
	c.cycles += 8
}
func (c *CPU) opINC_HL() { // 23
	c.SetHL(c.GetHL() + 1)
	c.cycles += 8
}
func (c *CPU) opINC_SP() { // 33
	c.sp++
	c.cycles += 8
}

// ------------------------------ INC [HL] -------------------------------------
func (c *CPU) opINC_aHL() { // 34
	addr := c.GetHL()
	val := c.read(addr)
	result := val + 1

	c.SetFlagZ(result == 0)
	c.SetFlagN(false)
	c.SetFlagH(((val & 0xF) + 1) > 0xF)

	c.write(addr, result)
	c.cycles += 12
}

// ------------------------------ DEC r8 ---------------------------------------
func (c *CPU) dec_r(r *byte) {
	result := byte(*r - 1)

	c.SetFlagZ(result == 0)
	c.SetFlagN(true)
	c.SetFlagH((*r & 0xF) < 1)

	*r = result
	c.cycles += 4
}
func (c *CPU) opDEC_B() { c.dec_r(&c.b) } // 05
func (c *CPU) opDEC_D() { c.dec_r(&c.d) } // 15
func (c *CPU) opDEC_H() { c.dec_r(&c.h) } // 25
func (c *CPU) opDEC_C() { c.dec_r(&c.c) } // 0D
func (c *CPU) opDEC_E() { c.dec_r(&c.e) } // 1D
func (c *CPU) opDEC_L() { c.dec_r(&c.l) } // 2D
func (c *CPU) opDEC_A() { c.dec_r(&c.a) } // 3D

// ------------------------------ DEC r16 --------------------------------------
func (c *CPU) opDEC_BC() { // 0B
	c.SetBC(c.GetBC() - 1)
	c.cycles += 8
}
func (c *CPU) opDEC_DE() { // 1B
	c.SetDE(c.GetDE() - 1)
	c.cycles += 8
}
func (c *CPU) opDEC_HL() { // 2B
	c.SetHL(c.GetHL() - 1)
	c.cycles += 8
}
func (c *CPU) opDEC_SP() { // 3B
	c.sp--
	c.cycles += 8
}

// ------------------------------ DEC [HL] -------------------------------------
func (c *CPU) opDEC_aHL() { // 35
	addr := c.GetHL()
	val := c.read(addr)
	result := byte(val - 1)

	c.SetFlagZ(result == 0)
	c.SetFlagN(true)
	c.SetFlagH((val & 0xF) < 1)

	c.write(addr, result)
	c.cycles += 12
}

// ------------------------------- RLCA ----------------------------------------
func (c *CPU) opRLCA() { // 07
	c.SetFlagZ(false)
	c.SetFlagN(false)
	c.SetFlagH(false)

	bit7 := (c.a & 0x80) >> 7
	c.a = (c.a << 1) + bit7
	c.SetFlagC(bit7 == 1)
	c.cycles += 4
}

// -------------------------------- RLA ----------------------------------------
func (c *CPU) opRLA() { // 17
	c.SetFlagZ(false)
	c.SetFlagN(false)
	c.SetFlagH(false)

	bit7 := (c.a & 0x80) >> 7

	carry := byte(0)
	if c.GetFlagC() {
		carry = 1
	}

	c.a = (c.a << 1) + carry
	c.SetFlagC(bit7 == 1)
	c.cycles += 4
}

// --------------------------------- RRCA --------------------------------------
func (c *CPU) opRRCA() { // 0F
	c.SetFlagZ(false)
	c.SetFlagN(false)
	c.SetFlagH(false)

	bit0 := c.a & 1
	c.a = (c.a >> 1) | (bit0 << 7)
	c.SetFlagC(bit0 == 1)
	c.cycles += 4
}

// -------------------------------- RRA ----------------------------------------
func (c *CPU) opRRA() { // 1F
	c.SetFlagZ(false)
	c.SetFlagN(false)
	c.SetFlagH(false)

	bit0 := c.a & 1

	carry := byte(0)
	if c.GetFlagC() {
		carry = 1
	}

	c.a = (c.a >> 1) | (carry << 7)
	c.SetFlagC(bit0 == 1)
	c.cycles += 4
}

// --------------------------------- DAA ---------------------------------------
func (c *CPU) opDAA() { // 27
	a := c.a
	oldA := a

	if c.GetFlagN() {
		// After SUB
		if c.GetFlagH() {
			a -= 0x06
		}
		if c.GetFlagC() {
			a -= 0x60
		}
	} else {
		// After ADD
		if c.GetFlagH() || ((a & 0x0F) > 9) {
			a += 0x06
		}
		if c.GetFlagC() || oldA > 0x99 {
			a += 0x60
			c.SetFlagC(true)
		}
	}

	c.a = a

	c.SetFlagZ(c.a == 0)
	c.SetFlagH(false)
	c.cycles += 4
}

// -------------------------------- SCF ----------------------------------------
func (c *CPU) opSCF() { // 37
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(true)
	c.cycles += 4
}

// --------------------------------- CPL ---------------------------------------
func (c *CPU) opCPL() { // 2F
	c.a = ^c.a
	c.SetFlagN(true)
	c.SetFlagH(true)
	c.cycles += 4
}

// -------------------------------- CCF ----------------------------------------
func (c *CPU) opCCF() { // 3F
	if c.GetFlagC() {
		c.SetFlagC(false)
	} else {
		c.SetFlagC(true)
	}
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.cycles += 4
}

// ---------------------------- PUSH r16 ---------------------------------------
func (c *CPU) push(rr uint16) {
	high := byte(rr >> 8)
	low := byte(rr & 0x00FF)

	c.sp--
	c.write(c.sp, high)
	c.sp--
	c.write(c.sp, low)

}

func (c *CPU) opPUSH_BC() { // C5
	c.push(c.GetBC())
	c.cycles += 16
}
func (c *CPU) opPUSH_DE() { // D5
	c.push(c.GetDE())
	c.cycles += 16
}
func (c *CPU) opPUSH_HL() { // E5
	c.push(c.GetHL())
	c.cycles += 16
}
func (c *CPU) opPUSH_AF() { // F5
	c.push(c.GetAF())
	c.cycles += 16
}

// ---------------------------- CALL ------------------------------------------
func (c *CPU) opCALL_NZ_a16() { // C4
	addr := c.fetch16()
	if !c.GetFlagZ() {
		c.push(c.pc)
		c.pc = addr
		c.cycles += 24
	} else {
		c.cycles += 12
	}
}

func (c *CPU) opCALL_Z_a16() { // CC
	addr := c.fetch16()
	if c.GetFlagZ() {
		c.push(c.pc)
		c.pc = addr
		c.cycles += 24
	} else {
		c.cycles += 12
	}
}

func (c *CPU) opCALL_a16() { // CD
	addr := c.fetch16()
	c.push(c.pc)
	c.pc = addr
	c.cycles += 24
}

func (c *CPU) opCALL_NC_a16() { // D4
	addr := c.fetch16()
	if !c.GetFlagC() {
		c.push(c.pc)
		c.pc = addr
		c.cycles += 24
	} else {
		c.cycles += 12
	}
}

func (c *CPU) opCALL_C_a16() { // DC
	addr := c.fetch16()
	if c.GetFlagC() {
		c.push(c.pc)
		c.pc = addr
		c.cycles += 24
	} else {
		c.cycles += 12
	}
}

// ------------------------------- RST -----------------------------------------
func (c *CPU) opRST_00() { // C7
	c.push(c.pc)
	c.pc = 0x0000
	c.cycles += 16
}
func (c *CPU) opRST_10() { // D7
	c.push(c.pc)
	c.pc = 0x0010
	c.cycles += 16
}
func (c *CPU) opRST_20() { // E7
	c.push(c.pc)
	c.pc = 0x0020
	c.cycles += 16
}
func (c *CPU) opRST_30() { // F7
	c.push(c.pc)
	c.pc = 0x0030
	c.cycles += 16
}
func (c *CPU) opRST_08() { // CF
	c.push(c.pc)
	c.pc = 0x0008
	c.cycles += 16
}
func (c *CPU) opRST_18() { // DF
	c.push(c.pc)
	c.pc = 0x0018
	c.cycles += 16
}
func (c *CPU) opRST_28() { // EF
	c.push(c.pc)
	c.pc = 0x0028
	c.cycles += 16
}
func (c *CPU) opRST_38() { // FF
	c.push(c.pc)
	c.pc = 0x0038
	c.cycles += 16
}

// ---------------------------- POP r16 ----------------------------------------
func (c *CPU) pop() uint16 {
	lo := c.read(c.sp)
	c.sp++
	hi := c.read(c.sp)
	c.sp++
	return (uint16(hi) << 8) | uint16(lo)
}
func (c *CPU) opPOP_BC() { // C1
	c.SetBC(c.pop())
	c.cycles += 12
}
func (c *CPU) opPOP_DE() { // D1
	c.SetDE(c.pop())
	c.cycles += 12
}
func (c *CPU) opPOP_HL() { // E1
	c.SetHL(c.pop())
	c.cycles += 12
}
func (c *CPU) opPOP_AF() { // F1
	c.SetAF(c.pop() & 0xFFF0)
	c.cycles += 12
}

// -------------------------------- RET ----------------------------------------
func (c *CPU) opRET() { // C9
	c.pc = c.pop()
	c.cycles += 16
}
func (c *CPU) opRET_Z() { // C8
	if c.GetFlagZ() {
		c.pc = c.pop()
		c.cycles += 20
	} else {
		c.cycles += 8
	}
}
func (c *CPU) opRET_C() { // D8
	if c.GetFlagC() {
		c.pc = c.pop()
		c.cycles += 20
	} else {
		c.cycles += 8
	}
}
func (c *CPU) opRET_NZ() { // C0
	if !c.GetFlagZ() {
		c.pc = c.pop()
		c.cycles += 20
	} else {
		c.cycles += 8
	}
}
func (c *CPU) opRET_NC() { // D0
	if !c.GetFlagC() {
		c.pc = c.pop()
		c.cycles += 20
	} else {
		c.cycles += 8
	}
}

// --------------------------------- RETI --------------------------------------
func (c *CPU) opRETI() { // D9
	c.pc = c.pop()
	c.isIMEEnabled = true
	c.cycles += 16
}
