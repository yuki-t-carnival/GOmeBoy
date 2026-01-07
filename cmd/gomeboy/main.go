package main

import (
	"bytes"
	"fmt"
	"gomeboy/config"
	"gomeboy/internal/emulator"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var isShowDebug bool
var cfg *config.Config

type Game struct {
	emu     *emulator.Emulator
	img     []byte
	palette [4][4]byte
	font    *text.GoTextFaceSource
}

func NewGame(rom []byte) *Game {
	font, _ := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))

	imgSize := 160 * 144 * 4
	if isShowDebug {
		imgSize *= 2
	}
	g := &Game{
		emu:  emulator.NewEmulator(rom),
		img:  make([]byte, imgSize),
		font: font,
		palette: [4][4]byte{
			{255, 255, 255, 255},
			{191, 191, 191, 255},
			{127, 127, 127, 255},
			{0, 0, 0, 255},
		},
	}
	// Joypadにconfig.tomlのGamepad関連の設定を渡す
	g.emu.CPU.Bus.Joypad.SetIsGamepadEnabled(cfg.Gamepad.IsEnabled)
	g.emu.CPU.Bus.Joypad.SetIsGamepadBind(cfg.Gamepad.Bind)
	return g
}

func (g *Game) Update() error {
	if g.emu.RunFrame() == -1 {
		return ebiten.Termination
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.updateImgFromFB(g.emu.CPU.Bus.PPU.FB)
	screen.WritePixels(g.img)
	if isShowDebug {
		g.drawDebugMonitor(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// 内部解像度（アスペクト比を維持した範囲でウィンドウサイズに伸縮する）
	screenHeight := 144
	screenWidth := 160
	if isShowDebug {
		screenWidth *= 2
	}
	return screenWidth, screenHeight
}

func main() {
	// config.tomlの読み込み
	var err error
	cfg, err = config.Load("config.toml")
	if err != nil {
		panic(err)
	}
	scale := cfg.Video.Scale
	scale = max(scale, 0)
	scale = min(scale, 4)
	isShowDebug = cfg.Video.IsShowDebug

	// ウィンドウの設定
	windowHeight := 144 * scale
	windowWidth := 160 * scale
	if isShowDebug {
		windowWidth *= 2
	}
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("GoMeboy")

	// .gbファイルの読み込み
	if len(os.Args) < 2 {
		fmt.Println("usage: gomeboy <romfile>")
		return
	}
	romPath := os.Args[1]
	rom, err := os.ReadFile(romPath)
	if err != nil {
		log.Fatal(err)
	}

	// NewGame()の後､Update()とDraw()ループへ
	if err := ebiten.RunGame(NewGame(rom)); err != nil {
		panic(err)
	}
}

// ゲーム画面右側にデバッグモニタを出力する
func (g *Game) drawDebugMonitor(screen *ebiten.Image) {
	strs := []string{} // 最大20文字*18行

	state := string("     ")
	if g.emu.IsPaused {
		state = "PAUSE"
	}
	top := fmt.Sprintf(state+"        FPS:%3.0f", ebiten.ActualFPS())
	strs = append(strs, top)
	strs = append(strs, g.emu.CPU.Tracer.GetCPUInfo()...) // CPU情報を受け取る
	strs = append(strs, "")
	strs = append(strs, "Cart: "+g.emu.CPU.Bus.Memory.CartTypeName)
	strs = append(strs, "ROM: "+g.emu.CPU.Bus.Memory.ROMSizeName)
	strs = append(strs, "RAM: "+g.emu.CPU.Bus.Memory.RAMSizeName)

	white := color.RGBA{255, 255, 255, 255}
	red := color.RGBA{255, 0, 0, 255}
	var cr = color.RGBA{}
	for i, s := range strs {
		if i == 0 {
			cr = red
		} else {
			cr = white
		}
		g.drawText(screen, s, 160, i*8, 8, cr)
	}
}

// FB（値=色番号）から、RGBA形式に変換
func (g *Game) updateImgFromFB(fb []byte) {
	srcBase := 0
	dstBase := 0
	for y := 0; y < 144; y++ {
		srcBase = y * 160
		dstBase = srcBase
		if isShowDebug {
			dstBase *= 2
		}
		for x := 0; x < 160; x++ {
			colorNum := fb[srcBase+x]
			dst := (dstBase + x) * 4
			copy(g.img[dst:dst+4], g.palette[colorNum][:])
		}
	}
}

// DebugPrintだと物足りないので､自作する｡
func (g *Game) drawText(dst *ebiten.Image, msg string, x, y, size int, cr color.RGBA) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(cr)
	text.Draw(dst, msg, &text.GoTextFace{
		Source: g.font,
		Size:   float64(size),
	}, op)

}
