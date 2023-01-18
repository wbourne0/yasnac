/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get {JOBNAME}",
	Short: "Reads a job's content from a yanac.",
	// Long: ``,

	Run: func(cmd *cobra.Command, args []string) {
		out := getOutputFile(cmd)

		err := machine.ReadFileTo(args[0], out)

		if err != nil {
			fmt.Println("unable to read job:", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringP("out", "o", "-", "File job content should be written to.  Defaults to stdout.")
}
