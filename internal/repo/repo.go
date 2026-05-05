package repo

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"symdoc/internal/config"
)

var fileCount int = 0

func BuildMarkdown(root string) error {
	if err := filepath.WalkDir(os.ExpandEnv(root), rstToMd); err != nil {
		return err
	}

	configObj := config.Get()
	configObj.MdGenerated = true
	configObj.MdGeneratedAt = fmt.Sprint(time.Now().Unix())
	config.Save(configObj)

	log.Println(fileCount, " files converted to markdown.")
	log.Println("Config file updated !")
	return nil
}

func rstToMd(path string, d fs.DirEntry, err error) error {
	if filepath.Ext(path) == ".rst" {
		markdownFile := strings.Replace(d.Name(), ".rst", ".md", 1)
		if err := exec.Command("pandoc", path, "-f", "rst", "-t", "markdown", "-s", "-o", fmt.Sprint(filepath.Join(filepath.Dir(path), markdownFile))).Run(); err != nil {
			return err
		}
		fmt.Println(fmt.Sprint(filepath.Dir(path), "/", markdownFile, " created."))
		fileCount++

	}
	return nil
}
