package cmd

import (
	"os"
	"path/filepath"
	"reflect"

	"github.com/Station-Manager/config"
	"github.com/Station-Manager/database/sqlite"
	"github.com/Station-Manager/iocdi"
	"github.com/Station-Manager/logging"
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
var container *iocdi.Container
var cfg *config.Service
var db *sqlite.Service

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
	workingDir = filepath.Join(workingDir, "..")

	container = iocdi.New()
	err := container.RegisterInstance("WorkingDir", workingDir)
	cobra.CheckErr(err)
	err = container.Register(config.ServiceName, reflect.TypeOf((*config.Service)(nil)))
	cobra.CheckErr(err)
	err = container.Register(sqlite.ServiceName, reflect.TypeOf((*sqlite.Service)(nil)))
	cobra.CheckErr(err)
	err = container.Register(logging.ServiceName, reflect.TypeOf((*logging.Service)(nil)))
	cobra.CheckErr(err)

	err = container.Build()
	cobra.CheckErr(err)

	cfgInstance, err := container.ResolveSafe(config.ServiceName)
	cobra.CheckErr(err)
	var ok bool
	cfg, ok = cfgInstance.(*config.Service)
	if !ok {
		panic("cfg is not of type *config.Service")
	}
	dbInstance, err := container.ResolveSafe(sqlite.ServiceName)
	cobra.CheckErr(err)
	db, ok = dbInstance.(*sqlite.Service)
	if !ok {
		panic("db is not of type *sqlite.Service")
	}

	importCmd.Flags().StringP("file", "f", "", "ADIF file to import")
}
