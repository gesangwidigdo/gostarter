package templates

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

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
		srcDir = "templates/gin"
	} else {
		srcDir = "templates/other"
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

	// Join the srcDir with the current working directory
	absSrcDir, err := filepath.Abs(filepath.Join("./src", srcDir))
	if err != nil {
		log.Fatalf("Failed to get absolute path for %s: %v", srcDir, err)
	}

	// Check if the srcDir exists
	if _, err := os.Stat(absSrcDir); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %s\n", absSrcDir)
	}

	files := template.New("")
	err = filepath.Walk(absSrcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".tmpl") {
			relPath, err := filepath.Rel(absSrcDir, path)
			if err != nil {
				return err
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			tmpl := files.New(relPath)
			if _, err := tmpl.Parse(string(content)); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to walk through templates files: %v", err)
	}

	for _, tmpl := range files.Templates() {
		targetPath := filepath.Join(projectDir, tmpl.Name())

		// Create the target directory if it doesn't exist
		targetDir := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create target directory: %v", err)
		}

		targetFile, err := os.Create(targetPath)
		if err != nil {
			log.Fatalf("Failed to create target file: %v", err)
		}
		defer targetFile.Close()

		if err := tmpl.Execute(targetFile, data); err != nil {
			log.Fatalf("Failed to execute template: %v", err)
		}

		fmt.Printf("Generated file: %s\n", targetPath)
	}
}
