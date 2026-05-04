package repo

import (
	"fmt"
	"io/fs"
	"strings"

	"os/exec"
	"path/filepath"
)

var fileCount int = 0

func BuildMarkdown(root string) error {
	if err := filepath.WalkDir(root, rstToMd); err != nil {
		return err
	}

	fmt.Println(fileCount, " files converted to markdown.")
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
