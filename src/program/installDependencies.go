package program

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gesangwidigdo/gostarter/src/dependencies"
)

func InstallDependencies(selectedFramework, selectedDBMS string) error {
	var frameworkURL, dbmsURL string
	// Find Framework dependencies
	for _, framework := range dependencies.Frameworks {
		if framework.Name == selectedFramework {
			frameworkURL = framework.URL
			break
		}
	}

	// Find DBMS dependencies
	for _, dbms := range dependencies.DBMSs {
		if dbms.DBMS == selectedDBMS {
			dbmsURL = dbms.URL
			break
		}
	}

	dependencies := []string{
		frameworkURL,
		dbmsURL,
		"gorm.io/gorm",
		"github.com/joho/godotenv",
		"github.com/golang-jwt/jwt/v5",
	}
	for _, dependency := range dependencies {
		cmd := exec.Command("go", "get", "-u", dependency)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error installing dependency: %v", err)
		}
	}

	return nil
}
