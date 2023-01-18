/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


var getPosCmd = &cobra.Command{
	Use:   "get-pos",
	Short: "Gets the robot's current position",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := machine.GetPosition()

		if err != nil {
			fmt.Println("Unable to get position:", err.Error())
			os.Exit(1)
		}

		// From section 6.9.8.2.4 of yasnac erc communications
		// http://spaz.org/~jake/robot/479236-17-Communications.pdf
		// The response data is in the following format:
		// Data(1): X axis (mm)
		// Data(2): Y axis (mm)
		// Data(3): Z axis (mm)
		// Data(4): List angle TX (°)
		// Data(5): List angle TY (°)
		// Data(6): List angle TZ (°)
		// Data(7): Type 1 ("0" flip; "1" no flip)
		// Data(8): Type 2 ("0" upper arm; "1" lower arm)
		// Data(9): Type 3 ("0" front; "1" back)
		// Data(10): Pulse number of 7th axis (Traverse axis indicated in mm)
		// Data(11): Pulse number of 8th axis (Traverse axis indicated in mm)
		// Data(12): Pulse number of 9th axis (Traverse axis indicated in mm)
		// Data(13): Pulse number of 10th axis
		// Data(14): Pulse number of 11th axis
		// Data(15): Pulse number of 12th axis

		fmt.Println(result)

		q.Z += 100

		fmt.Println(q.String())

		if err = machine.MoveTo(q, 100); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(getPosCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getPosCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getPosCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
