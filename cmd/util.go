package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func getInputFile(cmd *cobra.Command) (in io.ReadCloser) {

	outpath, err := cmd.Flags().GetString("file")

	if err != nil {
		fmt.Println("unable to parse args:", err.Error())
		os.Exit(1)
	}

	if outpath == "-" {
		in = os.Stdin
	} else {
		file, err := os.Open(outpath)

		if err != nil {
			fmt.Printf("unable to open %s: %s\n", outpath, err.Error())
			os.Exit(1)
		}

		in = file
	}

	return
}

func getOutputFile(cmd *cobra.Command) (out io.WriteCloser) {

	outpath, err := cmd.Flags().GetString("out")

	if err != nil {
		fmt.Println("unable to parse args:", err.Error())
		os.Exit(1)
	}

	if outpath == "-" {
		out = os.Stdout
	} else {
		file, err := os.Create(outpath)

		if err != nil {
			fmt.Printf("unable to create %s: %s\n", outpath, err.Error())
			os.Exit(1)
		}

		out = file
	}

	return
}
