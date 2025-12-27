package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import ADIF files",
	Long:  "Import ADIF files into the Station Manager database.",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}
		if err = importer(file); err != nil {
			return err
		}
		return nil
	},
}

var workingDir string

func Execute() {
	cobra.CheckErr(importCmd.Execute())
}

func exeDir() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	resolved, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		panic(err)
	}
	return filepath.Dir(resolved)
}

func init() {
	workingDir = exeDir()
	importCmd.Flags().StringP("file", "f", "", "ADIF file to import")
}
