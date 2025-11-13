package cmd

import (
	"fmt"
	"github.com/Station-Manager/apikey"
	"github.com/Station-Manager/config"
	"github.com/Station-Manager/database"
	"github.com/Station-Manager/iocdi"
	"github.com/Station-Manager/logging"
	"github.com/Station-Manager/utils"
	"github.com/spf13/cobra"
	"os"
	"reflect"
	"time"
)

var container *iocdi.Container

// var log *logging.Service
// var cfg *config.Service
var db *database.Service

var generateCmd = &cobra.Command{
	Use: "generate",
	RunE: func(cmd *cobra.Command, args []string) error {
		plain, hash, expiresAt, err := apikey.GenerateBootstrap()
		if err != nil {
			return err
		}

		fmt.Println("generate called:")
		fmt.Println("plain:", plain)
		fmt.Println("hash:", hash)
		fmt.Println("expiresAt:", expiresAt.Format(time.RFC3339))

		db.StoreBootstrap()
		return db.Close()
	},
}

func Execute() {
	if err := generateCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	_ = os.Setenv("SM_DEFAULT_DB", "pg")

	workingDir, err := utils.WorkingDir()
	cobra.CheckErr(err)

	container = iocdi.New()
	err = container.RegisterInstance("workingdir", workingDir)
	cobra.CheckErr(err)
	err = container.Register(config.ServiceName, reflect.TypeOf((*config.Service)(nil)))
	cobra.CheckErr(err)
	err = container.Register(logging.ServiceName, reflect.TypeOf((*logging.Service)(nil)))
	cobra.CheckErr(err)
	err = container.Register(database.ServiceName, reflect.TypeOf((*database.Service)(nil)))
	cobra.CheckErr(err)
	err = container.Build()
	cobra.CheckErr(err)

	//cfgInstance, err := container.ResolveSafe(config.ServiceName)
	//cobra.CheckErr(err)
	//var ok bool
	//cfg, ok = cfgInstance.(*config.Service)
	//if !ok {
	//	panic("cfg is not of type *config.Service")
	//}
	//
	//logInstance, err := container.ResolveSafe(logging.ServiceName)
	//cobra.CheckErr(err)
	//log, ok = logInstance.(*logging.Service)
	//if !ok {
	//	panic("db is not of type *database.Service")
	//}

	dbInstance, err := container.ResolveSafe(database.ServiceName)
	cobra.CheckErr(err)
	var ok bool
	db, ok = dbInstance.(*database.Service)
	if !ok {
		panic("db is not of type *database.Service")
	}

	err = db.Open()
	cobra.CheckErr(err)

	err = db.Migrate()
	cobra.CheckErr(err)
}
