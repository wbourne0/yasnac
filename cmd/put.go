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
var putCmd = &cobra.Command{
	Use:   "put {JOBNAME}",
	Short: "Reads a job's content from a yanac.",
	// Long: ``,

	Run: func(cmd *cobra.Command, args []string) {
		in := getInputFile(cmd)
		defer in.Close()

		err := machine.WriteFile(args[0], in)

		if err != nil {
			fmt.Println("unable to write job:", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(putCmd)

	putCmd.Flags().StringP("file", "f", "-", "File job content should be written to.  Defaults to stdin.")
}
