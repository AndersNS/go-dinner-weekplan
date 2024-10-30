package main

import (
	"bytes"
	"fmt"
	"math/rand"
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
	folderPath := "./example_files" // TODO: make configurable

	recipes, err := processMarkdownFiles(folderPath)
	if err != nil {
		fmt.Printf("error processing markdown files: %v\n", err)
	}

	fmt.Printf("Recipes found: %v\n", len(recipes))

	fmt.Printf("Your weely menu: \n")
	weekMeny := selectRandomRecipes(recipes, 7)
	for i, r := range weekMeny {
		fmt.Printf("%v: %s\n", i+1, r.name)
	}
}

func selectRandomRecipes(recipes []Recipe, n int) []Recipe {
	temp := make([]Recipe, len(recipes))
	copy(temp, recipes)

	// Fisher-Yates shuffle 👯
	for i := 0; i < n; i++ {
		j := rand.Intn(len(temp)-i) + i
		temp[i], temp[j] = temp[j], temp[i]
	}

	return temp[:n]
}

func processMarkdownFiles(folderPath string) ([]Recipe, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	recipes := make([]Recipe, 0)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			filePath := filepath.Join(folderPath, file.Name())
			recipe, err := processMarkdownFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("error processing file %s: %v\n", file.Name(), err)
			}

			recipes = append(recipes, *recipe)
		}
	}

	return recipes, nil
}

func processMarkdownFile(filePath string) (*Recipe, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	base := filepath.Base(filePath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	fmt.Printf("Processing file: %s\n", name)

	frontMatter, err := getFrontMatter(content)
	if err != nil {
		return nil, fmt.Errorf("error extracting properties: %v", err)
	}

	recipe := Recipe{name: name, frontMatter: frontMatter}

	return &recipe, nil
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
