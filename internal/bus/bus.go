package bus

import (
	"gomeboy/internal/joypad"
	"gomeboy/internal/memory"
	"gomeboy/internal/ppu"
	"gomeboy/internal/timer"
)

type Bus struct {
	PPU    *ppu.PPU
	Timer  *timer.Timer
	Joypad *joypad.Joypad
	Memory *memory.Memory

	// DMA Transfer
	IsDMATransferInProgress bool
	DMATransferIndex        int // 0 ~ 159
}

// Memory Map Adress
const (
	// I/O Registers (FF00 ~ FF7F)
	P1_JOYP      uint16 = 0xFF00
	SB           uint16 = 0xFF01
	SC           uint16 = 0xFF02
	DIV          uint16 = 0xFF04
	TIMA         uint16 = 0xFF05
	TMA          uint16 = 0xFF06
	TAC          uint16 = 0xFF07
	IF           uint16 = 0xFF0F
	NR10         uint16 = 0xFF10
	NR11         uint16 = 0xFF11
	NR12         uint16 = 0xFF12
	NR13         uint16 = 0xFF13
	NR14         uint16 = 0xFF14
	NR21         uint16 = 0xFF16
	NR22         uint16 = 0xFF17
	NR23         uint16 = 0xFF18
	NR24         uint16 = 0xFF19
	NR30         uint16 = 0xFF1A
	NR31         uint16 = 0xFF1B
	NR32         uint16 = 0xFF1C
	NR33         uint16 = 0xFF1D
	NR34         uint16 = 0xFF1E
	NR41         uint16 = 0xFF20
	NR42         uint16 = 0xFF21
	NR43         uint16 = 0xFF22
	NR44         uint16 = 0xFF23
	NR50         uint16 = 0xFF24
	NR51         uint16 = 0xFF25
	NR52         uint16 = 0xFF26
	WaveRAMStart uint16 = 0xFF30
	LCDC         uint16 = 0xFF40
	STAT         uint16 = 0xFF41
	SCY          uint16 = 0xFF42
	SCX          uint16 = 0xFF43
	LY           uint16 = 0xFF44
	LYC          uint16 = 0xFF45
	DMA          uint16 = 0xFF46
	BGP          uint16 = 0xFF47
	OBP0         uint16 = 0xFF48
	OBP1         uint16 = 0xFF49
	WY           uint16 = 0xFF4A
	WX           uint16 = 0xFF4B
	KEY0_SYS     uint16 = 0xFF4C
	KEY1_SPD     uint16 = 0xFF4D
	VBK          uint16 = 0xFF4F
	BANK         uint16 = 0xFF50
	HDMA1        uint16 = 0xFF51
	HDMA2        uint16 = 0xFF52
	HDMA3        uint16 = 0xFF53
	HDMA4        uint16 = 0xFF54
	HDMA5        uint16 = 0xFF55
	RP           uint16 = 0xFF56
	BCPS_BGPI    uint16 = 0xFF68
	BCPD_BGPD    uint16 = 0xFF69
	OCPS_OBPI    uint16 = 0xFF6A
	OCPD_OBPD    uint16 = 0xFF6B
	OPRI         uint16 = 0xFF6C
	SVBK_WBK     uint16 = 0xFF70
	PCM12        uint16 = 0xFF76
	PCM34        uint16 = 0xFF77

	// Interrupt Enable Register
	IE uint16 = 0xFFFF
)

func NewBus(m *memory.Memory) *Bus {
	// Serial
	m.Write(SB, 0x00)
	m.Write(SC, 0x7E)

	// Interrupt
	m.Write(IF, 0x01)
	m.Write(IE, 0x00)

	// Sound
	m.Write(NR10, 0x80)
	m.Write(NR11, 0xBF)
	m.Write(NR12, 0xF3)
	m.Write(NR14, 0xBF)
	m.Write(NR21, 0x3F)
	m.Write(NR22, 0x00)
	m.Write(NR24, 0xBF)
	m.Write(NR30, 0x7F)
	m.Write(NR31, 0xFF)
	m.Write(NR32, 0x9F)
	m.Write(NR34, 0xBF)
	m.Write(NR41, 0xFF)
	m.Write(NR42, 0x00)
	m.Write(NR43, 0x00)
	m.Write(NR44, 0xBF)
	m.Write(NR50, 0x77)
	m.Write(NR51, 0xF3)
	m.Write(NR52, 0xF1)

	bus := &Bus{
		PPU:    ppu.NewPPU(),
		Timer:  timer.NewTimer(),
		Joypad: joypad.NewJoypad(),
		Memory: m,
	}
	return bus
}

// Bus accesses the I/O, VRAM, OAM,
// and request other accesses to Memory
func (b *Bus) Read(addr uint16) byte {
	switch {
	// PPU
	case addr >= 0x8000 && addr < 0xA000:
		return b.PPU.ReadVRAM(addr - 0x8000)
	case addr >= 0xFE00 && addr < 0xFEA0:
		return b.PPU.ReadOAM(addr - 0xFE00)
	case addr == VBK:
		return b.PPU.GetVBK()
	case addr == DMA:
		return b.PPU.GetDMA()
	case addr == LCDC:
		return b.PPU.GetLCDC()
	case addr == STAT:
		return b.PPU.GetSTAT()
	case addr == LY:
		return b.PPU.GetLY()
	case addr == LYC:
		return b.PPU.GetLYC()
	case addr == OBP0:
		return b.PPU.GetOBP0()
	case addr == OBP1:
		return b.PPU.GetOBP1()
	case addr == BGP:
		return b.PPU.GetBGP()
	case addr == WY:
		return b.PPU.GetWY()
	case addr == WX:
		return b.PPU.GetWX()
	case addr == SCY:
		return b.PPU.GetSCY()
	case addr == SCX:
		return b.PPU.GetSCX()

	// Timer
	case addr == DIV:
		return b.Timer.GetDIV()
	case addr == TIMA:
		return b.Timer.GetTIMA()
	case addr == TMA:
		return b.Timer.GetTMA()
	case addr == TAC:
		return b.Timer.GetTAC()

	// Joypad
	case addr == P1_JOYP:
		return b.Joypad.GetP1JOYP()

	// Interrupt
	case addr == IF:
		return b.Memory.Read(addr) | 0xE0
	case addr == IE:
		return b.Memory.Read(addr) | 0xE0

	// Memory
	case addr == SVBK_WBK:
		return b.Memory.ReadWRAMBank()
	default:
		return b.Memory.Read(addr)
	}
}

// Bus accesses the I/O, VRAM, OAM,
// and request other accesses to Memory
func (b *Bus) Write(addr uint16, val byte) {
	switch {
	// PPU
	case addr >= 0x8000 && addr < 0xA000:
		b.PPU.WriteVRAM(addr-0x8000, val)
	case addr >= 0xFE00 && addr < 0xFEA0:
		b.PPU.WriteOAM(addr-0xFE00, val)
	case addr == VBK:
		b.PPU.SetVBK(val)
	case addr == DMA:
		b.PPU.SetDMA(val) // Acts as a trigger to start a DMA transfer
		if val <= 0xDF {
			b.IsDMATransferInProgress = true
			b.DMATransferIndex = 0
		}
	case addr == LCDC:
		b.PPU.SetLCDC(val) // TODO: LCD&PPU can be disabled only during VBlank period
	case addr == STAT:
		b.PPU.SetSTAT(val)
	case addr == SCY:
		b.PPU.SetSCY(val)
	case addr == SCX:
		b.PPU.SetSCX(val)
	case addr == LY:
		// Writing is prohibited
	case addr == LYC:
		b.PPU.SetLYC(val)
	case addr == BGP:
		b.PPU.SetBGP(val)
	case addr == OBP0:
		b.PPU.SetOBP0(val)
	case addr == OBP1:
		b.PPU.SetOBP1(val)
	case addr == WY:
		b.PPU.SetWY(val)
	case addr == WX:
		b.PPU.SetWX(val)

	// Timer
	case addr == DIV:
		b.Timer.ResetDiv()
	case addr == TIMA:
		b.Timer.SetTIMA(val)
	case addr == TMA:
		b.Timer.SetTMA(val)
	case addr == TAC:
		b.Timer.SetTAC(val)

	// Joypad
	case addr == P1_JOYP:
		b.Joypad.SetP1JOYP(val)

	// Interrupt
	case addr == IF:
		b.Memory.Write(IF, val&0x1F)
	case addr == IE:
		b.Memory.Write(IE, val&0x1F)

	// Memory
	case addr == SVBK_WBK:
		b.Memory.WriteWRAMBank(val)
	default:
		b.Memory.Write(addr, val)
	}
}

// During DMA Transfer, CPU is stopped and 4 cycles elapse for each byte transferred
func (b *Bus) DMATransfer() {
	if !(b.PPU.GetDMA() <= 0xDF) || !b.IsDMATransferInProgress {
		panic("DMA transfer error")
	}
	srcBase := uint16(b.PPU.GetDMA()) << 8
	i := uint16(b.DMATransferIndex)
	v := b.Read(srcBase + i)
	b.PPU.WriteOAM(i, v)
	b.DMATransferIndex++
	if b.DMATransferIndex == 160 {
		b.IsDMATransferInProgress = false
		b.DMATransferIndex = 0
		return
	}
}
