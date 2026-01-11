package main

import (
	"bytes"
	"fmt"
	"gomeboy/config"
	"gomeboy/internal/emulator"
	"image/color"
	"log"
	"os"
	"strings"

	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var isShowDebug bool
var cfg *config.Config

type Game struct {
	emu      *emulator.Emulator
	img      []byte
	palette  [4][4]byte
	font     *text.GoTextFaceSource
	romTitle string
}

func NewGame(g *Game, rom, sav []byte) *Game {
	font, _ := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))

	imgSize := 160 * 144 * 4
	if isShowDebug {
		imgSize *= 2
	}

	g.emu = emulator.NewEmulator(rom, sav)
	g.img = make([]byte, imgSize)
	g.font = font
	g.palette = [4][4]byte{
		{255, 255, 128, 255},
		{160, 192, 64, 255},
		{64, 128, 64, 255},
		{0, 24, 0, 255},
	}
	// Pass Joypad-related settings in config.toml
	g.emu.CPU.Bus.Joypad.SetIsGamepadEnabled(cfg.Gamepad.IsEnabled)
	g.emu.CPU.Bus.Joypad.SetIsGamepadBind(cfg.Gamepad.Bind)
	return g
}

func (g *Game) Update() error {
	g.setWindowTitle()
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
	// Internal Resolution (Stretch to fit window size while maintaining aspect ratio)
	screenHeight := 144
	screenWidth := 160
	if isShowDebug {
		screenWidth *= 2
	}
	return screenWidth, screenHeight
}

func main() {
	// Load config.toml settings
	var err error
	if cfg, err = config.Load("config.toml"); err != nil {
		panic(err)
	}
	scale := cfg.Video.Scale
	scale = max(scale, 0)
	scale = min(scale, 4)
	isShowDebug = cfg.Video.IsShowDebug

	// Load .gb file
	if len(os.Args) < 2 {
		fmt.Println("usage: gomeboy <romfile>")
		return
	}
	romPath := os.Args[1]
	rom, err := os.ReadFile(romPath)
	if err != nil {
		log.Fatal(err)
	}

	// Load .sav file (if exists)
	savPath := getSavePathFromROM(romPath)
	sav, _ := os.ReadFile(savPath)

	// Set window size
	windowHeight := 144 * scale
	windowWidth := 160 * scale
	if isShowDebug {
		windowWidth *= 2
	}
	ebiten.SetWindowSize(windowWidth, windowHeight)

	g := &Game{}

	// Get ROM title
	s := string(rom[0x0134:0x0144])
	firstNullIdx := strings.IndexByte(s, 0)
	if firstNullIdx != -1 {
		s = s[:firstNullIdx]
	}
	g.romTitle = s

	if err := ebiten.RunGame(NewGame(g, rom, sav)); err != nil && err != ebiten.Termination {
		panic(err)
	} else {
		// When the emulator is closed, save ERAM(save) data
		data := make([]byte, len(g.emu.CPU.Bus.Memory.ERAM))
		copy(data[:], g.emu.CPU.Bus.Memory.ERAM[:])
		os.WriteFile(savPath, data, 0644)
	}
}

// Output to the right of the screen.
func (g *Game) drawDebugMonitor(screen *ebiten.Image) {
	strs := []string{} // Max 20chars * 18rows
	state := string("     ")
	if g.emu.IsPaused {
		state = "PAUSE"
	}
	top := fmt.Sprintf(state+"        FPS:%3.0f", ebiten.ActualFPS())
	strs = append(strs, top)
	strs = append(strs, g.emu.CPU.Tracer.GetCPUInfo()...)
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

// Also converts to RGBA
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

// Use instead of ebiten.DebugPrint
func (g *Game) drawText(dst *ebiten.Image, msg string, x, y, size int, cr color.RGBA) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(cr)
	text.Draw(dst, msg, &text.GoTextFace{
		Source: g.font,
		Size:   float64(size),
	}, op)
}

func getSavePathFromROM(romPath string) string {
	ext := filepath.Ext(romPath)
	base := romPath[:len(romPath)-len(ext)]
	return base + ".sav"
}

func (g *Game) setWindowTitle() {
	emuState := ""
	if g.emu.IsPaused {
		emuState = "(paused)"
	}
	if len(g.romTitle) > 0 {
		ebiten.SetWindowTitle(emuState + "GOmeBoy - " + g.romTitle)
	} else {
		ebiten.SetWindowTitle(emuState + "GOmeBoy")
	}
}
