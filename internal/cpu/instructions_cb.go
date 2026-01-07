package cpu

// ------------------------------- RLC r8 --------------------------------------
func (c *CPU) rlc_r(r *byte) {
	bit7 := (*r & 0x80) >> 7
	*r = *r<<1 | bit7
	c.SetFlagZ(*r == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(bit7 == 1)
	c.cycles += 8
}
func (c *CPU) opRLC_B() { c.rlc_r(&c.b) } // 00
func (c *CPU) opRLC_C() { c.rlc_r(&c.c) } // 01
func (c *CPU) opRLC_D() { c.rlc_r(&c.d) } // 02
func (c *CPU) opRLC_E() { c.rlc_r(&c.e) } // 03
func (c *CPU) opRLC_H() { c.rlc_r(&c.h) } // 04
func (c *CPU) opRLC_L() { c.rlc_r(&c.l) } // 05
func (c *CPU) opRLC_A() { c.rlc_r(&c.a) } // 07

// ------------------------------- RLC [HL] -------------------------------------
func (c *CPU) opRLC_aHL() { // 06
	addr := c.GetHL()
	oldVal := c.read(addr)
	bit7 := (oldVal & 0x80) >> 7

	newVal := oldVal<<1 | bit7
	c.write(addr, newVal)

	c.SetFlagZ(newVal == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(bit7 == 1)

	c.cycles += 16
}

// ------------------------------- RRC r8 --------------------------------------
func (c *CPU) rrc_r(r *byte) {
	bit0 := *r & 1
	*r = (bit0 << 7) | (*r >> 1)
	c.SetFlagZ(*r == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(bit0 == 1)
	c.cycles += 8
}
func (c *CPU) opRRC_B() { c.rrc_r(&c.b) } // 08
func (c *CPU) opRRC_C() { c.rrc_r(&c.c) } // 09
func (c *CPU) opRRC_D() { c.rrc_r(&c.d) } // 0A
func (c *CPU) opRRC_E() { c.rrc_r(&c.e) } // 0B
func (c *CPU) opRRC_H() { c.rrc_r(&c.h) } // 0C
func (c *CPU) opRRC_L() { c.rrc_r(&c.l) } // 0D
func (c *CPU) opRRC_A() { c.rrc_r(&c.a) } // 0F

// ------------------------------- RRC [HL] -------------------------------------
func (c *CPU) opRRC_aHL() { // 0E
	addr := c.GetHL()
	oldVal := c.read(addr)
	bit0 := oldVal & 1

	newVal := (bit0 << 7) | (oldVal >> 1)
	c.write(addr, newVal)

	c.SetFlagZ(newVal == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(bit0 == 1)

	c.cycles += 16
}

// ------------------------------- RL r8 ---------------------------------------
func (c *CPU) rl_r(r *byte) {
	bit7 := (*r & 0x80) >> 7

	carry := byte(0)
	if c.GetFlagC() {
		carry = 1
	}

	*r = *r<<1 | carry
	c.SetFlagZ(*r == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(bit7 == 1)
	c.cycles += 8
}
func (c *CPU) opRL_B() { c.rl_r(&c.b) } // 10
func (c *CPU) opRL_C() { c.rl_r(&c.c) } // 11
func (c *CPU) opRL_D() { c.rl_r(&c.d) } // 12
func (c *CPU) opRL_E() { c.rl_r(&c.e) } // 13
func (c *CPU) opRL_H() { c.rl_r(&c.h) } // 14
func (c *CPU) opRL_L() { c.rl_r(&c.l) } // 15
func (c *CPU) opRL_A() { c.rl_r(&c.a) } // 17

// ------------------------------- RL [HL] --------------------------------------
func (c *CPU) opRL_aHL() { // 16
	addr := c.GetHL()
	oldVal := c.read(addr)
	bit7 := (oldVal & 0x80) >> 7

	carry := byte(0)
	if c.GetFlagC() {
		carry = 1
	}

	newVal := oldVal<<1 | carry
	c.write(addr, newVal)

	c.SetFlagZ(newVal == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(bit7 == 1)

	c.cycles += 16
}

// ------------------------------- RR r8 --------------------------------------
func (c *CPU) rr_r(r *byte) {
	bit0 := *r & 1

	carry := byte(0)
	if c.GetFlagC() {
		carry = 1
	}

	*r = (carry << 7) | (*r >> 1)
	c.SetFlagZ(*r == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(bit0 == 1)
	c.cycles += 8
}
func (c *CPU) opRR_B() { c.rr_r(&c.b) } // 18
func (c *CPU) opRR_C() { c.rr_r(&c.c) } // 19
func (c *CPU) opRR_D() { c.rr_r(&c.d) } // 1A
func (c *CPU) opRR_E() { c.rr_r(&c.e) } // 1B
func (c *CPU) opRR_H() { c.rr_r(&c.h) } // 1C
func (c *CPU) opRR_L() { c.rr_r(&c.l) } // 1D
func (c *CPU) opRR_A() { c.rr_r(&c.a) } // 1F

// ------------------------------- RR [HL] -------------------------------------
func (c *CPU) opRR_aHL() { // 1E
	addr := c.GetHL()
	oldVal := c.read(addr)
	bit0 := oldVal & 1

	carry := byte(0)
	if c.GetFlagC() {
		carry = 1
	}

	newVal := (carry << 7) | (oldVal >> 1)
	c.write(addr, newVal)

	c.SetFlagZ(newVal == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(bit0 == 1)

	c.cycles += 16
}

// ------------------------------- SLA r8 ---------------------------------------
func (c *CPU) sla_r(r *byte) {
	bit7 := (*r & 0x80) >> 7
	*r = *r << 1
	c.SetFlagZ(*r == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(bit7 == 1)
	c.cycles += 8
}
func (c *CPU) opSLA_B() { c.sla_r(&c.b) } // 20
func (c *CPU) opSLA_C() { c.sla_r(&c.c) } // 21
func (c *CPU) opSLA_D() { c.sla_r(&c.d) } // 22
func (c *CPU) opSLA_E() { c.sla_r(&c.e) } // 23
func (c *CPU) opSLA_H() { c.sla_r(&c.h) } // 24
func (c *CPU) opSLA_L() { c.sla_r(&c.l) } // 25
func (c *CPU) opSLA_A() { c.sla_r(&c.a) } // 27

// -------------------------------- SLA [HL] ------------------------------------
func (c *CPU) opSLA_aHL() { // 26
	addr := c.GetHL()
	oldVal := c.read(addr)
	bit7 := (oldVal & 0x80) >> 7

	newVal := oldVal << 1
	c.write(addr, newVal)

	c.SetFlagZ(newVal == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(bit7 == 1)

	c.cycles += 16
}

// ------------------------------- SRA r8 ---------------------------------------
func (c *CPU) sra_r(r *byte) {
	bit0 := *r & 1
	msbMask := *r & 0x80
	*r = msbMask | (*r >> 1)
	c.SetFlagZ(*r == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(bit0 == 1)
	c.cycles += 8
}
func (c *CPU) opSRA_B() { c.sra_r(&c.b) } // 28
func (c *CPU) opSRA_C() { c.sra_r(&c.c) } // 29
func (c *CPU) opSRA_D() { c.sra_r(&c.d) } // 2A
func (c *CPU) opSRA_E() { c.sra_r(&c.e) } // 2B
func (c *CPU) opSRA_H() { c.sra_r(&c.h) } // 2C
func (c *CPU) opSRA_L() { c.sra_r(&c.l) } // 2D
func (c *CPU) opSRA_A() { c.sra_r(&c.a) } // 2F

// -------------------------------- SRA [HL] ------------------------------------
func (c *CPU) opSRA_aHL() { // 2E
	addr := c.GetHL()
	oldVal := c.read(addr)
	bit0 := oldVal & 1
	msbMask := oldVal & 0x80

	newVal := msbMask | (oldVal >> 1)
	c.write(addr, newVal)

	c.SetFlagZ(newVal == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(bit0 == 1)

	c.cycles += 16
}

// -------------------------------- SWAP r8 -------------------------------------
func (c *CPU) swap_r(r *byte) {
	high := (*r & 0xF0) >> 4
	low := *r & 0x0F
	*r = (low << 4) | high

	c.SetFlagZ(*r == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(false)
	c.cycles += 8
}

func (c *CPU) opSWAP_B() { c.swap_r(&c.b) } // 30
func (c *CPU) opSWAP_C() { c.swap_r(&c.c) } // 31
func (c *CPU) opSWAP_D() { c.swap_r(&c.d) } // 32
func (c *CPU) opSWAP_E() { c.swap_r(&c.e) } // 33
func (c *CPU) opSWAP_H() { c.swap_r(&c.h) } // 34
func (c *CPU) opSWAP_L() { c.swap_r(&c.l) } // 35
func (c *CPU) opSWAP_A() { c.swap_r(&c.a) } // 37

// -------------------------------- SWAP [HL] -----------------------------------
func (c *CPU) opSWAP_aHL() { // 36
	addr := c.GetHL()
	oldVal := c.read(addr)
	high := (oldVal & 0xF0) >> 4
	low := oldVal & 0x0F
	newVal := (low << 4) | high
	c.write(addr, newVal)

	c.SetFlagZ(newVal == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(false)

	c.cycles += 16
}

// ------------------------------- SRL r8 ---------------------------------------
func (c *CPU) srl_r(r *byte) {
	bit0 := *r & 1
	*r = *r >> 1
	c.SetFlagZ(*r == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(bit0 == 1)
	c.cycles += 8
}
func (c *CPU) opSRL_B() { c.srl_r(&c.b) } // 38
func (c *CPU) opSRL_C() { c.srl_r(&c.c) } // 39
func (c *CPU) opSRL_D() { c.srl_r(&c.d) } // 3A
func (c *CPU) opSRL_E() { c.srl_r(&c.e) } // 3B
func (c *CPU) opSRL_H() { c.srl_r(&c.h) } // 3C
func (c *CPU) opSRL_L() { c.srl_r(&c.l) } // 3D
func (c *CPU) opSRL_A() { c.srl_r(&c.a) } // 3F

// -------------------------------- SRL [HL] ------------------------------------
func (c *CPU) opSRL_aHL() { // 3E
	addr := c.GetHL()
	oldVal := c.read(c.GetHL())
	bit0 := oldVal & 1

	newVal := oldVal >> 1
	c.write(addr, newVal)

	c.SetFlagZ(newVal == 0)
	c.SetFlagN(false)
	c.SetFlagH(false)
	c.SetFlagC(bit0 == 1)

	c.cycles += 16
}

// ---------------------------------- BIT n, r8 ---------------------------------
func (c *CPU) bit_n_r(n, r byte) {
	bitSet := (r >> n) & 1
	c.SetFlagZ(bitSet == 0)
	c.SetFlagN(false)
	c.SetFlagH(true)
	c.cycles += 8
}
func (c *CPU) opBIT_0_B() { c.bit_n_r(0, c.b) } // 40
func (c *CPU) opBIT_0_C() { c.bit_n_r(0, c.c) } // 41
func (c *CPU) opBIT_0_D() { c.bit_n_r(0, c.d) } // 42
func (c *CPU) opBIT_0_E() { c.bit_n_r(0, c.e) } // 43
func (c *CPU) opBIT_0_H() { c.bit_n_r(0, c.h) } // 44
func (c *CPU) opBIT_0_L() { c.bit_n_r(0, c.l) } // 45
func (c *CPU) opBIT_0_A() { c.bit_n_r(0, c.a) } // 47
func (c *CPU) opBIT_1_B() { c.bit_n_r(1, c.b) } // 48
func (c *CPU) opBIT_1_C() { c.bit_n_r(1, c.c) } // 49
func (c *CPU) opBIT_1_D() { c.bit_n_r(1, c.d) } // 4A
func (c *CPU) opBIT_1_E() { c.bit_n_r(1, c.e) } // 4B
func (c *CPU) opBIT_1_H() { c.bit_n_r(1, c.h) } // 4C
func (c *CPU) opBIT_1_L() { c.bit_n_r(1, c.l) } // 4D
func (c *CPU) opBIT_1_A() { c.bit_n_r(1, c.a) } // 4F
func (c *CPU) opBIT_2_B() { c.bit_n_r(2, c.b) } // 50
func (c *CPU) opBIT_2_C() { c.bit_n_r(2, c.c) } // 51
func (c *CPU) opBIT_2_D() { c.bit_n_r(2, c.d) } // 52
func (c *CPU) opBIT_2_E() { c.bit_n_r(2, c.e) } // 53
func (c *CPU) opBIT_2_H() { c.bit_n_r(2, c.h) } // 54
func (c *CPU) opBIT_2_L() { c.bit_n_r(2, c.l) } // 55
func (c *CPU) opBIT_2_A() { c.bit_n_r(2, c.a) } // 57
func (c *CPU) opBIT_3_B() { c.bit_n_r(3, c.b) } // 58
func (c *CPU) opBIT_3_C() { c.bit_n_r(3, c.c) } // 59
func (c *CPU) opBIT_3_D() { c.bit_n_r(3, c.d) } // 5A
func (c *CPU) opBIT_3_E() { c.bit_n_r(3, c.e) } // 5B
func (c *CPU) opBIT_3_H() { c.bit_n_r(3, c.h) } // 5C
func (c *CPU) opBIT_3_L() { c.bit_n_r(3, c.l) } // 5D
func (c *CPU) opBIT_3_A() { c.bit_n_r(3, c.a) } // 5F
func (c *CPU) opBIT_4_B() { c.bit_n_r(4, c.b) } // 60
func (c *CPU) opBIT_4_C() { c.bit_n_r(4, c.c) } // 61
func (c *CPU) opBIT_4_D() { c.bit_n_r(4, c.d) } // 62
func (c *CPU) opBIT_4_E() { c.bit_n_r(4, c.e) } // 63
func (c *CPU) opBIT_4_H() { c.bit_n_r(4, c.h) } // 64
func (c *CPU) opBIT_4_L() { c.bit_n_r(4, c.l) } // 65
func (c *CPU) opBIT_4_A() { c.bit_n_r(4, c.a) } // 67
func (c *CPU) opBIT_5_B() { c.bit_n_r(5, c.b) } // 68
func (c *CPU) opBIT_5_C() { c.bit_n_r(5, c.c) } // 69
func (c *CPU) opBIT_5_D() { c.bit_n_r(5, c.d) } // 6A
func (c *CPU) opBIT_5_E() { c.bit_n_r(5, c.e) } // 6B
func (c *CPU) opBIT_5_H() { c.bit_n_r(5, c.h) } // 6C
func (c *CPU) opBIT_5_L() { c.bit_n_r(5, c.l) } // 6D
func (c *CPU) opBIT_5_A() { c.bit_n_r(5, c.a) } // 6F
func (c *CPU) opBIT_6_B() { c.bit_n_r(6, c.b) } // 70
func (c *CPU) opBIT_6_C() { c.bit_n_r(6, c.c) } // 71
func (c *CPU) opBIT_6_D() { c.bit_n_r(6, c.d) } // 72
func (c *CPU) opBIT_6_E() { c.bit_n_r(6, c.e) } // 73
func (c *CPU) opBIT_6_H() { c.bit_n_r(6, c.h) } // 74
func (c *CPU) opBIT_6_L() { c.bit_n_r(6, c.l) } // 75
func (c *CPU) opBIT_6_A() { c.bit_n_r(6, c.a) } // 77
func (c *CPU) opBIT_7_B() { c.bit_n_r(7, c.b) } // 78
func (c *CPU) opBIT_7_C() { c.bit_n_r(7, c.c) } // 79
func (c *CPU) opBIT_7_D() { c.bit_n_r(7, c.d) } // 7A
func (c *CPU) opBIT_7_E() { c.bit_n_r(7, c.e) } // 7B
func (c *CPU) opBIT_7_H() { c.bit_n_r(7, c.h) } // 7C
func (c *CPU) opBIT_7_L() { c.bit_n_r(7, c.l) } // 7D
func (c *CPU) opBIT_7_A() { c.bit_n_r(7, c.a) } // 7F

// --------------------------------- BIT n, [HL] --------------------------------
func (c *CPU) bit_n_ahl(n byte) {
	addr := c.GetHL()
	val := c.read(addr)
	bitSet := (val >> n) & 1

	c.SetFlagZ(bitSet == 0)
	c.SetFlagN(false)
	c.SetFlagH(true)
	c.cycles += 12
}
func (c *CPU) opBIT_0_aHL() { c.bit_n_ahl(0) } // 46
func (c *CPU) opBIT_1_aHL() { c.bit_n_ahl(1) } // 4E
func (c *CPU) opBIT_2_aHL() { c.bit_n_ahl(2) } // 56
func (c *CPU) opBIT_3_aHL() { c.bit_n_ahl(3) } // 5E
func (c *CPU) opBIT_4_aHL() { c.bit_n_ahl(4) } // 66
func (c *CPU) opBIT_5_aHL() { c.bit_n_ahl(5) } // 6E
func (c *CPU) opBIT_6_aHL() { c.bit_n_ahl(6) } // 76
func (c *CPU) opBIT_7_aHL() { c.bit_n_ahl(7) } // 7E

// -------------------------------- RES n, r8 -----------------------------------
func (c *CPU) res_n_r(n uint8, r *byte) {
	*r = *r &^ (1 << n)
	c.cycles += 8
}
func (c *CPU) opRES_0_B() { c.res_n_r(0, &c.b) } // 80
func (c *CPU) opRES_0_C() { c.res_n_r(0, &c.c) } // 81
func (c *CPU) opRES_0_D() { c.res_n_r(0, &c.d) } // 82
func (c *CPU) opRES_0_E() { c.res_n_r(0, &c.e) } // 83
func (c *CPU) opRES_0_H() { c.res_n_r(0, &c.h) } // 84
func (c *CPU) opRES_0_L() { c.res_n_r(0, &c.l) } // 85
func (c *CPU) opRES_0_A() { c.res_n_r(0, &c.a) } // 87
func (c *CPU) opRES_1_B() { c.res_n_r(1, &c.b) } // 88
func (c *CPU) opRES_1_C() { c.res_n_r(1, &c.c) } // 89
func (c *CPU) opRES_1_D() { c.res_n_r(1, &c.d) } // 8A
func (c *CPU) opRES_1_E() { c.res_n_r(1, &c.e) } // 8B
func (c *CPU) opRES_1_H() { c.res_n_r(1, &c.h) } // 8C
func (c *CPU) opRES_1_L() { c.res_n_r(1, &c.l) } // 8D
func (c *CPU) opRES_1_A() { c.res_n_r(1, &c.a) } // 8F
func (c *CPU) opRES_2_B() { c.res_n_r(2, &c.b) } // 90
func (c *CPU) opRES_2_C() { c.res_n_r(2, &c.c) } // 91
func (c *CPU) opRES_2_D() { c.res_n_r(2, &c.d) } // 92
func (c *CPU) opRES_2_E() { c.res_n_r(2, &c.e) } // 93
func (c *CPU) opRES_2_H() { c.res_n_r(2, &c.h) } // 94
func (c *CPU) opRES_2_L() { c.res_n_r(2, &c.l) } // 95
func (c *CPU) opRES_2_A() { c.res_n_r(2, &c.a) } // 97
func (c *CPU) opRES_3_B() { c.res_n_r(3, &c.b) } // 98
func (c *CPU) opRES_3_C() { c.res_n_r(3, &c.c) } // 99
func (c *CPU) opRES_3_D() { c.res_n_r(3, &c.d) } // 9A
func (c *CPU) opRES_3_E() { c.res_n_r(3, &c.e) } // 9B
func (c *CPU) opRES_3_H() { c.res_n_r(3, &c.h) } // 9C
func (c *CPU) opRES_3_L() { c.res_n_r(3, &c.l) } // 9D
func (c *CPU) opRES_3_A() { c.res_n_r(3, &c.a) } // 9F
func (c *CPU) opRES_4_B() { c.res_n_r(4, &c.b) } // A0
func (c *CPU) opRES_4_C() { c.res_n_r(4, &c.c) } // A1
func (c *CPU) opRES_4_D() { c.res_n_r(4, &c.d) } // A2
func (c *CPU) opRES_4_E() { c.res_n_r(4, &c.e) } // A3
func (c *CPU) opRES_4_H() { c.res_n_r(4, &c.h) } // A4
func (c *CPU) opRES_4_L() { c.res_n_r(4, &c.l) } // A5
func (c *CPU) opRES_4_A() { c.res_n_r(4, &c.a) } // A7
func (c *CPU) opRES_5_B() { c.res_n_r(5, &c.b) } // A8
func (c *CPU) opRES_5_C() { c.res_n_r(5, &c.c) } // A9
func (c *CPU) opRES_5_D() { c.res_n_r(5, &c.d) } // AA
func (c *CPU) opRES_5_E() { c.res_n_r(5, &c.e) } // AB
func (c *CPU) opRES_5_H() { c.res_n_r(5, &c.h) } // AC
func (c *CPU) opRES_5_L() { c.res_n_r(5, &c.l) } // AD
func (c *CPU) opRES_5_A() { c.res_n_r(5, &c.a) } // AF
func (c *CPU) opRES_6_B() { c.res_n_r(6, &c.b) } // B0
func (c *CPU) opRES_6_C() { c.res_n_r(6, &c.c) } // B1
func (c *CPU) opRES_6_D() { c.res_n_r(6, &c.d) } // B2
func (c *CPU) opRES_6_E() { c.res_n_r(6, &c.e) } // B3
func (c *CPU) opRES_6_H() { c.res_n_r(6, &c.h) } // B4
func (c *CPU) opRES_6_L() { c.res_n_r(6, &c.l) } // B5
func (c *CPU) opRES_6_A() { c.res_n_r(6, &c.a) } // B7
func (c *CPU) opRES_7_B() { c.res_n_r(7, &c.b) } // B8
func (c *CPU) opRES_7_C() { c.res_n_r(7, &c.c) } // B9
func (c *CPU) opRES_7_D() { c.res_n_r(7, &c.d) } // BA
func (c *CPU) opRES_7_E() { c.res_n_r(7, &c.e) } // BB
func (c *CPU) opRES_7_H() { c.res_n_r(7, &c.h) } // BC
func (c *CPU) opRES_7_L() { c.res_n_r(7, &c.l) } // BD
func (c *CPU) opRES_7_A() { c.res_n_r(7, &c.a) } // BF
// -------------------------------- RES n, [HL] ---------------------------------
func (c *CPU) res_n_ahl(n uint8) {
	addr := c.GetHL()
	val := c.read(addr)
	val = val &^ (1 << n)
	c.write(addr, val)
	c.cycles += 16
}
func (c *CPU) opRES_0_aHL() { c.res_n_ahl(0) } // 86
func (c *CPU) opRES_1_aHL() { c.res_n_ahl(1) } // 8E
func (c *CPU) opRES_2_aHL() { c.res_n_ahl(2) } // 96
func (c *CPU) opRES_3_aHL() { c.res_n_ahl(3) } // 9E
func (c *CPU) opRES_4_aHL() { c.res_n_ahl(4) } // A6
func (c *CPU) opRES_5_aHL() { c.res_n_ahl(5) } // AE
func (c *CPU) opRES_6_aHL() { c.res_n_ahl(6) } // B6
func (c *CPU) opRES_7_aHL() { c.res_n_ahl(7) } // BE

// -------------------------------- SET n, r8 -----------------------------------
func (c *CPU) set_n_r(n uint8, r *byte) {
	*r = *r | (1 << n)
	c.cycles += 8
}
func (c *CPU) opSET_0_B() { c.set_n_r(0, &c.b) } // C0
func (c *CPU) opSET_0_C() { c.set_n_r(0, &c.c) } // C1
func (c *CPU) opSET_0_D() { c.set_n_r(0, &c.d) } // C2
func (c *CPU) opSET_0_E() { c.set_n_r(0, &c.e) } // C3
func (c *CPU) opSET_0_H() { c.set_n_r(0, &c.h) } // C4
func (c *CPU) opSET_0_L() { c.set_n_r(0, &c.l) } // C5
func (c *CPU) opSET_0_A() { c.set_n_r(0, &c.a) } // C7
func (c *CPU) opSET_1_B() { c.set_n_r(1, &c.b) } // C8
func (c *CPU) opSET_1_C() { c.set_n_r(1, &c.c) } // C9
func (c *CPU) opSET_1_D() { c.set_n_r(1, &c.d) } // CA
func (c *CPU) opSET_1_E() { c.set_n_r(1, &c.e) } // CB
func (c *CPU) opSET_1_H() { c.set_n_r(1, &c.h) } // CC
func (c *CPU) opSET_1_L() { c.set_n_r(1, &c.l) } // CD
func (c *CPU) opSET_1_A() { c.set_n_r(1, &c.a) } // CF
func (c *CPU) opSET_2_B() { c.set_n_r(2, &c.b) } // D0
func (c *CPU) opSET_2_C() { c.set_n_r(2, &c.c) } // D1
func (c *CPU) opSET_2_D() { c.set_n_r(2, &c.d) } // D2
func (c *CPU) opSET_2_E() { c.set_n_r(2, &c.e) } // D3
func (c *CPU) opSET_2_H() { c.set_n_r(2, &c.h) } // D4
func (c *CPU) opSET_2_L() { c.set_n_r(2, &c.l) } // D5
func (c *CPU) opSET_2_A() { c.set_n_r(2, &c.a) } // D7
func (c *CPU) opSET_3_B() { c.set_n_r(3, &c.b) } // D8
func (c *CPU) opSET_3_C() { c.set_n_r(3, &c.c) } // D9
func (c *CPU) opSET_3_D() { c.set_n_r(3, &c.d) } // DA
func (c *CPU) opSET_3_E() { c.set_n_r(3, &c.e) } // DB
func (c *CPU) opSET_3_H() { c.set_n_r(3, &c.h) } // DC
func (c *CPU) opSET_3_L() { c.set_n_r(3, &c.l) } // DD
func (c *CPU) opSET_3_A() { c.set_n_r(3, &c.a) } // DF
func (c *CPU) opSET_4_B() { c.set_n_r(4, &c.b) } // E0
func (c *CPU) opSET_4_C() { c.set_n_r(4, &c.c) } // E1
func (c *CPU) opSET_4_D() { c.set_n_r(4, &c.d) } // E2
func (c *CPU) opSET_4_E() { c.set_n_r(4, &c.e) } // E3
func (c *CPU) opSET_4_H() { c.set_n_r(4, &c.h) } // E4
func (c *CPU) opSET_4_L() { c.set_n_r(4, &c.l) } // E5
func (c *CPU) opSET_4_A() { c.set_n_r(4, &c.a) } // E7
func (c *CPU) opSET_5_B() { c.set_n_r(5, &c.b) } // E8
func (c *CPU) opSET_5_C() { c.set_n_r(5, &c.c) } // E9
func (c *CPU) opSET_5_D() { c.set_n_r(5, &c.d) } // EA
func (c *CPU) opSET_5_E() { c.set_n_r(5, &c.e) } // EB
func (c *CPU) opSET_5_H() { c.set_n_r(5, &c.h) } // EC
func (c *CPU) opSET_5_L() { c.set_n_r(5, &c.l) } // ED
func (c *CPU) opSET_5_A() { c.set_n_r(5, &c.a) } // EF
func (c *CPU) opSET_6_B() { c.set_n_r(6, &c.b) } // F0
func (c *CPU) opSET_6_C() { c.set_n_r(6, &c.c) } // F1
func (c *CPU) opSET_6_D() { c.set_n_r(6, &c.d) } // F2
func (c *CPU) opSET_6_E() { c.set_n_r(6, &c.e) } // F3
func (c *CPU) opSET_6_H() { c.set_n_r(6, &c.h) } // F4
func (c *CPU) opSET_6_L() { c.set_n_r(6, &c.l) } // F5
func (c *CPU) opSET_6_A() { c.set_n_r(6, &c.a) } // F7
func (c *CPU) opSET_7_B() { c.set_n_r(7, &c.b) } // F8
func (c *CPU) opSET_7_C() { c.set_n_r(7, &c.c) } // F9
func (c *CPU) opSET_7_D() { c.set_n_r(7, &c.d) } // FA
func (c *CPU) opSET_7_E() { c.set_n_r(7, &c.e) } // FB
func (c *CPU) opSET_7_H() { c.set_n_r(7, &c.h) } // FC
func (c *CPU) opSET_7_L() { c.set_n_r(7, &c.l) } // FD
func (c *CPU) opSET_7_A() { c.set_n_r(7, &c.a) } // FF

// -------------------------------- SET n, [HL] ---------------------------------
func (c *CPU) set_n_ahl(n uint8) {
	addr := c.GetHL()
	val := c.read(addr)
	val = val | (1 << n)
	c.write(addr, val)
	c.cycles += 16
}
func (c *CPU) opSET_0_aHL() { c.set_n_ahl(0) } // C6
func (c *CPU) opSET_1_aHL() { c.set_n_ahl(1) } // CE
func (c *CPU) opSET_2_aHL() { c.set_n_ahl(2) } // D6
func (c *CPU) opSET_3_aHL() { c.set_n_ahl(3) } // DE
func (c *CPU) opSET_4_aHL() { c.set_n_ahl(4) } // E6
func (c *CPU) opSET_5_aHL() { c.set_n_ahl(5) } // EE
func (c *CPU) opSET_6_aHL() { c.set_n_ahl(6) } // F6
func (c *CPU) opSET_7_aHL() { c.set_n_ahl(7) } // FE
