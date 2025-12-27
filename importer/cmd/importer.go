package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Station-Manager/utils"
)

func importer(filename string) error {
	if workingDir == "" {
		return fmt.Errorf("working directory not set")
	}

	if !strings.HasSuffix(workingDir, "tools") {
		return fmt.Errorf("working directory must be in the tools directory")
	}

	if err := checkConfig(workingDir); err != nil {
		return err
	}

	if err := checkDbDir(workingDir); err != nil {
		return err
	}

	filename = strings.TrimSpace(filename)
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	fpath := filepath.Join(workingDir, "..", filename)
	fmt.Printf("Importing file: %s\n", fpath)

	return nil
}

func checkConfig(workingDir string) error {
	cfgFile := filepath.Join(workingDir, "..", "config.json")

	exists, err := utils.PathExists(cfgFile)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("config file not found at %s", cfgFile)
	}
	return nil
}

func checkDbDir(workingDir string) error {
	dbDir := filepath.Join(workingDir, "..", "db")
	exists, err := utils.PathExists(dbDir)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("database directory not found at %s", dbDir)
	}

	return nil
}
