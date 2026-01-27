package main

import (
	"bytes"
	"fmt"
	"gomeboy/config"
	"gomeboy/internal/apu"
	"gomeboy/internal/emulator"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"time"

	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var screenFont *text.GoTextFaceSource

type Game struct {
	emu                  *emulator.Emulator
	ebitenImage          *ebiten.Image
	imageRGBA            *image.RGBA
	audioCtx             *audio.Context
	audioPlayer          *audio.Player
	cfg                  *config.Config
	pixelScale           int
	isDebugScreenEnabled bool
	debugLog             []string
}

func newGame(g *Game, rom, sav []byte) *Game {
	screenFont, _ = text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))

	debuggerWidth := 0
	if g.isDebugScreenEnabled {
		debuggerWidth = 160
	}
	g.imageRGBA = image.NewRGBA(image.Rect(0, 0, 160+debuggerWidth, 144))
	g.ebitenImage = ebiten.NewImage(160+debuggerWidth, 144)

	g.emu = emulator.NewEmulator(rom, sav)

	g.emu.CPU.Bus.Joypad.SetIsGamepadEnabled(g.cfg.Gamepad.IsEnabled)
	g.emu.CPU.Bus.Joypad.SetIsGamepadBind(g.cfg.Gamepad.Bind)

	g.audioCtx = audio.NewContext(int(apu.SampleRate))
	g.audioPlayer, _ = g.audioCtx.NewPlayerF32(g.emu.CPU.Bus.APU.AudioStream)
	g.audioPlayer.SetBufferSize(40 * time.Millisecond)
	g.audioPlayer.SetVolume(0.5)
	g.audioPlayer.Play()

	return g
}

// Game.Update() calls Emulator.RunFrame() at 60FPS.
func (g *Game) Update() error {
	g.setWindowTitle()
	if !ebiten.IsFocused() || g.emu.IsPauseMode {
		g.audioPlayer.Pause()
	} else {
		g.audioPlayer.Play()
	}
	if ebiten.IsFocused() {
		if g.emu.RunFrame() == -1 {
			return ebiten.Termination
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	gameScreen := g.emu.CPU.Bus.PPU.GetGameScreen()
	draw.Draw(g.imageRGBA, image.Rect(0, 0, 160, 144), gameScreen, gameScreen.Rect.Min, draw.Src)
	g.ebitenImage = ebiten.NewImageFromImage(g.imageRGBA)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(g.pixelScale), float64(g.pixelScale))
	screen.DrawImage(g.ebitenImage, op)

	if g.isDebugScreenEnabled {
		strs := g.emu.GetDebugLog()
		for i, s := range strs {
			white := color.RGBA{255, 255, 255, 255}
			fontSize := 16
			g.drawText(screen, s, 160*g.pixelScale+fontSize, (i+1)*fontSize, fontSize, white)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	screenHeight := 144 * g.pixelScale
	screenWidth := 160 * g.pixelScale
	if g.isDebugScreenEnabled {
		screenWidth *= 2
	}

	return screenWidth, screenHeight
}

func main() {
	g := &Game{}

	var err error
	if g.cfg, err = config.Load("config.toml"); err != nil {
		panic(err)
	}
	g.pixelScale = g.cfg.Video.Scale
	g.pixelScale = max(g.pixelScale, 1)
	g.pixelScale = min(g.pixelScale, 4)
	g.isDebugScreenEnabled = g.cfg.Video.IsShowDebug

	if len(os.Args) < 2 {
		fmt.Println("usage: gomeboy <romfile>")
		return
	}
	romPath := os.Args[1]
	rom, err := os.ReadFile(romPath)
	if err != nil {
		log.Fatal(err)
	}

	savPath := getSavePathFromROM(romPath)
	sav, _ := os.ReadFile(savPath)

	windowHeight := 144 * g.pixelScale
	windowWidth := 160 * g.pixelScale
	if g.isDebugScreenEnabled {
		windowWidth *= 2
	}
	ebiten.SetWindowSize(windowWidth, windowHeight)

	err = ebiten.RunGame(newGame(g, rom, sav))
	if err != nil && err != ebiten.Termination {
		panic(err)
	} else {
		// When the emulator is closed, save ERAM(save) data.
		savData := g.emu.CPU.Bus.Memory.GetSaveData()
		os.WriteFile(savPath, savData, 0644)
	}
}

func (g *Game) drawText(dst *ebiten.Image, msg string, x, y, size int, cr color.RGBA) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(cr)
	text.Draw(dst, msg, &text.GoTextFace{
		Source: screenFont,
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
	if len(g.emu.ROMTitle) > 0 {
		ebiten.SetWindowTitle(emuState + "GOmeBoy - " + g.emu.ROMTitle)
	} else {
		ebiten.SetWindowTitle(emuState + "GOmeBoy")
	}
}
