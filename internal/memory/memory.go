package memory

type Memory struct {
	rom  []byte          // =.gb data
	ERAM [0x8000]byte    // =External RAM, SRAM
	wram [8][0x1000]byte // DMG=1bank, CGB=8bank
	hram [0x7F]byte
	io   [0x80]byte
	ie   byte

	// Cartridge Header
	cartTypeID   byte   // $0147
	romSizeID    byte   // $0148
	ramSizeID    byte   // $0149
	CartTypeName string // Description of the above ID values
	ROMSizeName  string
	RAMSizeName  string

	// MBC Registers
	bankingMode byte
	bankHigh    byte
	romBankLow  byte
	isRAMEnable bool

	wramBank     byte // only CGB mode
	romBankCount byte

	// For debug
	//ROMBankHistory []byte
	//RAMBankHistory []byte
}

func NewMemory(rom, sav []byte) *Memory {
	mem := &Memory{
		rom:        rom,
		romBankLow: 1,
	}
	copy(mem.ERAM[:], sav)
	mem.cartTypeID = mem.rom[0x0147]
	mem.CartTypeName = cartTypeNames[mem.cartTypeID]
	mem.romSizeID = mem.rom[0x0148]
	mem.ROMSizeName = romSizeNames[mem.romSizeID]
	mem.ramSizeID = mem.rom[0x0149]
	mem.RAMSizeName = ramSizeNames[mem.ramSizeID]
	mem.romBankCount = byte(cap(mem.rom) / 0x4000)
	return mem
}

var cartTypeNames = [256]string{
	0x00: "ROM ONLY",
	0x01: "MBC1",
	0x02: "MBC1+RAM",
	0x03: "MBC1+RAM+BT",
	0x05: "MBC2",
	0x06: "MBC2+BT",
	0x08: "ROM+RAM 11",
	0x09: "ROM+RAM+BT 11",
	0x0B: "MMM01",
	0x0C: "MMM01+RAM",
	0x0D: "MMM01+RAM+BT",
	0x0F: "MBC3+T+BT",
	0x10: "MBC3+T+RAM+BT 12",
	0x11: "MBC3",
	0x12: "MBC3+RAM 12",
	0x13: "MBC3+RAM+BT 12",
	0x19: "MBC5",
	0x1A: "MBC5+RAM",
	0x1B: "MBC5+RAM+BT",
	0x1C: "MBC5+RBL",
	0x1D: "MBC5+RBL+RAM",
	0x1E: "MBC5+RBL+RAM+BT",
	0x20: "MBC6",
	0x22: "MBC7+SEN+RBL+RAM+BT",
	0xFC: "POCKET CAMERA",
	0xFD: "BANDAI TAMA5",
	0xFE: "HuC3",
	0xFF: "HuC1+RAM+BT",
}

var romSizeNames = [0x55]string{
	0x00: "32KiB(2banks,noBanking)",
	0x01: "64KiB(4banks)",
	0x02: "128KiB(8banks)",
	0x03: "256KiB(16banks)",
	0x04: "512KiB(32banks)",
	0x05: "1MiB(64banks)",
	0x06: "2MiB(128banks)",
	0x07: "4MiB(256banks)",
	0x08: "8MiB(512banks)",
	0x52: "1.1MiB(72banks)",
	0x53: "1.2MiB(80banks)",
	0x54: "1.5MiB(96banks)",
}

var ramSizeNames = [6]string{
	0x00: "0(No RAM)",
	0x01: "-(Unused)",
	0x02: "8KiB(1bank)",
	0x03: "32KiB(4banks/8KiB)",
	0x04: "128KiB(16banks/8KiB)",
	0x05: "64KiB(8banks/8KiB)",
}

// Called from Bus.Read()
func (m *Memory) Read(addr uint16) byte {
	switch {
	case addr < 0x8000:
		return m.readROM(addr)

	case addr >= 0x8000 && addr < 0xA000:
		return 0xFF // Access VRAM via Bus.Read()

	case addr >= 0xA000 && addr < 0xC000:
		return m.readERAM(addr)

	case addr >= 0xC000 && addr < 0xD000:
		return m.wram[0][addr-0xC000]
	case addr >= 0xD000 && addr < 0xE000:
		return m.wram[max(m.wramBank, 1)][addr-0xD000]

	case addr >= 0xE000 && addr < 0xFE00:
		mirror := addr - 0x2000
		return m.Read(mirror)

	case addr >= 0xFE00 && addr < 0xFEA0:
		return 0xFF // Access OAM via Bus.Read()

	case addr >= 0xFEA0 && addr < 0xFF00:
		return 0xFF

	case addr >= 0xFF00 && addr < 0xFF80:
		return m.io[addr-0xFF00]

	case addr >= 0xFF80 && addr < 0xFFFF:
		return m.hram[addr-0xFF80]

	case addr == 0xFFFF:
		return m.ie

	default:
		return 0xFF
	}
}

// Called from Bus.Write()
func (m *Memory) Write(addr uint16, val byte) {
	switch {
	case addr < 0x8000:
		m.writeROMArea(addr, val) // Consider banking

	case addr >= 0x8000 && addr < 0xA000:
		return // Access VRAM via Bus.Write()

	case addr >= 0xA000 && addr < 0xC000:
		m.writeERAMArea(addr, val)

	case addr >= 0xC000 && addr < 0xD000:
		m.wram[0][addr-0xC000] = val
	case addr >= 0xD000 && addr < 0xE000:
		m.wram[max(m.wramBank, 1)][addr-0xD000] = val

	case addr >= 0xE000 && addr < 0xFE00:
		mirror := addr - 0x2000
		m.Write(mirror, val)

	case addr >= 0xFE00 && addr < 0xFEA0:
		return // Access OAM via Bus.Write()

	case addr >= 0xFEA0 && addr < 0xFF00:
		return

	case addr >= 0xFF00 && addr < 0xFF80:
		m.io[addr-0xFF00] = val

	case addr >= 0xFF80 && addr < 0xFFFF:
		m.hram[addr-0xFF80] = val

	case addr == 0xFFFF:
		m.ie = val
	}
}

func (m *Memory) ReadWRAMBank() byte {
	return m.wramBank & 0x07
}

func (m *Memory) WriteWRAMBank(val byte) {
	m.wramBank = val & 0x07
}

// Read from ROM in the current bank
// (Only compatible with MBC1)
func (m *Memory) readROM(addr uint16) byte {
	switch {
	case addr < 0x4000: // ROM Bank $20/$40/$60
		bank := byte(0)
		if m.bankingMode == 1 {
			bank = m.bankHigh << 5
		}
		if bank >= m.romBankCount {
			bank %= m.romBankCount
		}
		return m.rom[0x4000*uint32(bank)+uint32(addr)]

	case addr >= 0x4000 && addr < 0x8000: // ROM Bank 01-7F
		bank := (m.bankHigh << 5) | m.romBankLow
		if bank >= m.romBankCount {
			bank %= m.romBankCount
		}
		return m.rom[0x4000*uint32(bank)+uint32(addr-0x4000)]
	default:
		return 0xFF
	}
}

// Read from ERAM in the current bank
// (Only compatible with MBC1)
func (m *Memory) readERAM(addr uint16) byte {
	switch {
	case addr >= 0xA000 && addr < 0xC000:
		if !m.isRAMEnable {
			return 0xFF
		}
		bank := byte(0)
		if m.bankingMode == 1 {
			bank = m.bankHigh
		}
		return m.ERAM[uint16(bank)*0x2000+addr-0xA000]
	default:
		return 0xFF
	}
}

// Write to ROM area
// (it is not a write to the ROM, but a write to the MBC register)
// (Only compatible with MBC1)
func (m *Memory) writeROMArea(addr uint16, val byte) {
	switch {
	case addr < 0x2000: // RAM Enable
		m.isRAMEnable = val&0x0F == 0x0A

	// ROM Bank Number
	case addr >= 0x2000 && addr < 0x4000:
		m.romBankLow = val & 0x1F
		if m.romBankLow == 0 {
			m.romBankLow = 1
		}

	// RAM Bank Number or Upper Bits of ROM Bank Number
	case addr >= 0x4000 && addr < 0x6000:
		m.bankHigh = val & 0x03

	// Banking Mode Select
	case addr >= 0x6000 && addr < 0x8000:
		m.bankingMode = val & 0x01
	}
}

// Write to ERAM area
// (Only compatible with MBC1)
func (m *Memory) writeERAMArea(addr uint16, val byte) {
	switch {
	case addr >= 0xA000 && addr < 0xC000:
		if !m.isRAMEnable {
			return
		}
		bank := byte(0)
		if m.bankingMode == 1 {
			bank = m.bankHigh
		}
		m.ERAM[uint16(bank)*0x2000+addr-0xA000] = val
	}
}
