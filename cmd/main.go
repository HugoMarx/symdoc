package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"symdoc/internal/repo"
)

func main() {
	fmt.Println("Hello Symdoc")
	configDir := os.ExpandEnv(filepath.Join("$HOME", ".symdoc", "symfony-docs"))
	if err := repo.BuildMarkdown(configDir); err != nil {
		log.Fatal(err)
	}
}
