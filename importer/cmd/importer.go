package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Station-Manager/adif"
	"github.com/Station-Manager/utils"
)

func importer(filename string) error {
	if workingDir == "" {
		return fmt.Errorf("working directory not set")
	}

	if err := checkConfig(workingDir); err != nil {
		return err
	}

	if err := checkDbDir(workingDir); err != nil {
		return err
	}

	fpath, err := checkAdiFile(workingDir, filename)
	if err != nil {
		return err
	}

	data, err := readFile(fpath)
	if err != nil {
		return err
	}

	adiObj, err := adif.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Println(adiObj)
	//	sessionID, err = s.DatabaseService.GenerateSession()

	return nil
}

func checkConfig(workingDir string) error {
	cfgFile := filepath.Join(workingDir, "config.json")

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
	dbDir := filepath.Join(workingDir, "db")
	exists, err := utils.PathExists(dbDir)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("database directory not found at %s", dbDir)
	}

	return nil
}

func checkAdiFile(workingDir, filename string) (string, error) {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}

	fpath := filepath.Join(workingDir, filename)
	exists, err := utils.PathExists(fpath)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("ADI file not found at %s", fpath)
	}

	fmt.Printf("Importing file: %s\n", fpath)

	return fpath, nil
}

func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	var data []byte
	if data, err = io.ReadAll(file); err != nil {
		return nil, err
	}

	return data, nil
}
