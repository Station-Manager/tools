package cmd

import (
	"fmt"
	"github.com/Station-Manager/apikey"
	"github.com/spf13/cobra"
	"os"
)

var hashCmd = &cobra.Command{
	Use: "password",
	RunE: func(cmd *cobra.Command, args []string) error {
		pass, err := cmd.Flags().GetString("pass")
		cobra.CheckErr(err)
		if pass == "" {
			return fmt.Errorf("password must be provided")
		}

		hash := apikey.HashSecret(pass)
		fmt.Println(hash)

		return nil
	},
}

func Execute() {
	if err := hashCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	hashCmd.Flags().StringP("pass", "p", "", "password to hash")
}
