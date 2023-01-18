/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	// "fmt"
	// "fmt"
	"time"
	"yasnac/lib"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		machine.RunCommand("DELETE DOUTS")
		j := lib.Job{
			Name: "DOUTS",
			Date: time.Now(),
		}

		// j.Steps = append(j.Steps, "*1")

		for i := 1; i <= 24; i++ {
			j.Dout(i, true)
			// j.Pause()
			j.Steps = append(j.Steps, "TIMER T=0.25")
			// j.Dout(i, false)

			if i % 8 == 0 {
				j.Pause()
			}
		}
		
		// for i := 3; i <= 3; i++ {
		// 	for v := 0; v < 256; v = (v << 1) + 1 {
		// 		j.Steps = append(j.Steps, fmt.Sprintf("DOUT OG#%d %d", i, v))
		// 		// if v == 0 {
		// 		// 	j.Steps = append(j.Steps, "TIMER T=1")
		// 		// 	j.Steps = append(j.Steps, fmt.Sprintf("DOUT OG#%d %d", i, 255))
		// 		// 	j.Steps = append(j.Steps, "TIMER T=1")
		// 		// 	j.Steps = append(j.Steps, fmt.Sprintf("DOUT OG#%d %d", i, 0))
		// 		// 	j.Steps = append(j.Steps, "TIMER T=1")
		// 		// } else {
		// 		// j.Steps = append(j.Steps, "TIMER T=1")
		// 		// }
		// 	}

		// 	// j.Dout(i, true)
		// 	// j.Pause()
		// 	// j.Steps = append(j.Steps, "TIMER T=0.25")
		// 	// j.Dout(i, false)

		// }

		// j.Steps = append(j.Steps, "JUMP *1")

		buf := j.Buffer()
		if err := machine.WriteFile("DOUTS", &buf); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
