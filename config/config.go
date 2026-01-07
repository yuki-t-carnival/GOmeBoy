package config

import "github.com/BurntSushi/toml"

func Load(path string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

type Config struct {
	Video   VideoConfig   `toml:"video"`
	Gamepad GamepadConfig `toml:"gamepad"`
}

type VideoConfig struct {
	Scale       int  `toml:"scale"`
	IsShowDebug bool `toml:"show_debug"`
}

type GamepadConfig struct {
	IsEnabled bool   `toml:"enabled"`
	Bind      [8]int `toml:"bind"`
}
