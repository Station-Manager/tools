package cmd

import (
	"fmt"
	"path/filepath"
)

func importer(filename string) error {
	if workingDir == "" {
		return fmt.Errorf("working directory not set")
	}

	toolsDir := filepath.Join(workingDir, "tools")
	fmt.Println(">", toolsDir)
	//exists, err := utils.PathExists(toolsDir)
	//if err != nil {
	//	return err
	//}

	//if !exists {
	//	if err = os.MkdirAll(toolsDir, 0755); err != nil {
	//		return err
	//	}
	//}
	//
	//filename = strings.TrimSpace(filename)
	//if filename == "" {
	//	return fmt.Errorf("filename cannot be empty")
	//}

	//	fpath := filepath.Join(workingDir, filename)

	return nil
}

func checkConfig(workingDir string) error {
	//	cfgFile := filepath.Join(workingDir, "..", "config.json")

	return nil
}
