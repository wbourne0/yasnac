/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// holdCmd represents the hold command
var holdCmd = &cobra.Command{
	Use:   "hold",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := machine.RunCommand("HOLD 0"); err != nil {
			fmt.Println("unable to run job:", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(holdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// holdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// holdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
