package cpu

type OpEntry struct {
	fn   func(c *CPU)
	Name string
}

var OpTable = [256]OpEntry{
	0x00: {fn: func(c *CPU) { c.opNOP() }, Name: "NOP"},
	0x01: {fn: func(c *CPU) { c.opLD_BC_n16() }, Name: "LD BC, n16"},
	0x02: {fn: func(c *CPU) { c.opLD_aBC_A() }, Name: "LD [BC], A"},
	0x03: {fn: func(c *CPU) { c.opINC_BC() }, Name: "INC BC"},
	0x04: {fn: func(c *CPU) { c.opINC_B() }, Name: "INC B"},
	0x05: {fn: func(c *CPU) { c.opDEC_B() }, Name: "DEC B"},
	0x06: {fn: func(c *CPU) { c.opLD_B_n8() }, Name: "LD B, n8"},
	0x07: {fn: func(c *CPU) { c.opRLCA() }, Name: "RLCA"},
	0x08: {fn: func(c *CPU) { c.opLD_aa16_SP() }, Name: "LD [a16], SP"},
	0x09: {fn: func(c *CPU) { c.opADD_HL_BC() }, Name: "ADD HL, BC"},
	0x0A: {fn: func(c *CPU) { c.opLD_A_aBC() }, Name: "LD A, [BC]"},
	0x0B: {fn: func(c *CPU) { c.opDEC_BC() }, Name: "DEC BC"},
	0x0C: {fn: func(c *CPU) { c.opINC_C() }, Name: "INC C"},
	0x0D: {fn: func(c *CPU) { c.opDEC_C() }, Name: "DEC C"},
	0x0E: {fn: func(c *CPU) { c.opLD_C_n8() }, Name: "LD C, n8"},
	0x0F: {fn: func(c *CPU) { c.opRRCA() }, Name: "RRCA"},
	0x10: {fn: func(c *CPU) { c.opSTOP_n8() }, Name: "STOP n8"},
	0x11: {fn: func(c *CPU) { c.opLD_DE_n16() }, Name: "LD DE, n16"},
	0x12: {fn: func(c *CPU) { c.opLD_aDE_A() }, Name: "LD [DE], A"},
	0x13: {fn: func(c *CPU) { c.opINC_DE() }, Name: "INC DE"},
	0x14: {fn: func(c *CPU) { c.opINC_D() }, Name: "INC D"},
	0x15: {fn: func(c *CPU) { c.opDEC_D() }, Name: "DEC D"},
	0x16: {fn: func(c *CPU) { c.opLD_D_n8() }, Name: "LD D, n8"},
	0x17: {fn: func(c *CPU) { c.opRLA() }, Name: "RLA"},
	0x18: {fn: func(c *CPU) { c.opJR_e8() }, Name: "JR e8"},
	0x19: {fn: func(c *CPU) { c.opADD_HL_DE() }, Name: "ADD HL, DE"},
	0x1A: {fn: func(c *CPU) { c.opLD_A_aDE() }, Name: "LD A, [DE]"},
	0x1B: {fn: func(c *CPU) { c.opDEC_DE() }, Name: "DEC DE"},
	0x1C: {fn: func(c *CPU) { c.opINC_E() }, Name: "INC E"},
	0x1D: {fn: func(c *CPU) { c.opDEC_E() }, Name: "DEC E"},
	0x1E: {fn: func(c *CPU) { c.opLD_E_n8() }, Name: "LD E, n8"},
	0x1F: {fn: func(c *CPU) { c.opRRA() }, Name: "RRA"},
	0x20: {fn: func(c *CPU) { c.opJR_NZ_e8() }, Name: "JR NZ, e8"},
	0x21: {fn: func(c *CPU) { c.opLD_HL_n16() }, Name: "LD HL, n16"},
	0x22: {fn: func(c *CPU) { c.opLD_aHLp_A() }, Name: "LD [HL+], A"},
	0x23: {fn: func(c *CPU) { c.opINC_HL() }, Name: "INC HL"},
	0x24: {fn: func(c *CPU) { c.opINC_H() }, Name: "INC H"},
	0x25: {fn: func(c *CPU) { c.opDEC_H() }, Name: "DEC H"},
	0x26: {fn: func(c *CPU) { c.opLD_H_n8() }, Name: "LD H, n8"},
	0x27: {fn: func(c *CPU) { c.opDAA() }, Name: "DAA"},
	0x28: {fn: func(c *CPU) { c.opJR_Z_e8() }, Name: "JR Z, e8"},
	0x29: {fn: func(c *CPU) { c.opADD_HL_HL() }, Name: "ADD HL, HL"},
	0x2A: {fn: func(c *CPU) { c.opLD_A_aHLp() }, Name: "LD A, [HL+]"},
	0x2B: {fn: func(c *CPU) { c.opDEC_HL() }, Name: "DEC HL"},
	0x2C: {fn: func(c *CPU) { c.opINC_L() }, Name: "INC L"},
	0x2D: {fn: func(c *CPU) { c.opDEC_L() }, Name: "DEC L"},
	0x2E: {fn: func(c *CPU) { c.opLD_L_n8() }, Name: "LD L, n8"},
	0x2F: {fn: func(c *CPU) { c.opCPL() }, Name: "CPL"},
	0x30: {fn: func(c *CPU) { c.opJR_NC_e8() }, Name: "JR NC, e8"},
	0x31: {fn: func(c *CPU) { c.opLD_SP_n16() }, Name: "LD SP, n16"},
	0x32: {fn: func(c *CPU) { c.opLD_aHLm_A() }, Name: "LD [HL-], A"},
	0x33: {fn: func(c *CPU) { c.opINC_SP() }, Name: "INC SP"},
	0x34: {fn: func(c *CPU) { c.opINC_aHL() }, Name: "INC [HL]"},
	0x35: {fn: func(c *CPU) { c.opDEC_aHL() }, Name: "DEC [HL]"},
	0x36: {fn: func(c *CPU) { c.opLD_aHL_n8() }, Name: "LD [HL], n8"},
	0x37: {fn: func(c *CPU) { c.opSCF() }, Name: "SCF"},
	0x38: {fn: func(c *CPU) { c.opJR_C_e8() }, Name: "JR C, e8"},
	0x39: {fn: func(c *CPU) { c.opADD_HL_SP() }, Name: "ADD HL, SP"},
	0x3A: {fn: func(c *CPU) { c.opLD_A_aHLm() }, Name: "LD A, [HL-]"},
	0x3B: {fn: func(c *CPU) { c.opDEC_SP() }, Name: "DEC SP"},
	0x3C: {fn: func(c *CPU) { c.opINC_A() }, Name: "INC A"},
	0x3D: {fn: func(c *CPU) { c.opDEC_A() }, Name: "DEC A"},
	0x3E: {fn: func(c *CPU) { c.opLD_A_n8() }, Name: "LD A, n8"},
	0x3F: {fn: func(c *CPU) { c.opCCF() }, Name: "CCF"},
	0x40: {fn: func(c *CPU) { c.opLD_B_B() }, Name: "LD B, B"},
	0x41: {fn: func(c *CPU) { c.opLD_B_C() }, Name: "LD B, C"},
	0x42: {fn: func(c *CPU) { c.opLD_B_D() }, Name: "LD B, D"},
	0x43: {fn: func(c *CPU) { c.opLD_B_E() }, Name: "LD B, E"},
	0x44: {fn: func(c *CPU) { c.opLD_B_H() }, Name: "LD B, H"},
	0x45: {fn: func(c *CPU) { c.opLD_B_L() }, Name: "LD B, L"},
	0x46: {fn: func(c *CPU) { c.opLD_B_aHL() }, Name: "LD B, [HL]"},
	0x47: {fn: func(c *CPU) { c.opLD_B_A() }, Name: "LD B, A"},
	0x48: {fn: func(c *CPU) { c.opLD_C_B() }, Name: "LD C, B"},
	0x49: {fn: func(c *CPU) { c.opLD_C_C() }, Name: "LD C, C"},
	0x4A: {fn: func(c *CPU) { c.opLD_C_D() }, Name: "LD C, D"},
	0x4B: {fn: func(c *CPU) { c.opLD_C_E() }, Name: "LD C, E"},
	0x4C: {fn: func(c *CPU) { c.opLD_C_H() }, Name: "LD C, H"},
	0x4D: {fn: func(c *CPU) { c.opLD_C_L() }, Name: "LD C, L"},
	0x4E: {fn: func(c *CPU) { c.opLD_C_aHL() }, Name: "LD C, [HL]"},
	0x4F: {fn: func(c *CPU) { c.opLD_C_A() }, Name: "LD C, A"},
	0x50: {fn: func(c *CPU) { c.opLD_D_B() }, Name: "LD D, B"},
	0x51: {fn: func(c *CPU) { c.opLD_D_C() }, Name: "LD D, C"},
	0x52: {fn: func(c *CPU) { c.opLD_D_D() }, Name: "LD D, D"},
	0x53: {fn: func(c *CPU) { c.opLD_D_E() }, Name: "LD D, E"},
	0x54: {fn: func(c *CPU) { c.opLD_D_H() }, Name: "LD D, H"},
	0x55: {fn: func(c *CPU) { c.opLD_D_L() }, Name: "LD D, L"},
	0x56: {fn: func(c *CPU) { c.opLD_D_aHL() }, Name: "LD D, [HL]"},
	0x57: {fn: func(c *CPU) { c.opLD_D_A() }, Name: "LD D, A"},
	0x58: {fn: func(c *CPU) { c.opLD_E_B() }, Name: "LD E, B"},
	0x59: {fn: func(c *CPU) { c.opLD_E_C() }, Name: "LD E, C"},
	0x5A: {fn: func(c *CPU) { c.opLD_E_D() }, Name: "LD E, D"},
	0x5B: {fn: func(c *CPU) { c.opLD_E_E() }, Name: "LD E, E"},
	0x5C: {fn: func(c *CPU) { c.opLD_E_H() }, Name: "LD E, H"},
	0x5D: {fn: func(c *CPU) { c.opLD_E_L() }, Name: "LD E, L"},
	0x5E: {fn: func(c *CPU) { c.opLD_E_aHL() }, Name: "LD E, [HL]"},
	0x5F: {fn: func(c *CPU) { c.opLD_E_A() }, Name: "LD E, A"},
	0x60: {fn: func(c *CPU) { c.opLD_H_B() }, Name: "LD H, B"},
	0x61: {fn: func(c *CPU) { c.opLD_H_C() }, Name: "LD H, C"},
	0x62: {fn: func(c *CPU) { c.opLD_H_D() }, Name: "LD H, D"},
	0x63: {fn: func(c *CPU) { c.opLD_H_E() }, Name: "LD H, E"},
	0x64: {fn: func(c *CPU) { c.opLD_H_H() }, Name: "LD H, H"},
	0x65: {fn: func(c *CPU) { c.opLD_H_L() }, Name: "LD H, L"},
	0x66: {fn: func(c *CPU) { c.opLD_H_aHL() }, Name: "LD H, [HL]"},
	0x67: {fn: func(c *CPU) { c.opLD_H_A() }, Name: "LD H, A"},
	0x68: {fn: func(c *CPU) { c.opLD_L_B() }, Name: "LD L, B"},
	0x69: {fn: func(c *CPU) { c.opLD_L_C() }, Name: "LD L, C"},
	0x6A: {fn: func(c *CPU) { c.opLD_L_D() }, Name: "LD L, D"},
	0x6B: {fn: func(c *CPU) { c.opLD_L_E() }, Name: "LD L, E"},
	0x6C: {fn: func(c *CPU) { c.opLD_L_H() }, Name: "LD L, H"},
	0x6D: {fn: func(c *CPU) { c.opLD_L_L() }, Name: "LD L, L"},
	0x6E: {fn: func(c *CPU) { c.opLD_L_aHL() }, Name: "LD L, [HL]"},
	0x6F: {fn: func(c *CPU) { c.opLD_L_A() }, Name: "LD L, A"},
	0x70: {fn: func(c *CPU) { c.opLD_aHL_B() }, Name: "LD [HL], B"},
	0x71: {fn: func(c *CPU) { c.opLD_aHL_C() }, Name: "LD [HL], C"},
	0x72: {fn: func(c *CPU) { c.opLD_aHL_D() }, Name: "LD [HL], D"},
	0x73: {fn: func(c *CPU) { c.opLD_aHL_E() }, Name: "LD [HL], E"},
	0x74: {fn: func(c *CPU) { c.opLD_aHL_H() }, Name: "LD [HL], H"},
	0x75: {fn: func(c *CPU) { c.opLD_aHL_L() }, Name: "LD [HL], L"},
	0x76: {fn: func(c *CPU) { c.opHALT() }, Name: "HALT"},
	0x77: {fn: func(c *CPU) { c.opLD_aHL_A() }, Name: "LD [HL], A"},
	0x78: {fn: func(c *CPU) { c.opLD_A_B() }, Name: "LD A, B"},
	0x79: {fn: func(c *CPU) { c.opLD_A_C() }, Name: "LD A, C"},
	0x7A: {fn: func(c *CPU) { c.opLD_A_D() }, Name: "LD A, D"},
	0x7B: {fn: func(c *CPU) { c.opLD_A_E() }, Name: "LD A, E"},
	0x7C: {fn: func(c *CPU) { c.opLD_A_H() }, Name: "LD A, H"},
	0x7D: {fn: func(c *CPU) { c.opLD_A_L() }, Name: "LD A, L"},
	0x7E: {fn: func(c *CPU) { c.opLD_A_aHL() }, Name: "LD A, [HL]"},
	0x7F: {fn: func(c *CPU) { c.opLD_A_A() }, Name: "LD A, A"},
	0x80: {fn: func(c *CPU) { c.opADD_B() }, Name: "ADD A, B"},
	0x81: {fn: func(c *CPU) { c.opADD_C() }, Name: "ADD A, C"},
	0x82: {fn: func(c *CPU) { c.opADD_D() }, Name: "ADD A, D"},
	0x83: {fn: func(c *CPU) { c.opADD_E() }, Name: "ADD A, E"},
	0x84: {fn: func(c *CPU) { c.opADD_H() }, Name: "ADD A, H"},
	0x85: {fn: func(c *CPU) { c.opADD_L() }, Name: "ADD A, L"},
	0x86: {fn: func(c *CPU) { c.opADD_A_aHL() }, Name: "ADD A, [HL]"},
	0x87: {fn: func(c *CPU) { c.opADD_A() }, Name: "ADD A, A"},
	0x88: {fn: func(c *CPU) { c.opADC_B() }, Name: "ADC A, B"},
	0x89: {fn: func(c *CPU) { c.opADC_C() }, Name: "ADC A, C"},
	0x8A: {fn: func(c *CPU) { c.opADC_D() }, Name: "ADC A, D"},
	0x8B: {fn: func(c *CPU) { c.opADC_E() }, Name: "ADC A, E"},
	0x8C: {fn: func(c *CPU) { c.opADC_H() }, Name: "ADC A, H"},
	0x8D: {fn: func(c *CPU) { c.opADC_L() }, Name: "ADC A, L"},
	0x8E: {fn: func(c *CPU) { c.opADC_A_aHL() }, Name: "ADC A, [HL]"},
	0x8F: {fn: func(c *CPU) { c.opADC_A() }, Name: "ADC A, A"},
	0x90: {fn: func(c *CPU) { c.opSUB_B() }, Name: "SUB A, B"},
	0x91: {fn: func(c *CPU) { c.opSUB_C() }, Name: "SUB A, C"},
	0x92: {fn: func(c *CPU) { c.opSUB_D() }, Name: "SUB A, D"},
	0x93: {fn: func(c *CPU) { c.opSUB_E() }, Name: "SUB A, E"},
	0x94: {fn: func(c *CPU) { c.opSUB_H() }, Name: "SUB A, H"},
	0x95: {fn: func(c *CPU) { c.opSUB_L() }, Name: "SUB A, L"},
	0x96: {fn: func(c *CPU) { c.opSUB_A_aHL() }, Name: "SUB A, [HL]"},
	0x97: {fn: func(c *CPU) { c.opSUB_A() }, Name: "SUB A, A"},
	0x98: {fn: func(c *CPU) { c.opSBC_B() }, Name: "SBC A, B"},
	0x99: {fn: func(c *CPU) { c.opSBC_C() }, Name: "SBC A, C"},
	0x9A: {fn: func(c *CPU) { c.opSBC_D() }, Name: "SBC A, D"},
	0x9B: {fn: func(c *CPU) { c.opSBC_E() }, Name: "SBC A, E"},
	0x9C: {fn: func(c *CPU) { c.opSBC_H() }, Name: "SBC A, H"},
	0x9D: {fn: func(c *CPU) { c.opSBC_L() }, Name: "SBC A, L"},
	0x9E: {fn: func(c *CPU) { c.opSBC_A_aHL() }, Name: "SBC A, [HL]"},
	0x9F: {fn: func(c *CPU) { c.opSBC_A() }, Name: "SBC A, A"},
	0xA0: {fn: func(c *CPU) { c.opAND_B() }, Name: "AND A, B"},
	0xA1: {fn: func(c *CPU) { c.opAND_C() }, Name: "AND A, C"},
	0xA2: {fn: func(c *CPU) { c.opAND_D() }, Name: "AND A, D"},
	0xA3: {fn: func(c *CPU) { c.opAND_E() }, Name: "AND A, E"},
	0xA4: {fn: func(c *CPU) { c.opAND_H() }, Name: "AND A, H"},
	0xA5: {fn: func(c *CPU) { c.opAND_L() }, Name: "AND A, L"},
	0xA6: {fn: func(c *CPU) { c.opAND_A_aHL() }, Name: "AND A, [HL]"},
	0xA7: {fn: func(c *CPU) { c.opAND_A() }, Name: "AND A, A"},
	0xA8: {fn: func(c *CPU) { c.opXOR_B() }, Name: "XOR A, B"},
	0xA9: {fn: func(c *CPU) { c.opXOR_C() }, Name: "XOR A, C"},
	0xAA: {fn: func(c *CPU) { c.opXOR_D() }, Name: "XOR A, D"},
	0xAB: {fn: func(c *CPU) { c.opXOR_E() }, Name: "XOR A, E"},
	0xAC: {fn: func(c *CPU) { c.opXOR_H() }, Name: "XOR A, H"},
	0xAD: {fn: func(c *CPU) { c.opXOR_L() }, Name: "XOR A, L"},
	0xAE: {fn: func(c *CPU) { c.opXOR_A_aHL() }, Name: "XOR A, [HL]"},
	0xAF: {fn: func(c *CPU) { c.opXOR_A() }, Name: "XOR A, A"},
	0xB0: {fn: func(c *CPU) { c.opOR_B() }, Name: "OR A, B"},
	0xB1: {fn: func(c *CPU) { c.opOR_C() }, Name: "OR A, C"},
	0xB2: {fn: func(c *CPU) { c.opOR_D() }, Name: "OR A, D"},
	0xB3: {fn: func(c *CPU) { c.opOR_E() }, Name: "OR A, E"},
	0xB4: {fn: func(c *CPU) { c.opOR_H() }, Name: "OR A, H"},
	0xB5: {fn: func(c *CPU) { c.opOR_L() }, Name: "OR A, L"},
	0xB6: {fn: func(c *CPU) { c.opOR_A_aHL() }, Name: "OR A, [HL]"},
	0xB7: {fn: func(c *CPU) { c.opOR_A() }, Name: "OR A, A"},
	0xB8: {fn: func(c *CPU) { c.opCP_B() }, Name: "CP A, B"},
	0xB9: {fn: func(c *CPU) { c.opCP_C() }, Name: "CP A, C"},
	0xBA: {fn: func(c *CPU) { c.opCP_D() }, Name: "CP A, D"},
	0xBB: {fn: func(c *CPU) { c.opCP_E() }, Name: "CP A, E"},
	0xBC: {fn: func(c *CPU) { c.opCP_H() }, Name: "CP A, H"},
	0xBD: {fn: func(c *CPU) { c.opCP_L() }, Name: "CP A, L"},
	0xBE: {fn: func(c *CPU) { c.opCP_A_aHL() }, Name: "CP A, [HL]"},
	0xBF: {fn: func(c *CPU) { c.opCP_A() }, Name: "CP A, A"},
	0xC0: {fn: func(c *CPU) { c.opRET_NZ() }, Name: "RET NZ"},
	0xC1: {fn: func(c *CPU) { c.opPOP_BC() }, Name: "POP BC"},
	0xC2: {fn: func(c *CPU) { c.opJP_NZ_a16() }, Name: "JP NZ, a16"},
	0xC3: {fn: func(c *CPU) { c.opJP_a16() }, Name: "JP a16"},
	0xC4: {fn: func(c *CPU) { c.opCALL_NZ_a16() }, Name: "CALL NZ, a16"},
	0xC5: {fn: func(c *CPU) { c.opPUSH_BC() }, Name: "PUSH BC"},
	0xC6: {fn: func(c *CPU) { c.opADD_A_n8() }, Name: "ADD A, n8"},
	0xC7: {fn: func(c *CPU) { c.opRST_00() }, Name: "RST $00"},
	0xC8: {fn: func(c *CPU) { c.opRET_Z() }, Name: "RET Z"},
	0xC9: {fn: func(c *CPU) { c.opRET() }, Name: "RET"},
	0xCA: {fn: func(c *CPU) { c.opJP_Z_a16() }, Name: "JP Z, a16"},
	0xCB: {fn: func(c *CPU) { c.opPREFIX() }, Name: "PREFIX"},
	0xCC: {fn: func(c *CPU) { c.opCALL_Z_a16() }, Name: "CALL Z, a16"},
	0xCD: {fn: func(c *CPU) { c.opCALL_a16() }, Name: "CALL a16"},
	0xCE: {fn: func(c *CPU) { c.opADC_A_n8() }, Name: "ADC A, n8"},
	0xCF: {fn: func(c *CPU) { c.opRST_08() }, Name: "RST $08"},
	0xD0: {fn: func(c *CPU) { c.opRET_NC() }, Name: "RET NC"},
	0xD1: {fn: func(c *CPU) { c.opPOP_DE() }, Name: "POP DE"},
	0xD2: {fn: func(c *CPU) { c.opJP_NC_a16() }, Name: "JP NC, a16"},
	0xD3: {fn: func(c *CPU) { c.opNOP() }, Name: ""}, // 未定義
	0xD4: {fn: func(c *CPU) { c.opCALL_NC_a16() }, Name: "CALL NC, a16"},
	0xD5: {fn: func(c *CPU) { c.opPUSH_DE() }, Name: "PUSH DE"},
	0xD6: {fn: func(c *CPU) { c.opSUB_A_n8() }, Name: "SUB A, n8"},
	0xD7: {fn: func(c *CPU) { c.opRST_10() }, Name: "RST $10"},
	0xD8: {fn: func(c *CPU) { c.opRET_C() }, Name: "RET C"},
	0xD9: {fn: func(c *CPU) { c.opRETI() }, Name: "RETI"},
	0xDA: {fn: func(c *CPU) { c.opJP_C_a16() }, Name: "JP C, a16"},
	0xDB: {fn: func(c *CPU) { c.opNOP() }, Name: ""}, // 未定義
	0xDC: {fn: func(c *CPU) { c.opCALL_C_a16() }, Name: "CALL C, a16"},
	0xDD: {fn: func(c *CPU) { c.opNOP() }, Name: ""}, // 未定義
	0xDE: {fn: func(c *CPU) { c.opSBC_A_n8() }, Name: "SBC A, n8"},
	0xDF: {fn: func(c *CPU) { c.opRST_18() }, Name: "RST $18"},
	0xE0: {fn: func(c *CPU) { c.opLDH_aa8_A() }, Name: "LDH [a8], A"},
	0xE1: {fn: func(c *CPU) { c.opPOP_HL() }, Name: "POP HL"},
	0xE2: {fn: func(c *CPU) { c.opLDH_aC_A() }, Name: "LDH [C], A"},
	0xE3: {fn: func(c *CPU) { c.opNOP() }, Name: ""}, // 未定義
	0xE4: {fn: func(c *CPU) { c.opNOP() }, Name: ""}, // 未定義
	0xE5: {fn: func(c *CPU) { c.opPUSH_HL() }, Name: "PUSH HL"},
	0xE6: {fn: func(c *CPU) { c.opAND_A_n8() }, Name: "AND A, n8"},
	0xE7: {fn: func(c *CPU) { c.opRST_20() }, Name: "RST $20"},
	0xE8: {fn: func(c *CPU) { c.opADD_SP_e8() }, Name: "ADD SP, e8"},
	0xE9: {fn: func(c *CPU) { c.opJP_HL() }, Name: "JP HL"},
	0xEA: {fn: func(c *CPU) { c.opLD_aa16_A() }, Name: "LD [a16], A"},
	0xEB: {fn: func(c *CPU) { c.opNOP() }, Name: ""}, // 未定義
	0xEC: {fn: func(c *CPU) { c.opNOP() }, Name: ""}, // 未定義
	0xED: {fn: func(c *CPU) { c.opNOP() }, Name: ""}, // 未定義
	0xEE: {fn: func(c *CPU) { c.opXOR_A_n8() }, Name: "XOR A, n8"},
	0xEF: {fn: func(c *CPU) { c.opRST_28() }, Name: "RST $28"},
	0xF0: {fn: func(c *CPU) { c.opLDH_A_aa8() }, Name: "LDH A, [a8]"},
	0xF1: {fn: func(c *CPU) { c.opPOP_AF() }, Name: "POP AF"},
	0xF2: {fn: func(c *CPU) { c.opLDH_A_aC() }, Name: "LDH A, [C]"},
	0xF3: {fn: func(c *CPU) { c.opDI() }, Name: "DI"},
	0xF4: {fn: func(c *CPU) { c.opNOP() }, Name: ""}, // 未定義
	0xF5: {fn: func(c *CPU) { c.opPUSH_AF() }, Name: "PUSH AF"},
	0xF6: {fn: func(c *CPU) { c.opOR_A_n8() }, Name: "OR A, n8"},
	0xF7: {fn: func(c *CPU) { c.opRST_30() }, Name: "RST $30"},
	0xF8: {fn: func(c *CPU) { c.opLD_HL_SP_p_e8() }, Name: "LD HL, SP + e8"},
	0xF9: {fn: func(c *CPU) { c.opLD_SP_HL() }, Name: "LD SP, HL"},
	0xFA: {fn: func(c *CPU) { c.opLD_A_aa16() }, Name: "LD A, [a16]"},
	0xFB: {fn: func(c *CPU) { c.opEI() }, Name: "EI"},
	0xFC: {fn: func(c *CPU) { c.opNOP() }, Name: ""}, // 未定義
	0xFD: {fn: func(c *CPU) { c.opNOP() }, Name: ""}, // 未定義
	0xFE: {fn: func(c *CPU) { c.opCP_A_n8() }, Name: "CP A, n8"},
	0xFF: {fn: func(c *CPU) { c.opRST_38() }, Name: "RST $38"},
}

var cbTable = [256]OpEntry{
	0x00: {fn: func(c *CPU) { c.opRLC_B() }, Name: "RLC B"},           // 00
	0x01: {fn: func(c *CPU) { c.opRLC_C() }, Name: "RLC C"},           // 01
	0x02: {fn: func(c *CPU) { c.opRLC_D() }, Name: "RLC D"},           // 02
	0x03: {fn: func(c *CPU) { c.opRLC_E() }, Name: "RLC E"},           // 03
	0x04: {fn: func(c *CPU) { c.opRLC_H() }, Name: "RLC H"},           // 04
	0x05: {fn: func(c *CPU) { c.opRLC_L() }, Name: "RLC L"},           // 05
	0x06: {fn: func(c *CPU) { c.opRLC_aHL() }, Name: "RLC [HL]"},      // 06
	0x07: {fn: func(c *CPU) { c.opRLC_A() }, Name: "RLC A"},           // 07
	0x08: {fn: func(c *CPU) { c.opRRC_B() }, Name: "RRC B"},           // 08
	0x09: {fn: func(c *CPU) { c.opRRC_C() }, Name: "RRC C"},           // 09
	0x0A: {fn: func(c *CPU) { c.opRRC_D() }, Name: "RRC D"},           // 0A
	0x0B: {fn: func(c *CPU) { c.opRRC_E() }, Name: "RRC E"},           // 0B
	0x0C: {fn: func(c *CPU) { c.opRRC_H() }, Name: "RRC H"},           // 0C
	0x0D: {fn: func(c *CPU) { c.opRRC_L() }, Name: "RRC L"},           // 0D
	0x0E: {fn: func(c *CPU) { c.opRRC_aHL() }, Name: "RRC [HL]"},      // 0E
	0x0F: {fn: func(c *CPU) { c.opRRC_A() }, Name: "RRC A"},           // 0F
	0x10: {fn: func(c *CPU) { c.opRL_B() }, Name: "RL B"},             // 10
	0x11: {fn: func(c *CPU) { c.opRL_C() }, Name: "RL C"},             // 11
	0x12: {fn: func(c *CPU) { c.opRL_D() }, Name: "RL D"},             // 12
	0x13: {fn: func(c *CPU) { c.opRL_E() }, Name: "RL E"},             // 13
	0x14: {fn: func(c *CPU) { c.opRL_H() }, Name: "RL H"},             // 14
	0x15: {fn: func(c *CPU) { c.opRL_L() }, Name: "RL L"},             // 15
	0x16: {fn: func(c *CPU) { c.opRL_aHL() }, Name: "RL [HL]"},        // 16
	0x17: {fn: func(c *CPU) { c.opRL_A() }, Name: "RL A"},             // 17
	0x18: {fn: func(c *CPU) { c.opRR_B() }, Name: "RR B"},             // 18
	0x19: {fn: func(c *CPU) { c.opRR_C() }, Name: "RR C"},             // 19
	0x1A: {fn: func(c *CPU) { c.opRR_D() }, Name: "RR D"},             // 1A
	0x1B: {fn: func(c *CPU) { c.opRR_E() }, Name: "RR E"},             // 1B
	0x1C: {fn: func(c *CPU) { c.opRR_H() }, Name: "RR H"},             // 1C
	0x1D: {fn: func(c *CPU) { c.opRR_L() }, Name: "RR L"},             // 1D
	0x1E: {fn: func(c *CPU) { c.opRR_aHL() }, Name: "RR [HL]"},        // 1E
	0x1F: {fn: func(c *CPU) { c.opRR_A() }, Name: "RR A"},             // 1F
	0x20: {fn: func(c *CPU) { c.opSLA_B() }, Name: "SLA B"},           // 20
	0x21: {fn: func(c *CPU) { c.opSLA_C() }, Name: "SLA C"},           // 21
	0x22: {fn: func(c *CPU) { c.opSLA_D() }, Name: "SLA D"},           // 22
	0x23: {fn: func(c *CPU) { c.opSLA_E() }, Name: "SLA E"},           // 23
	0x24: {fn: func(c *CPU) { c.opSLA_H() }, Name: "SLA H"},           // 24
	0x25: {fn: func(c *CPU) { c.opSLA_L() }, Name: "SLA L"},           // 25
	0x26: {fn: func(c *CPU) { c.opSLA_aHL() }, Name: "SLA [HL]"},      // 26
	0x27: {fn: func(c *CPU) { c.opSLA_A() }, Name: "SLA A"},           // 27
	0x28: {fn: func(c *CPU) { c.opSRA_B() }, Name: "SRA B"},           // 28
	0x29: {fn: func(c *CPU) { c.opSRA_C() }, Name: "SRA C"},           // 29
	0x2A: {fn: func(c *CPU) { c.opSRA_D() }, Name: "SRA D"},           // 2A
	0x2B: {fn: func(c *CPU) { c.opSRA_E() }, Name: "SRA E"},           // 2B
	0x2C: {fn: func(c *CPU) { c.opSRA_H() }, Name: "SRA H"},           // 2C
	0x2D: {fn: func(c *CPU) { c.opSRA_L() }, Name: "SRA L"},           // 2D
	0x2E: {fn: func(c *CPU) { c.opSRA_aHL() }, Name: "SRA [HL]"},      // 2E
	0x2F: {fn: func(c *CPU) { c.opSRA_A() }, Name: "SRA A"},           // 2F
	0x30: {fn: func(c *CPU) { c.opSWAP_B() }, Name: "SWAP B"},         // 30
	0x31: {fn: func(c *CPU) { c.opSWAP_C() }, Name: "SWAP C"},         // 31
	0x32: {fn: func(c *CPU) { c.opSWAP_D() }, Name: "SWAP D"},         // 32
	0x33: {fn: func(c *CPU) { c.opSWAP_E() }, Name: "SWAP E"},         // 33
	0x34: {fn: func(c *CPU) { c.opSWAP_H() }, Name: "SWAP H"},         // 34
	0x35: {fn: func(c *CPU) { c.opSWAP_L() }, Name: "SWAP L"},         // 35
	0x36: {fn: func(c *CPU) { c.opSWAP_aHL() }, Name: "SWAP [HL]"},    // 36
	0x37: {fn: func(c *CPU) { c.opSWAP_A() }, Name: "SWAP A"},         // 37
	0x38: {fn: func(c *CPU) { c.opSRL_B() }, Name: "SRL B"},           // 38
	0x39: {fn: func(c *CPU) { c.opSRL_C() }, Name: "SRL C"},           // 39
	0x3A: {fn: func(c *CPU) { c.opSRL_D() }, Name: "SRL D"},           // 3A
	0x3B: {fn: func(c *CPU) { c.opSRL_E() }, Name: "SRL E"},           // 3B
	0x3C: {fn: func(c *CPU) { c.opSRL_H() }, Name: "SRL H"},           // 3C
	0x3D: {fn: func(c *CPU) { c.opSRL_L() }, Name: "SRL L"},           // 3D
	0x3E: {fn: func(c *CPU) { c.opSRL_aHL() }, Name: "SRL [HL]"},      // 3E
	0x3F: {fn: func(c *CPU) { c.opSRL_A() }, Name: "SRL A"},           // 3F
	0x40: {fn: func(c *CPU) { c.opBIT_0_B() }, Name: "BIT 0, B"},      // 40
	0x41: {fn: func(c *CPU) { c.opBIT_0_C() }, Name: "BIT 0, C"},      // 41
	0x42: {fn: func(c *CPU) { c.opBIT_0_D() }, Name: "BIT 0, D"},      // 42
	0x43: {fn: func(c *CPU) { c.opBIT_0_E() }, Name: "BIT 0, E"},      // 43
	0x44: {fn: func(c *CPU) { c.opBIT_0_H() }, Name: "BIT 0, H"},      // 44
	0x45: {fn: func(c *CPU) { c.opBIT_0_L() }, Name: "BIT 0, L"},      // 45
	0x46: {fn: func(c *CPU) { c.opBIT_0_aHL() }, Name: "BIT 0, [HL]"}, // 46
	0x47: {fn: func(c *CPU) { c.opBIT_0_A() }, Name: "BIT 0, A"},      // 47
	0x48: {fn: func(c *CPU) { c.opBIT_1_B() }, Name: "BIT 1, B"},      // 48
	0x49: {fn: func(c *CPU) { c.opBIT_1_C() }, Name: "BIT 1, C"},      // 49
	0x4A: {fn: func(c *CPU) { c.opBIT_1_D() }, Name: "BIT 1, D"},      // 4A
	0x4B: {fn: func(c *CPU) { c.opBIT_1_E() }, Name: "BIT 1, E"},      // 4B
	0x4C: {fn: func(c *CPU) { c.opBIT_1_H() }, Name: "BIT 1, H"},      // 4C
	0x4D: {fn: func(c *CPU) { c.opBIT_1_L() }, Name: "BIT 1, L"},      // 4D
	0x4E: {fn: func(c *CPU) { c.opBIT_1_aHL() }, Name: "BIT 1, [HL]"}, // 4E
	0x4F: {fn: func(c *CPU) { c.opBIT_1_A() }, Name: "BIT 1, A"},      // 4F
	0x50: {fn: func(c *CPU) { c.opBIT_2_B() }, Name: "BIT 2, B"},      // 50
	0x51: {fn: func(c *CPU) { c.opBIT_2_C() }, Name: "BIT 2, C"},      // 51
	0x52: {fn: func(c *CPU) { c.opBIT_2_D() }, Name: "BIT 2, D"},      // 52
	0x53: {fn: func(c *CPU) { c.opBIT_2_E() }, Name: "BIT 2, E"},      // 53
	0x54: {fn: func(c *CPU) { c.opBIT_2_H() }, Name: "BIT 2, H"},      // 54
	0x55: {fn: func(c *CPU) { c.opBIT_2_L() }, Name: "BIT 2, L"},      // 55
	0x56: {fn: func(c *CPU) { c.opBIT_2_aHL() }, Name: "BIT 2, [HL]"}, // 56
	0x57: {fn: func(c *CPU) { c.opBIT_2_A() }, Name: "BIT 2, A"},      // 57
	0x58: {fn: func(c *CPU) { c.opBIT_3_B() }, Name: "BIT 3, B"},      // 58
	0x59: {fn: func(c *CPU) { c.opBIT_3_C() }, Name: "BIT 3, C"},      // 59
	0x5A: {fn: func(c *CPU) { c.opBIT_3_D() }, Name: "BIT 3, D"},      // 5A
	0x5B: {fn: func(c *CPU) { c.opBIT_3_E() }, Name: "BIT 3, E"},      // 5B
	0x5C: {fn: func(c *CPU) { c.opBIT_3_H() }, Name: "BIT 3, H"},      // 5C
	0x5D: {fn: func(c *CPU) { c.opBIT_3_L() }, Name: "BIT 3, L"},      // 5D
	0x5E: {fn: func(c *CPU) { c.opBIT_3_aHL() }, Name: "BIT 3, [HL]"}, // 5E
	0x5F: {fn: func(c *CPU) { c.opBIT_3_A() }, Name: "BIT 3, A"},      // 5F
	0x60: {fn: func(c *CPU) { c.opBIT_4_B() }, Name: "BIT 4, B"},      // 60
	0x61: {fn: func(c *CPU) { c.opBIT_4_C() }, Name: "BIT 4, C"},      // 61
	0x62: {fn: func(c *CPU) { c.opBIT_4_D() }, Name: "BIT 4, D"},      // 62
	0x63: {fn: func(c *CPU) { c.opBIT_4_E() }, Name: "BIT 4, E"},      // 63
	0x64: {fn: func(c *CPU) { c.opBIT_4_H() }, Name: "BIT 4, H"},      // 64
	0x65: {fn: func(c *CPU) { c.opBIT_4_L() }, Name: "BIT 4, L"},      // 65
	0x66: {fn: func(c *CPU) { c.opBIT_4_aHL() }, Name: "BIT 4, [HL]"}, // 66
	0x67: {fn: func(c *CPU) { c.opBIT_4_A() }, Name: "BIT 4, A"},      // 67
	0x68: {fn: func(c *CPU) { c.opBIT_5_B() }, Name: "BIT 5, B"},      // 68
	0x69: {fn: func(c *CPU) { c.opBIT_5_C() }, Name: "BIT 5, C"},      // 69
	0x6A: {fn: func(c *CPU) { c.opBIT_5_D() }, Name: "BIT 5, D"},      // 6A
	0x6B: {fn: func(c *CPU) { c.opBIT_5_E() }, Name: "BIT 5, E"},      // 6B
	0x6C: {fn: func(c *CPU) { c.opBIT_5_H() }, Name: "BIT 5, H"},      // 6C
	0x6D: {fn: func(c *CPU) { c.opBIT_5_L() }, Name: "BIT 5, L"},      // 6D
	0x6E: {fn: func(c *CPU) { c.opBIT_5_aHL() }, Name: "BIT 5, [HL]"}, // 6E
	0x6F: {fn: func(c *CPU) { c.opBIT_5_A() }, Name: "BIT 5, A"},      // 6F
	0x70: {fn: func(c *CPU) { c.opBIT_6_B() }, Name: "BIT 6, B"},      // 70
	0x71: {fn: func(c *CPU) { c.opBIT_6_C() }, Name: "BIT 6, C"},      // 71
	0x72: {fn: func(c *CPU) { c.opBIT_6_D() }, Name: "BIT 6, D"},      // 72
	0x73: {fn: func(c *CPU) { c.opBIT_6_E() }, Name: "BIT 6, E"},      // 73
	0x74: {fn: func(c *CPU) { c.opBIT_6_H() }, Name: "BIT 6, H"},      // 74
	0x75: {fn: func(c *CPU) { c.opBIT_6_L() }, Name: "BIT 6, L"},      // 75
	0x76: {fn: func(c *CPU) { c.opBIT_6_aHL() }, Name: "BIT 6, [HL]"}, // 76
	0x77: {fn: func(c *CPU) { c.opBIT_6_A() }, Name: "BIT 6, A"},      // 77
	0x78: {fn: func(c *CPU) { c.opBIT_7_B() }, Name: "BIT 7, B"},      // 78
	0x79: {fn: func(c *CPU) { c.opBIT_7_C() }, Name: "BIT 7, C"},      // 79
	0x7A: {fn: func(c *CPU) { c.opBIT_7_D() }, Name: "BIT 7, D"},      // 7A
	0x7B: {fn: func(c *CPU) { c.opBIT_7_E() }, Name: "BIT 7, E"},      // 7B
	0x7C: {fn: func(c *CPU) { c.opBIT_7_H() }, Name: "BIT 7, H"},      // 7C
	0x7D: {fn: func(c *CPU) { c.opBIT_7_L() }, Name: "BIT 7, L"},      // 7D
	0x7E: {fn: func(c *CPU) { c.opBIT_7_aHL() }, Name: "BIT 7, [HL]"}, // 7E
	0x7F: {fn: func(c *CPU) { c.opBIT_7_A() }, Name: "BIT 7, A"},      // 7F
	0x80: {fn: func(c *CPU) { c.opRES_0_B() }, Name: "RES 0, B"},      // 80
	0x81: {fn: func(c *CPU) { c.opRES_0_C() }, Name: "RES 0, C"},      // 81
	0x82: {fn: func(c *CPU) { c.opRES_0_D() }, Name: "RES 0, D"},      // 82
	0x83: {fn: func(c *CPU) { c.opRES_0_E() }, Name: "RES 0, E"},      // 83
	0x84: {fn: func(c *CPU) { c.opRES_0_H() }, Name: "RES 0, H"},      // 84
	0x85: {fn: func(c *CPU) { c.opRES_0_L() }, Name: "RES 0, L"},      // 85
	0x86: {fn: func(c *CPU) { c.opRES_0_aHL() }, Name: "RES 0, [HL]"}, // 86
	0x87: {fn: func(c *CPU) { c.opRES_0_A() }, Name: "RES 0, A"},      // 87
	0x88: {fn: func(c *CPU) { c.opRES_1_B() }, Name: "RES 1, B"},      // 88
	0x89: {fn: func(c *CPU) { c.opRES_1_C() }, Name: "RES 1, C"},      // 89
	0x8A: {fn: func(c *CPU) { c.opRES_1_D() }, Name: "RES 1, D"},      // 8A
	0x8B: {fn: func(c *CPU) { c.opRES_1_E() }, Name: "RES 1, E"},      // 8B
	0x8C: {fn: func(c *CPU) { c.opRES_1_H() }, Name: "RES 1, H"},      // 8C
	0x8D: {fn: func(c *CPU) { c.opRES_1_L() }, Name: "RES 1, L"},      // 8D
	0x8E: {fn: func(c *CPU) { c.opRES_1_aHL() }, Name: "RES 1, [HL]"}, // 8E
	0x8F: {fn: func(c *CPU) { c.opRES_1_A() }, Name: "RES 1, A"},      // 8F
	0x90: {fn: func(c *CPU) { c.opRES_2_B() }, Name: "RES 2, B"},      // 90
	0x91: {fn: func(c *CPU) { c.opRES_2_C() }, Name: "RES 2, C"},      // 91
	0x92: {fn: func(c *CPU) { c.opRES_2_D() }, Name: "RES 2, D"},      // 92
	0x93: {fn: func(c *CPU) { c.opRES_2_E() }, Name: "RES 2, E"},      // 93
	0x94: {fn: func(c *CPU) { c.opRES_2_H() }, Name: "RES 2, H"},      // 94
	0x95: {fn: func(c *CPU) { c.opRES_2_L() }, Name: "RES 2, L"},      // 95
	0x96: {fn: func(c *CPU) { c.opRES_2_aHL() }, Name: "RES 2, [HL]"}, // 96
	0x97: {fn: func(c *CPU) { c.opRES_2_A() }, Name: "RES 2, A"},      // 97
	0x98: {fn: func(c *CPU) { c.opRES_3_B() }, Name: "RES 3, B"},      // 98
	0x99: {fn: func(c *CPU) { c.opRES_3_C() }, Name: "RES 3, C"},      // 99
	0x9A: {fn: func(c *CPU) { c.opRES_3_D() }, Name: "RES 3, D"},      // 9A
	0x9B: {fn: func(c *CPU) { c.opRES_3_E() }, Name: "RES 3, E"},      // 9B
	0x9C: {fn: func(c *CPU) { c.opRES_3_H() }, Name: "RES 3, H"},      // 9C
	0x9D: {fn: func(c *CPU) { c.opRES_3_L() }, Name: "RES 3, L"},      // 9D
	0x9E: {fn: func(c *CPU) { c.opRES_3_aHL() }, Name: "RES 3, [HL]"}, // 9E
	0x9F: {fn: func(c *CPU) { c.opRES_3_A() }, Name: "RES 3, A"},      // 9F
	0xA0: {fn: func(c *CPU) { c.opRES_4_B() }, Name: "RES 4, B"},      // A0
	0xA1: {fn: func(c *CPU) { c.opRES_4_C() }, Name: "RES 4, C"},      // A1
	0xA2: {fn: func(c *CPU) { c.opRES_4_D() }, Name: "RES 4, D"},      // A2
	0xA3: {fn: func(c *CPU) { c.opRES_4_E() }, Name: "RES 4, E"},      // A3
	0xA4: {fn: func(c *CPU) { c.opRES_4_H() }, Name: "RES 4, H"},      // A4
	0xA5: {fn: func(c *CPU) { c.opRES_4_L() }, Name: "RES 4, L"},      // A5
	0xA6: {fn: func(c *CPU) { c.opRES_4_aHL() }, Name: "RES 4, [HL]"}, // A6
	0xA7: {fn: func(c *CPU) { c.opRES_4_A() }, Name: "RES 4, A"},      // A7
	0xA8: {fn: func(c *CPU) { c.opRES_5_B() }, Name: "RES 5, B"},      // A8
	0xA9: {fn: func(c *CPU) { c.opRES_5_C() }, Name: "RES 5, C"},      // A9
	0xAA: {fn: func(c *CPU) { c.opRES_5_D() }, Name: "RES 5, D"},      // AA
	0xAB: {fn: func(c *CPU) { c.opRES_5_E() }, Name: "RES 5, E"},      // AB
	0xAC: {fn: func(c *CPU) { c.opRES_5_H() }, Name: "RES 5, H"},      // AC
	0xAD: {fn: func(c *CPU) { c.opRES_5_L() }, Name: "RES 5, L"},      // AD
	0xAE: {fn: func(c *CPU) { c.opRES_5_aHL() }, Name: "RES 5, [HL]"}, // AE
	0xAF: {fn: func(c *CPU) { c.opRES_5_A() }, Name: "RES 5, A"},      // AF
	0xB0: {fn: func(c *CPU) { c.opRES_6_B() }, Name: "RES 6, B"},      // B0
	0xB1: {fn: func(c *CPU) { c.opRES_6_C() }, Name: "RES 6, C"},      // B1
	0xB2: {fn: func(c *CPU) { c.opRES_6_D() }, Name: "RES 6, D"},      // B2
	0xB3: {fn: func(c *CPU) { c.opRES_6_E() }, Name: "RES 6, E"},      // B3
	0xB4: {fn: func(c *CPU) { c.opRES_6_H() }, Name: "RES 6, H"},      // B4
	0xB5: {fn: func(c *CPU) { c.opRES_6_L() }, Name: "RES 6, L"},      // B5
	0xB6: {fn: func(c *CPU) { c.opRES_6_aHL() }, Name: "RES 6, [HL]"}, // B6
	0xB7: {fn: func(c *CPU) { c.opRES_6_A() }, Name: "RES 6, A"},      // B7
	0xB8: {fn: func(c *CPU) { c.opRES_7_B() }, Name: "RES 7, B"},      // B8
	0xB9: {fn: func(c *CPU) { c.opRES_7_C() }, Name: "RES 7, C"},      // B9
	0xBA: {fn: func(c *CPU) { c.opRES_7_D() }, Name: "RES 7, D"},      // BA
	0xBB: {fn: func(c *CPU) { c.opRES_7_E() }, Name: "RES 7, E"},      // BB
	0xBC: {fn: func(c *CPU) { c.opRES_7_H() }, Name: "RES 7, H"},      // BC
	0xBD: {fn: func(c *CPU) { c.opRES_7_L() }, Name: "RES 7, L"},      // BD
	0xBE: {fn: func(c *CPU) { c.opRES_7_aHL() }, Name: "RES 7, [HL]"}, // BE
	0xBF: {fn: func(c *CPU) { c.opRES_7_A() }, Name: "RES 7, A"},      // BF
	0xC0: {fn: func(c *CPU) { c.opSET_0_B() }, Name: "SET 0, B"},      // C0
	0xC1: {fn: func(c *CPU) { c.opSET_0_C() }, Name: "SET 0, C"},      // C1
	0xC2: {fn: func(c *CPU) { c.opSET_0_D() }, Name: "SET 0, D"},      // C2
	0xC3: {fn: func(c *CPU) { c.opSET_0_E() }, Name: "SET 0, E"},      // C3
	0xC4: {fn: func(c *CPU) { c.opSET_0_H() }, Name: "SET 0, H"},      // C4
	0xC5: {fn: func(c *CPU) { c.opSET_0_L() }, Name: "SET 0, L"},      // C5
	0xC6: {fn: func(c *CPU) { c.opSET_0_aHL() }, Name: "SET 0, [HL]"}, // C6
	0xC7: {fn: func(c *CPU) { c.opSET_0_A() }, Name: "SET 0, A"},      // C7
	0xC8: {fn: func(c *CPU) { c.opSET_1_B() }, Name: "SET 1, B"},      // C8
	0xC9: {fn: func(c *CPU) { c.opSET_1_C() }, Name: "SET 1, C"},      // C9
	0xCA: {fn: func(c *CPU) { c.opSET_1_D() }, Name: "SET 1, D"},      // CA
	0xCB: {fn: func(c *CPU) { c.opSET_1_E() }, Name: "SET 1, E"},      // CB
	0xCC: {fn: func(c *CPU) { c.opSET_1_H() }, Name: "SET 1, H"},      // CC
	0xCD: {fn: func(c *CPU) { c.opSET_1_L() }, Name: "SET 1, L"},      // CD
	0xCE: {fn: func(c *CPU) { c.opSET_1_aHL() }, Name: "SET 1, [HL]"}, // CE
	0xCF: {fn: func(c *CPU) { c.opSET_1_A() }, Name: "SET 1, A"},      // CF
	0xD0: {fn: func(c *CPU) { c.opSET_2_B() }, Name: "SET 2, B"},      // D0
	0xD1: {fn: func(c *CPU) { c.opSET_2_C() }, Name: "SET 2, C"},      // D1
	0xD2: {fn: func(c *CPU) { c.opSET_2_D() }, Name: "SET 2, D"},      // D2
	0xD3: {fn: func(c *CPU) { c.opSET_2_E() }, Name: "SET 2, E"},      // D3
	0xD4: {fn: func(c *CPU) { c.opSET_2_H() }, Name: "SET 2, H"},      // D4
	0xD5: {fn: func(c *CPU) { c.opSET_2_L() }, Name: "SET 2, L"},      // D5
	0xD6: {fn: func(c *CPU) { c.opSET_2_aHL() }, Name: "SET 2, [HL]"}, // D6
	0xD7: {fn: func(c *CPU) { c.opSET_2_A() }, Name: "SET 2, A"},      // D7
	0xD8: {fn: func(c *CPU) { c.opSET_3_B() }, Name: "SET 3, B"},      // D8
	0xD9: {fn: func(c *CPU) { c.opSET_3_C() }, Name: "SET 3, C"},      // D9
	0xDA: {fn: func(c *CPU) { c.opSET_3_D() }, Name: "SET 3, D"},      // DA
	0xDB: {fn: func(c *CPU) { c.opSET_3_E() }, Name: "SET 3, E"},      // DB
	0xDC: {fn: func(c *CPU) { c.opSET_3_H() }, Name: "SET 3, H"},      // DC
	0xDD: {fn: func(c *CPU) { c.opSET_3_L() }, Name: "SET 3, L"},      // DD
	0xDE: {fn: func(c *CPU) { c.opSET_3_aHL() }, Name: "SET 3, [HL]"}, // DE
	0xDF: {fn: func(c *CPU) { c.opSET_3_A() }, Name: "SET 3, A"},      // DF
	0xE0: {fn: func(c *CPU) { c.opSET_4_B() }, Name: "SET 4, B"},      // E0
	0xE1: {fn: func(c *CPU) { c.opSET_4_C() }, Name: "SET 4, C"},      // E1
	0xE2: {fn: func(c *CPU) { c.opSET_4_D() }, Name: "SET 4, D"},      // E2
	0xE3: {fn: func(c *CPU) { c.opSET_4_E() }, Name: "SET 4, E"},      // E3
	0xE4: {fn: func(c *CPU) { c.opSET_4_H() }, Name: "SET 4, H"},      // E4
	0xE5: {fn: func(c *CPU) { c.opSET_4_L() }, Name: "SET 4, L"},      // E5
	0xE6: {fn: func(c *CPU) { c.opSET_4_aHL() }, Name: "SET 4, [HL]"}, // E6
	0xE7: {fn: func(c *CPU) { c.opSET_4_A() }, Name: "SET 4, A"},      // E7
	0xE8: {fn: func(c *CPU) { c.opSET_5_B() }, Name: "SET 5, B"},      // E8
	0xE9: {fn: func(c *CPU) { c.opSET_5_C() }, Name: "SET 5, C"},      // E9
	0xEA: {fn: func(c *CPU) { c.opSET_5_D() }, Name: "SET 5, D"},      // EA
	0xEB: {fn: func(c *CPU) { c.opSET_5_E() }, Name: "SET 5, E"},      // EB
	0xEC: {fn: func(c *CPU) { c.opSET_5_H() }, Name: "SET 5, H"},      // EC
	0xED: {fn: func(c *CPU) { c.opSET_5_L() }, Name: "SET 5, L"},      // ED
	0xEE: {fn: func(c *CPU) { c.opSET_5_aHL() }, Name: "SET 5, [HL]"}, // EE
	0xEF: {fn: func(c *CPU) { c.opSET_5_A() }, Name: "SET 5, A"},      // EF
	0xF0: {fn: func(c *CPU) { c.opSET_6_B() }, Name: "SET 6, B"},      // F0
	0xF1: {fn: func(c *CPU) { c.opSET_6_C() }, Name: "SET 6, C"},      // F1
	0xF2: {fn: func(c *CPU) { c.opSET_6_D() }, Name: "SET 6, D"},      // F2
	0xF3: {fn: func(c *CPU) { c.opSET_6_E() }, Name: "SET 6, E"},      // F3
	0xF4: {fn: func(c *CPU) { c.opSET_6_H() }, Name: "SET 6, H"},      // F4
	0xF5: {fn: func(c *CPU) { c.opSET_6_L() }, Name: "SET 6, L"},      // F5
	0xF6: {fn: func(c *CPU) { c.opSET_6_aHL() }, Name: "SET 6, [HL]"}, // F6
	0xF7: {fn: func(c *CPU) { c.opSET_6_A() }, Name: "SET 6, A"},      // F7
	0xF8: {fn: func(c *CPU) { c.opSET_7_B() }, Name: "SET 7, B"},      // F8
	0xF9: {fn: func(c *CPU) { c.opSET_7_C() }, Name: "SET 7, C"},      // F9
	0xFA: {fn: func(c *CPU) { c.opSET_7_D() }, Name: "SET 7, D"},      // FA
	0xFB: {fn: func(c *CPU) { c.opSET_7_E() }, Name: "SET 7, E"},      // FB
	0xFC: {fn: func(c *CPU) { c.opSET_7_H() }, Name: "SET 7, H"},      // FC
	0xFD: {fn: func(c *CPU) { c.opSET_7_L() }, Name: "SET 7, L"},      // FD
	0xFE: {fn: func(c *CPU) { c.opSET_7_aHL() }, Name: "SET 7, [HL]"}, // FE
	0xFF: {fn: func(c *CPU) { c.opSET_7_A() }, Name: "SET 7, A"},      // FF
}
