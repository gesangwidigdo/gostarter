package program

import (
	"fmt"
	"os"
	"os/exec"
)

func ProjectInitialization(projectName, moduleURL string) error {
	if err := os.Mkdir(projectName, 0755); err != nil {
		return fmt.Errorf("error creating project folder: %v", err)
	}

	if err := os.Chdir(projectName); err != nil {
		return fmt.Errorf("error changing directory: %v", err)
	}

	cmd := exec.Command("go", "mod", "init", moduleURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error initializing go module: %v", err)
	}

	return nil
}
