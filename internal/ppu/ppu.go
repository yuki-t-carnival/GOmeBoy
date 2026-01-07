// ================== PPU Memo =========================================
// ------------------ LCD Position and Scrolling -----------------------
// Viewport==Gameboy画面とする｡
// BGMap上の､Viewportの左上座標が(SCX, SCY)｡ wrapする｡
// Viewport上の､Windowの左上座標が(WX-7, WY)｡ wrapしない｡
//
// ------------------ Object -------------------------------------------
// PPUのOAMに40個分のobject情報（Y, X, TileIndex, Attributeの4Bytes）がある｡計160Bytes
// ObjectのTileDataのアドレス始点は､常に0x8000
// ObjcetはTileSizeが8x8､8x16と選べる(LCDCbit2)
//
// ------------------ Background / Window ------------------------------
// BG/WindowのTileMapは､32*32=256タイル分あり､1タイルにつき1ByteのTileDataIndexが入る｡計256Bytes
// BGのTileMapのアドレス始点は､LCDCbit3で設定
// WindowのTileMapのアドレス始点は､LCDCbit6で設定
// BG/WindowのTileDataのアドレス始点は､両方LCDCbit4で設定
// BG/WindowのTileSizeは常に8x8
//
// ------------------ Tile ---------------------------------------------
// Tileは1Byteずつ､Row0(Plane0), Row0(plane1), Row1(plane0)...とデータが並んでいる｡
// なので､8x8なら16Bytes､8x16(Objのみ)なら32Bytes｡
// plane0は重みが1､plane1は重みが2｡なので､0~3の計4階調｡
//
// ------------------ Cycles -------------------------------------------
// 4194304cycles / 59.73Frames
//	70221cycles / Frame(154Lines)
//	  456cycles / Line

package ppu

import "gomeboy/internal/util"

const (
	// BG/Window共通関数用の用途別ID
	BG     = 0
	Window = 1

	// Paletteの種類
	BGP  = 0
	OBP0 = 1
	OBP1 = 2
)

type PPU struct {
	// LCD
	vp [160 * 144]Pixel // Raw
	FB []byte           // 色変換済み

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
	cycles int // 実機は再現してない

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

// 1ByteをXFlip（水平反転）したときの結果を､負荷軽減のためLUT化する｡（Obj用）
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
	// PPU Disableの場合
	if p.lcdcBits[7] == 0 {
		p.prevLY = p.ly
		p.ly = 0
		p.wly = 0
		return
	}

	// LYC==LY Int監視
	p.stat = (p.stat &^ byte(1<<2)) // LYC==LYビットクリア
	if p.lyc == p.ly {
		p.stat |= 1 << 2
		isLYCIntSel := (p.stat & byte(1<<6)) != 0
		if isLYCIntSel && !p.isLockLYCInt { // LYC == LY && LYCint有効なら､STAT IRQ（初回のみ）
			p.HasSTATIRQ = true
			p.isLockLYCInt = true
		}
	}

	switch {
	// LY 0~143 の間は PPU Mode 2>3>0 LY++ 2>3>0 LY++ を繰り返す
	case p.ly <= 143 && p.ly != p.prevLY: // LYが変わった初回に一気に処理
		p.setMode(2)
		p.oamSearch()
		p.setMode(3)
		p.pixelTransfer()
		p.setMode(0)

	// LY 144~153は､ ずっとPPU Mode1 (VBlank)
	case p.ly == 144 && p.ly != p.prevLY: // LYが変わった初回に一気に処理
		p.setMode(1)
		p.HasVBlankIRQ = true
		p.wly = 0
	}
	p.prevLY = p.ly

	// Mode Int監視
	p.handleModeIRQ()

	// 456cycles経過したら､LY++
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

// FBをLYの1行分だけ更新
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

// STATの下位2bitに現在のPPU Modeをセットする
func (p *PPU) setMode(nextMode int) {
	p.stat = (p.stat & 0x7C) | (byte(nextMode) & 0x03)
	p.isLockModeInt = false
}

// PPU Modeに応じて、必要ならSTAT IRQを出す
func (p *PPU) handleModeIRQ() {
	mode := p.stat & 0x03
	// 現在のモードに対応したMode int selectのbitが立っていて、
	// かつ、前回からモードが変更されたか、Mode int selectが0から1に変わった初回のみ、STAT IRQ。
	bit := p.stat&(1<<(mode+3)) != 0
	if bit && !p.isLockModeInt {
		// STAT IRQ
		p.HasSTATIRQ = true
		p.isLockModeInt = true
	}
}

// 現在のLYで描画対象のObjectを最大10個リストアップする｡
func (p *PPU) oamSearch() {
	p.objList = p.objList[:0] // =[]int{}

	bigMode := p.lcdcBits[2] == 1
	ly := int(p.ly)

	// 選定の優先順位は完全にOAM昇順
	for i := 0; i < 40; i++ {

		y0 := int(p.oam[i<<2+0]) - 16 // Obj(i)の､Viewport上のY始点

		var objHeight int
		if bigMode {
			objHeight = 16
		} else {
			objHeight = 8
		}

		// Obj(i)のY座標が､LYにかかっていればリストに加える
		if ly >= y0 && ly < (y0+objHeight) {
			p.objList = append(p.objList, i)
			// 10個の時点で終了
			if len(p.objList) == 10 {
				return
			}
		}
	}
}

// 描画対象のObjectを､LY行分だけvpに描く
func (p *PPU) objectsTransfer() {
	ly := int(p.ly)

	// Objectの描画リストをX降順（同位ならOAMindex降順）で並び替える｡
	// （描画の優先順位は昇順なので､描画順としては降順になる）
	sortedList := []int{}
	for _, oami := range p.objList {
		pos := len(sortedList) // 並び替え後のスライスのindex(初期値は最後尾)
		for j, oamj := range sortedList {
			if p.oam[oami<<2+1] >= p.oam[oamj<<2+1] { // >=なのは､Xが同じの場合でも､OAMは大きいはずだから）
				pos = j
				break // 後ろはすべてx降順で勝ってるのでやらない
			}
		}
		sortedList = util.InsertSlice(sortedList, pos, oami) // ねじ込み処理
	}

	// 並び替えたリストでFBに書き込んでいく
	for _, oamIdx := range sortedList {

		// Objectごとに、4byte分のObject属性を取り出す
		base := uint16(oamIdx << 2)
		y0 := int(p.oam[base]) - 16
		x0 := int(p.oam[base+1]) - 8
		idx := p.oam[base+2]
		attr := p.oam[base+3]
		tilePy := ly - y0

		data := [2]byte{}
		data = p.getObjectTile(idx, attr, tilePy)
		pri := (attr & (1 << 7)) >> 7

		var pl uint8
		if (attr & (1 << 4)) == 0 {
			pl = OBP0
		} else {
			pl = OBP1
		}

		// 取得したタイルの1行をFBに貼る
		tgtY := ly * 160
		for b := 0; b < 8; b++ {
			tgtX := x0 + b
			if tgtX < 0 || tgtX >= 160 {
				continue
			}
			tgt := tgtY + tgtX
			bgwNum := p.vp[tgt].num

			// Objectが下に隠れる場合
			if pri == 1 && bgwNum != 0 {
				continue
			}
			// 4階調モノクロ（2bit）
			lo := (data[0] >> (7 - b)) & 1
			hi := (data[1] >> (7 - b)) & 1
			num := (hi << 1) | lo

			// Objectは0のピクセルは透明色
			if num == 0 {
				continue
			}

			p.vp[tgt].pl = pl
			p.vp[tgt].num = num
		}
	}
}

// TileDataをObjectとして､tilePyの行だけ取得する
func (p *PPU) getObjectTile(idx, attr byte, tilePy int) [2]byte {
	isYFlip := attr&(1<<6) != 0
	isXFlip := attr&(1<<5) != 0
	if p.lcdcBits[2] == 1 {
		idx &= 0xFE // OBJ size が 8x16 のときはindexは偶数のみ
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

// BackgroundのLY行分だけ､FBに描く
func (p *PPU) BGTransfer() {
	var data [2]byte
	vpIdxY := int(p.ly) * 160
	bgY := int(p.ly + p.scy)
	bgRow := bgY >> 3    // =/8
	tilePy := bgY & 0x07 // =%8
	for x := 0; x < 160; x++ {
		bgX := byte(x) + p.scx // wrapされる
		tilePx := bgX & 0x07   // =%8

		// Viewportの左端か､MapTileが変わるときだけ､tileDataを取得する（節約）
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

// WindowのLY行分だけ､FBに描く
func (p *PPU) WindowTransfer() bool {
	var isDrawn bool
	var data [2]byte
	vpIdxY := int(p.ly) * 160
	wRow := p.wly >> 3     // =/8
	tilePy := p.wly & 0x07 // =%8
	for x := 0; x < 160; x++ {
		bgX := x - (int(p.wx) - 7)
		tilePx := bgX & 0x07 // =%8

		// Viewportの左端か､MapTileが変わるときだけ､tileDataを取得する（節約）
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

// Map(idx)のTileDataをBG/Window MapTileとして､pyの行だけ取得する
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
	tileAddr := tileStart + (uint16(py*2) & 0x1FFF) // Out of range対策
	return [2]byte{
		p.vram[p.bank][tileAddr],
		p.vram[p.bank][tileAddr+1],
	}
}

// vpを色変換してFBに描く
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
	offset := addr & 0x1FFF // Out of range対策
	return p.vram[p.bank][offset]
}

func (p *PPU) WriteVRAM(addr uint16, val byte) {
	offset := addr & 0x1FFF // Out of range対策
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

func (p *PPU) DMATransfer(data *[160]byte) {
	p.oam = *data
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
	p.handleModeIRQ() // ModeごとにSTATIRQ判定
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
