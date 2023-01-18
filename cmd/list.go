/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists jobs on yasnac",
	// Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := machine.RunCommand("RJDIR *")

		if err != nil {
			fmt.Println("Unable to list jobs:", err.Error())
			os.Exit(1)
		}

		items := strings.Split(string(body), ",")

		fmt.Println(strings.Join(items, "\n"))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
