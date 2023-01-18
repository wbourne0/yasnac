/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"time"
	"yasnac/lib"

	"github.com/256dpi/gcode"
	"github.com/spf13/cobra"
)



var (
	bl = lib.Position{X: -34.54, Y: 762.54, Z: 45.29, TX: 95.34, TY: -86.42, TZ: 177.68}
	// bl = lib.Position{X: -105.49, Y: 513.85, Z: 459.82, TX: 180, TY: 0, TZ: 0}
	// tr = lib.Position{X: 394.51, Y: 1013.85, Z: 459.82, TX: 180, TY: 0, TZ: 0}
)

// gcodeCmd represents the gcode command
var gcodeCmd = &cobra.Command{
	Use:   "gcode",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// bl.Y += 100
		bl.Z -= 3
		in := getInputFile(cmd)
		defer in.Close()
		out := getOutputFile(cmd)
		defer out.Close()

		parsed, err := gcode.ParseFile(in)

		if err != nil {
			fmt.Printf("Unable to parse file: %s\n", err.Error())

			os.Exit(1)
		}

		job := lib.Job{
			Name: "GCODE",
			Date: time.Now(),
		}



		// job.MoveTo(bl)

		pos := bl

		for _, l := range parsed.Lines {
			if len(l.Codes) < 2 || l.Codes[0].Letter != "G" || l.Codes[0].Value != 1 {
				continue
			}

			didChange := false



			for _, v := range l.Codes[1:] {
				switch v.Letter{
				case "X":
					didChange = true
					pos.X = bl.X + v.Value * 2
				case "Y":
					didChange = true
					pos.Y = bl.Y + v.Value * 2
				case "Z":
					didChange = true
					pos.Z = bl.Z + v.Value
				}
			}

			if !didChange {
				continue
			}

			if job.MoveTo(pos, 3) {
				break
			}
		}

		// pos.Z += 50
		job.MoveTo(pos, 3)

		buf := job.Buffer()

		if _, err = io.Copy( out, &buf); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(gcodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gcodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gcodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func init() {
	rootCmd.AddCommand(putCmd)

	gcodeCmd.Flags().StringP("file", "f", "-", "Source file.  Defaults to stdin.")
	gcodeCmd.Flags().StringP("out", "o", "-", "Output file.  Defaults to stdout.")
}