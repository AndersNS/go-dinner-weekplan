package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

type FrontMatter struct {
	Tags []string
}

type Recipe struct {
	name        string
	frontMatter FrontMatter
}

func main() {
	folderPath := "./example_files" // Replace with your folder path

	err := processMarkdownFiles(folderPath)
	if err != nil {
		fmt.Printf("Error processing markdown files: %v\n", err)
	}
}

func processMarkdownFiles(folderPath string) error {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("Error reading directory: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			filePath := filepath.Join(folderPath, file.Name())
			err := processMarkdownFile(filePath)
			if err != nil {
				fmt.Printf("Error processing file %s: %v\n", file.Name(), err)
			}
		}
	}

	return nil
}

func processMarkdownFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Error reading file: %v", err)
	}

	base := filepath.Base(filePath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	fmt.Printf("Processing file: %s\n", name)

	properties, err := getFrontMatter(content)
	if err != nil {
		return fmt.Errorf("Error extracting properties: %v", err)
	}

	fmt.Printf("Tags found: %v\n", properties.Tags)

	return nil
}

func getFrontMatter(content []byte) (FrontMatter, error) {
	var properties FrontMatter

	// YAML front matter
	if bytes.HasPrefix(content, []byte("---\n")) {
		endIndex := bytes.Index(content[3:], []byte("\n---"))
		if endIndex != -1 {
			frontMatter := content[3 : endIndex+3]
			err := yaml.Unmarshal(frontMatter, &properties)
			if err != nil {
				return properties, fmt.Errorf("error parsing YAML front matter: %v", err)
			}
			return properties, nil
		}
	}

	// TOML front matter
	if bytes.HasPrefix(content, []byte("+++\n")) {
		endIndex := bytes.Index(content[3:], []byte("\n+++"))
		if endIndex != -1 {
			frontMatter := content[3 : endIndex+3]
			err := toml.Unmarshal(frontMatter, &properties)
			if err != nil {
				return properties, fmt.Errorf("error parsing TOML front matter: %v", err)
			}
			return properties, nil
		}
	}

	return properties, fmt.Errorf("no valid front matter found")
}
