package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"symdoc/internal/config"
	"symdoc/internal/repo"

	"github.com/charmbracelet/glamour"
)

func main() {
	fmt.Println("Hello Symdoc")

	fileName := os.Args[1:][0] // Starts at index 1 since 0 is the program path.
	config := config.Load()

	if !config.MdGenerated {
		if err := repo.BuildMarkdown(config.AssetsDir); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Print("Already converted !")
	}

	fileContent, err := os.ReadFile(filepath.Join(os.ExpandEnv(config.AssetsDir), fileName))
	if err != nil {
		log.Fatal("File not found : ", err.Error())
	}

	output, err := glamour.Render(string(fileContent), "dark")
	if err != nil {
		log.Fatal("Unable to render file content : ", err.Error())
	}

	fmt.Print(output)
}
