package config

import (
	"bytes"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type configuration struct {
	MdGenerated   bool   `toml:"md_generated"`
	MdGeneratedAt string `toml:"md_generated_at"`
	AssetsDir     string `toml:"assets_dir"`
}

var config = configuration{
	MdGenerated: false,
	AssetsDir:   "$HOME/.symdoc/symfony-docs",
}

func Load() configuration {
	configPath := GetConfigPath()
	_, err := toml.DecodeFile(configPath, &config)
	if os.IsNotExist(err) {
		Save(config)
		log.Println("Config file created !")
	}
	return config
}

func GetConfigPath() string {
	return os.ExpandEnv(filepath.Join("$HOME", ".symdoc", "config.toml"))
}

func Get() configuration {
	return config
}

func Save(config configuration) {
	buffer := bytes.NewBuffer([]byte{})
	if err := toml.NewEncoder(buffer).Encode(config); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(GetConfigPath(), buffer.Bytes(), 0o644); err != nil {
		log.Print("Unable to update config file !")
	}
}
