# GoMeBoy

## Description

This is a **Game Boy emulator** written in **Go**, using **Ebiten**.  
It was developed with assistance from **ChatGPT**.

⚠️ This emulator is still a work in progress and contains many bugs.  
🔇 Sound is **not supported**.  
🎮 Only **DMG** and **MBC1** cartridges are supported.

---

![GoMeBoy thumbnail](thumbnail.png)

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

## Passed ROMs

- Blargg's test ROMs  
  - cpu_instrs/01-specials.gb  
  - cpu_instrs/02-interrupts.gb  
  - instr_timing/instr_timing.gb

---

## Partially Working / Failed ROMs

- △ mattcurrie/dmg-acid2.gb
- △ Sushi Nights GB
- ✕ Tobu Tobu Girl
