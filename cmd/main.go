package main

import (
	"fmt"
	"log"

	"symdoc/internal/config"
	"symdoc/internal/repo"
)

func main() {
	fmt.Println("Hello Symdoc")
	config := config.Load()

	if !config.MdGenerated {
		if err := repo.BuildMarkdown(config.AssetsDir); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Print("Already converted !")
	}
}
