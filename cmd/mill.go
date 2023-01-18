/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"
	"yasnac/lib"

	"github.com/spf13/cobra"
)

// millCmd represents the mill command
var millCmd = &cobra.Command{
	Use:   "mill",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		br := lib.Position{X: 70.44, Y: 913.72, Z: 48.45, TX: 95.34, TY: -86.42, TZ: 177.68}
		tr := lib.Position{X: 285.5, Y: 970.60, Z: 58.45, TX: 95.34, TY: -86.42, TZ: 177.68}

		j := lib.Job{
			Name: "MILL",
			Date: time.Now(),
		}


		var isAlt bool
		const speed = 30
		const dis = 2
		br.Y -= dis * 4

		for y := br.Y; y < tr.Y; y += dis {
			posL := br
			posR := tr

			posL.Y = y
			posR.Y = y

			if isAlt {
				j.MoveTo(posR, speed)
				j.MoveTo(posL, speed)
			} else {
				j.MoveTo(posL, speed)
				j.MoveTo(posR, speed)
			}
			isAlt = !isAlt
		}

		posL := br
		posR := tr

		posL.Y = tr.Y
		posR.Y = tr.Y

		if isAlt {
			j.MoveTo(posR, speed)
			j.MoveTo(posL, speed)
		} else {
			j.MoveTo(posL, speed)
			j.MoveTo(posR, speed)
		}
		isAlt = !isAlt

		j.MoveTo(br, speed)

		buf := j.Buffer()

		fmt.Println(buf.String())

		if _, err := machine.RunCommand("DELETE MILL"); err != nil {
			panic(err)
		}

		if err := machine.WriteFile("MILL", &buf); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(millCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// millCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// millCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
