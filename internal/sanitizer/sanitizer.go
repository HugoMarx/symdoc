package sanitizer

import (
	"os"
	"regexp"
)

func SanitizeContent(filePath string) error {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	sanitizedContent := stripAdmonitions(stripAnchors(fileContent))
	if err := os.WriteFile(filePath, sanitizedContent, os.FileMode(os.O_TRUNC)); err != nil {
		return err
	}

	return nil
}

func stripAdmonitions(content []byte) []byte {
	regex := regexp.MustCompile(`(?s):::.*?:::\n?`)
	return regex.ReplaceAll(content, []byte(""))
}

func stripAnchors(content []byte) []byte {
	regex := regexp.MustCompile(`{#[^}]+\}`)
	return regex.ReplaceAll(content, []byte(""))
}
