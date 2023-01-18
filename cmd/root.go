/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"yasnac/lib"

	"github.com/spf13/cobra"
)

var machine *lib.ERCMachine

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yasnac",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("serial")
		_, err := os.Stat(port)

		if err != nil {
			fmt.Println("Unable to locate serial bus at", port)
			os.Exit(1)
		}

		machine, err = lib.NewERCMachine(port)

		if err != nil {
			fmt.Println("Unable to connect to device:", err.Error())
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("serial", "s", "/dev/ttyUSB0", "Serial port for yasnac.  Defaults to /dev/ttyUSB0")

}
