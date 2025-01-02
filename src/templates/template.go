package templates

import (
	"fmt"
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

	dstDir := data.ProjectName

	err := os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		fmt.Errorf("Error creating project directory: %v", err)
	}

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			relativePath := strings.TrimPrefix(path, srcDir)
			dstPath := filepath.Join(dstDir, relativePath)
			return os.MkdirAll(dstPath, os.ModePerm)
		}

		relativePath := strings.TrimPrefix(path, srcDir)
		dstPath := filepath.Join(dstDir, strings.TrimSuffix(relativePath, ".tmpl"))

		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return err
		}

		outputFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}

		defer outputFile.Close()

		return tmpl.Execute(outputFile, data)
	})

	if err != nil {
		fmt.Errorf("error generating project: %v", err)
	}

	fmt.Println("Project created successfully at ", dstDir)
}
