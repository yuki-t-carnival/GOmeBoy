# GOmeBoy

## Description

This is a **Game Boy emulator** written in **Go**, using **Ebiten**.  
It was developed with assistance from **ChatGPT**.

⚠️ This emulator is still a work in progress and contains many bugs.  
🔇 Sound is **not supported**.  
🎮 Only **DMG** and **MBC1** cartridges are supported.

![GOmeBoy thumbnail](thumbnail.png)

---

## How to Launch

    go run ./cmd/gomeboy/main.go <rom_path>

---

## How to Change Settings

Edit the configuration file:

    ./config.toml

---

## Default Game Button Bindings

| Game Boy | Key |
|----------|-----|
| A        | Z |
| B        | X |
| SELECT   | Left Shift |
| START    | Enter |
| D-Pad    | Arrow Keys |

---

## Emulator Control Keys

| Action | Key |
|--------|-----|
| Toggle Pause / Run | P |
| Step (while paused) | S |
| Exit | Esc |

---

## Playable / Passed ROMs

- Tobu Tobu Girl
- Sushi Nights GB
- mattcurrie/dmg-acid2.gb
- Blargg's test ROMs  
  - cpu_instrs.gb    
  - instr_timing.gb
