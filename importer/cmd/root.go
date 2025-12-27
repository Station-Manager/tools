package cmd

import (
	"os"

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

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}

func init() {
	workingDir = mustGetwd()
	importCmd.Flags().StringP("file", "f", "", "ADIF file to import")
}
