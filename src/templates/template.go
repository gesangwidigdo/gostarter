package templates

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed gin/*.tmpl
//go:embed gin/**/*.tmpl
//go:embed other/*.tmpl
//go:embed template.go
var templatesFS embed.FS

type TemplateData struct {
	ProjectName string
	ModuleName  string
	Framework   string
}

func GenerateTemplate(projectName, moduleName, framework string) {
	data := TemplateData{
		ProjectName: projectName,
		ModuleName:  moduleName,
		Framework:   framework,
	}

	var srcDir string
	if data.Framework == "Gin" {
		srcDir = "gin"
	} else {
		srcDir = "other"
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	// Buat folder berdasarkan nama ProjectName
	projectDir := filepath.Join(cwd, projectName)
	if err := os.MkdirAll(projectDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create project directory: %v", err)
	}

	// Rekursif untuk memproses semua file template di srcDir
	err = processDirectory(srcDir, projectDir, data)
	if err != nil {
		log.Fatalf("Error processing directory: %v", err)
	}
}

func processDirectory(srcDir, projectDir string, data TemplateData) error {
	entries, err := templatesFS.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", srcDir, err)
	}

	for _, entry := range entries {
		embeddedPath := path.Join(srcDir, entry.Name())
		targetPath := filepath.Join(projectDir, entry.Name())

		if entry.IsDir() {
			// Jika item adalah direktori, buat direktori target dan proses secara rekursif
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
			if err := processDirectory(embeddedPath, targetPath, data); err != nil {
				return err
			}
		} else if strings.HasSuffix(entry.Name(), ".tmpl") {
			// Jika item adalah file template, baca, parse, dan generate file
			content, err := templatesFS.ReadFile(embeddedPath)
			if err != nil {
				return fmt.Errorf("failed to read embedded file %s: %w", embeddedPath, err)
			}

			tmpl := template.New(entry.Name())
			_, err = tmpl.Parse(string(content))
			if err != nil {
				return fmt.Errorf("failed to parse template %s: %w", embeddedPath, err)
			}

			// Hapus ekstensi .tmpl dari nama file target
			targetFilePath := strings.TrimSuffix(targetPath, ".tmpl")
			targetFile, err := os.Create(targetFilePath)
			if err != nil {
				return fmt.Errorf("failed to create target file %s: %w", targetFilePath, err)
			}
			defer targetFile.Close()

			if err := tmpl.Execute(targetFile, data); err != nil {
				return fmt.Errorf("failed to execute template for %s: %w", targetFilePath, err)
			}

			fmt.Printf("Generated file: %s\n", targetFilePath)
		}
	}

	return nil
}
