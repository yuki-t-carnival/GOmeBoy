package ppu

import "gomeboy/internal/util"

const (
	// ID for BG/Window common functions
	BG     = 0
	Window = 1

	// Kinds of palette
	BGP  = 0
	OBP0 = 1
	OBP1 = 2
)

type PPU struct {
	// LCD
	vp [160 * 144]Pixel // Color Num
	FB []byte           // After Resolve palette

	// PPU Memory
	vram [2][0x2000]byte
	oam  [160]byte

	// I/O Registers
	bank     byte
	dma      byte
	lcdcBits [8]uint8
	stat     byte
	ly       byte
	lyc      byte
	obp0     byte
	obp1     byte
	bgp      byte
	wy       byte
	wx       byte
	scy      byte
	scx      byte

	// Prev PPU Registers
	prevLCDC5 uint8
	prevLY    byte

	// PPU Internal Counters
	wly    int
	cycles int

	// Interrupt
	HasSTATIRQ    bool
	HasVBlankIRQ  bool
	isLockLYCInt  bool
	isLockModeInt bool

	// Object
	objList  []int
	xFlipLUT [256]byte

	// Palette
	plBGP  [4]uint8
	plOBP0 [4]uint8
	plOBP1 [4]uint8

	hasTransferRq bool
}

type Pixel struct {
	num byte  // 0 ~ 3
	pl  uint8 // BGP, PlaletteOBP0, OBP1
}

func NewPPU() *PPU {
	p := &PPU{
		FB:   make([]byte, 160*144),
		stat: 0x85,
		ly:   0x00,
		lyc:  0x00,
		wy:   0x00,
		wx:   0x00,
		scy:  0x00,
		scx:  0x00,
	}
	p.SetLCDC(0x91)
	p.SetBGP(0xFC)
	p.SetOBP0(0xFF)
	p.SetOBP1(0xFF)
	p.createXFlipLUT()
	return p
}

// For Object
func (p *PPU) createXFlipLUT() {
	result := byte(0)
	for i := 0; i < 256; i++ {
		result = 0
		for b := 0; b < 8; b++ {
			if i&(1<<b) != 0 {
				result |= (1 << (7 - b))
			}
		}
		p.xFlipLUT[i] = result
	}
}

func (p *PPU) Step(cpuCycles int) {
	// In case of PPU disabled
	if p.lcdcBits[7] == 0 {
		p.prevLY = p.ly
		p.ly = 0
		p.wly = 0
		return
	}

	// Check LYC==LY Interrupts
	p.stat = (p.stat &^ byte(1<<2)) // LYC==LY bit clear
	if p.lyc == p.ly {
		p.stat |= 1 << 2
		isLYCIntSel := (p.stat & byte(1<<6)) != 0
		if isLYCIntSel && !p.isLockLYCInt { // LYC==LY && LYCint enabled
			p.HasSTATIRQ = true
			p.isLockLYCInt = true
		}
	}

	switch {
	// During the period LY=0~143, repeat Mode2>3>0> every LY
	case p.ly <= 143:
		if p.ly != p.prevLY {
			p.hasTransferRq = true
		}
		// When starting with cycle=0, some settings may not be applied
		if p.hasTransferRq && p.cycles >= 280 {
			p.setPPUMode(2)
			p.oamSearch()
			p.setPPUMode(3)
			p.pixelTransfer()
			p.setPPUMode(0)
			p.hasTransferRq = false
		}
	// VBlank during LY=144~153 (VBlank IRQ occurs only once when LY=144)
	case p.ly == 144:
		if p.ly != p.prevLY {
			p.setPPUMode(1)
			p.HasVBlankIRQ = true
			p.wly = 0
		}
	}
	p.prevLY = p.ly

	p.checkSTATInt()

	// 1 Line == 456 CPU cycles
	p.cycles += cpuCycles
	if p.cycles >= 456 {
		p.cycles -= 456
		p.ly++
		p.isLockLYCInt = false
		if p.ly == 154 {
			p.ly = 0
		}
	}
}

// Update the frame buffer by one line
func (p *PPU) pixelTransfer() {
	// init vp
	base := int(p.ly) * 160
	for i := 0; i < 160; i++ {
		vp := &p.vp[base+i]
		vp.pl = 0
		vp.num = 0
	}

	// BG/Window is On
	if p.lcdcBits[0] == 1 {

		// BG
		p.BGTransfer()

		// Window is On
		if p.lcdcBits[5] == 1 {
			if p.prevLCDC5 == 0 {
				p.wly = 0
			}
			if p.ly >= p.wy {
				isDrawn := p.WindowTransfer()
				if isDrawn {
					p.wly++
				}
			}
		}
		p.prevLCDC5 = p.lcdcBits[5]
	}

	// Obj is On
	if p.lcdcBits[1] == 1 {
		p.objectsTransfer()
	}

	// Convert vp to FB
	p.resolvePalette()
}

func (p *PPU) setPPUMode(nextMode int) {
	p.stat = (p.stat & 0x7C) | (byte(nextMode) & 0x03)
	p.isLockModeInt = false
}

func (p *PPU) checkSTATInt() {
	mode := p.stat & 0x03
	isIntEnabled := p.stat&(1<<(mode+3)) != 0
	if isIntEnabled && !p.isLockModeInt {
		p.HasSTATIRQ = true
		p.isLockModeInt = true
	}
}

// Lists 0~10 objects to be drawn on the current line
func (p *PPU) oamSearch() {
	p.objList = p.objList[:0] // =[]int{}

	bigMode := p.lcdcBits[2] == 1
	ly := int(p.ly)

	// List the objects in ascending OAM index order
	for i := 0; i < 40; i++ {

		// Y origin of Object(i) on the viewport
		y0 := int(p.oam[i<<2+0]) - 16

		var objHeight int
		if bigMode {
			objHeight = 16
		} else {
			objHeight = 8
		}

		// If part of the Y Position of Object(i) == LY, add it to the list
		if ly >= y0 && ly < (y0+objHeight) {
			p.objList = append(p.objList, i)
			// Ends when 10 objcets are found
			if len(p.objList) == 10 {
				return
			}
		}
	}
}

// Draw the current line objects listed by oamSearch() to vp
func (p *PPU) objectsTransfer() {
	ly := int(p.ly)

	// Sort the list X Position DESC or OAM index DESC
	sortedList := []int{}
	for _, oami := range p.objList {
		insPos := len(sortedList)
		for j, oamj := range sortedList {
			if p.oam[oami<<2+1] >= p.oam[oamj<<2+1] {
				insPos = j
				break
			}
		}
		sortedList = util.InsertSlice(sortedList, insPos, oami)
	}

	// Draw objects to vp in the order sorted above
	for _, oamIdx := range sortedList {

		// Get Object attribytes in the OAM
		base := uint16(oamIdx << 2)
		y0 := int(p.oam[base]) - 16  // Byte 0 - Y Position
		x0 := int(p.oam[base+1]) - 8 // Byte 1 - X Position
		idx := p.oam[base+2]         // Byte 2 - Tile Index
		attr := p.oam[base+3]        // Byte 3 - Attributes/Flags
		tilePy := ly - y0

		data := [2]byte{}
		data = p.getObjectTile(idx, attr, tilePy)

		var pl uint8
		if (attr & (1 << 4)) == 0 {
			pl = OBP0
		} else {
			pl = OBP1
		}

		// Draw tileData(one row) on vp
		tgtY := ly * 160
		for b := 0; b < 8; b++ {
			tgtX := x0 + b
			if tgtX < 0 || tgtX >= 160 {
				continue
			}
			tgt := tgtY + tgtX
			bgwNum := p.vp[tgt].num

			// In the case below, the Object pixel is hidden under the BG/W
			pri := (attr & (1 << 7)) >> 7
			if pri == 1 && bgwNum != 0 {
				continue
			}
			// 2-bit monochrome (DMG)
			lo := (data[0] >> (7 - b)) & 1
			hi := (data[1] >> (7 - b)) & 1
			num := (hi << 1) | lo

			// When a tile is used in an object, ID 0 means transparent
			if num == 0 {
				continue
			}

			p.vp[tgt].pl = pl
			p.vp[tgt].num = num
		}
	}
}

// Get one row of tileData as an object
func (p *PPU) getObjectTile(idx, attr byte, tilePy int) [2]byte {
	isYFlip := attr&(1<<6) != 0
	isXFlip := attr&(1<<5) != 0
	if p.lcdcBits[2] == 1 { // OBJ size
		idx &= 0xFE // When 8x16, idx is masked to even numbers only
		if isYFlip {
			tilePy = 15 - tilePy
		}
	} else {
		if isYFlip {
			tilePy = 7 - tilePy
		}
	}
	base := uint16(idx) << 4
	data := [2]byte{}
	for i := 0; i < 2; i++ {
		addr := base + uint16(tilePy<<1+i)
		data[i] = p.ReadVRAM(addr)
		if isXFlip {
			data[i] = p.xFlipLUT[data[i]]
		}
	}
	return data
}

// Draw the current line background to vp
func (p *PPU) BGTransfer() {
	var data [2]byte
	vpIdxY := int(p.ly) * 160
	bgY := int(p.ly + p.scy)
	bgRow := bgY >> 3    // =/8
	tilePy := bgY & 0x07 // =%8
	for x := 0; x < 160; x++ {
		bgX := byte(x) + p.scx // wrapされる
		tilePx := bgX & 0x07   // =%8

		// Get tileData only when needed
		if x == 0 || tilePx == 0 {
			bgCol := int(bgX) >> 3
			bgIdx := bgRow<<5 + bgCol // <<5 == *32
			data = p.getBGWTile(BG, bgIdx, tilePy)
		}
		lo := (data[0] >> (7 - tilePx)) & 1
		hi := (data[1] >> (7 - tilePx)) & 1
		num := (hi << 1) | lo
		p.vp[vpIdxY+x].num = num
		p.vp[vpIdxY+x].pl = BGP
	}
}

// Draw the current line Window to vp
func (p *PPU) WindowTransfer() bool {
	var isDrawn bool
	var data [2]byte
	vpIdxY := int(p.ly) * 160
	wRow := p.wly >> 3     // =/8
	tilePy := p.wly & 0x07 // =%8
	for x := 0; x < 160; x++ {
		bgX := x - (int(p.wx) - 7)
		tilePx := bgX & 0x07 // =%8

		// Get tileData only when needed
		if x == 0 || tilePx == 0 {
			wCol := bgX >> 3
			wIdx := wRow<<5 + wCol // <<5 == *32
			data = p.getBGWTile(Window, wIdx, tilePy)
		}
		if bgX < 0 || bgX >= 160 {
			continue
		}
		lo := (data[0] >> (7 - tilePx)) & 1
		hi := (data[1] >> (7 - tilePx)) & 1
		num := (hi << 1) | lo
		p.vp[vpIdxY+x].num = num
		p.vp[vpIdxY+x].pl = BGP
		isDrawn = true
	}
	return isDrawn
}

// Get one row of tileData as background or window
func (p *PPU) getBGWTile(bgwID, idx, py int) [2]byte {
	var mapStart uint16
	var addrType int
	switch bgwID {
	case BG:
		addrType = int(p.lcdcBits[3])
	case Window:
		addrType = int(p.lcdcBits[6])
	}
	if addrType == 0 {
		mapStart = 0x1800
	} else {
		mapStart = 0x1C00
	}
	mapAddr := mapStart + uint16(idx)
	tileIdx := p.ReadVRAM(mapAddr)

	tileStart := uint16(0)
	if p.lcdcBits[4] == 0 {
		tileStart = uint16(0x1000 + int(int8(tileIdx))<<4)
	} else {
		tileStart = uint16(tileIdx) << 4
	}
	tileAddr := tileStart + uint16(py*2)
	return [2]byte{
		p.ReadVRAM(tileAddr),
		p.ReadVRAM(tileAddr + 1),
	}
}

// After color conversion, transfer vp to the frame buffer
func (p *PPU) resolvePalette() {
	tgtBase := int(p.ly) * 160
	for x := 0; x < 160; x++ {
		tgt := tgtBase + x
		data := p.vp[tgt].num
		pl := p.vp[tgt].pl
		switch pl {
		case BGP:
			p.FB[tgt] = p.plBGP[data]
		case OBP0:
			p.FB[tgt] = p.plOBP0[data]
		case OBP1:
			p.FB[tgt] = p.plOBP1[data]
		default:
			p.FB[tgt] = 0
		}
	}
}

func (p *PPU) ReadVRAM(addr uint16) byte {
	offset := addr & 0x1FFF // To prevent out of range errors
	return p.vram[p.bank][offset]
}

func (p *PPU) WriteVRAM(addr uint16, val byte) {
	offset := addr & 0x1FFF // To prevent out of range errors
	p.vram[p.bank][offset] = val
}

func (p *PPU) GetVBK() byte {
	return 0xFE | (p.bank & 0x01)
}

func (p *PPU) SetVBK(val byte) {
	p.bank = val & 0x01
}

func (p *PPU) ReadOAM(addr uint16) byte {
	return p.oam[addr]
}
func (p *PPU) WriteOAM(addr uint16, val byte) {
	p.oam[addr] = val
}

func (p *PPU) GetDMA() byte {
	return p.dma
}

func (p *PPU) SetDMA(val byte) {
	p.dma = val
}

func (p *PPU) GetLCDC() byte {
	lcdc := byte(0)
	for i, v := range p.lcdcBits {
		lcdc |= v << i
	}
	return lcdc
}

func (p *PPU) SetLCDC(val byte) {
	p.lcdcBits = [8]uint8{}
	for i := 0; i < 8; i++ {
		if val&(1<<i) != 0 {
			p.lcdcBits[i] = 1
		}
	}
}

func (p *PPU) GetSTAT() byte {
	return p.stat
}

func (p *PPU) SetSTAT(val byte) {
	p.stat = (val & 0x7C) | (p.stat & 0x03)

	// If any Mode int select is changed, check for interrupts
	p.checkSTATInt()
}

func (p *PPU) GetLY() byte {
	return p.ly
}

func (p *PPU) GetLYC() byte {
	return p.lyc
}

func (p *PPU) SetLYC(val byte) {
	p.lyc = val
}

func (p *PPU) GetOBP0() byte {
	return p.obp0
}

func (p *PPU) SetOBP0(val byte) {
	for i := 0; i < 4; i++ {
		p.plOBP0[i] = (val & (0x03 << (2 * uint8(i)))) >> (2 * uint8(i))
	}
	p.obp0 = val
}

func (p *PPU) GetOBP1() byte {
	return p.obp1
}

func (p *PPU) SetOBP1(val byte) {
	for i := 0; i < 4; i++ {
		p.plOBP1[i] = (val & (0x03 << (2 * uint8(i)))) >> (2 * uint8(i))
	}
	p.obp1 = val
}

func (p *PPU) GetBGP() byte {
	return p.bgp
}

func (p *PPU) SetBGP(val byte) {
	for i := 0; i < 4; i++ {
		p.plBGP[i] = (val & (0x03 << (2 * uint8(i)))) >> (2 * uint8(i))
	}
	p.bgp = val
}

func (p *PPU) GetWY() byte {
	return p.wy
}

func (p *PPU) SetWY(val byte) {
	p.wy = val
}

func (p *PPU) GetWX() byte {
	return p.wx
}

func (p *PPU) SetWX(val byte) {
	p.wx = val
}

func (p *PPU) GetSCY() byte {
	return p.scy
}

func (p *PPU) SetSCY(val byte) {
	p.scy = val
}

func (p *PPU) GetSCX() byte {
	return p.scx
}

func (p *PPU) SetSCX(val byte) {
	p.scx = val
}
