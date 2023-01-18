/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"
	"time"
	"yasnac/lib"

	"github.com/spf13/cobra"
)

const disX = 3 / 4.0

// const offY = 1 / 4.0
const offZ = 2 // 2-3
const inToMM = 25.4

type KeyOff struct {
	x, y  float64
	shift bool
}

var arr = [...]KeyOff{
	'`':  {-1, -1, false},
	'~':  {-1, -1, true},
	'0':  {9, -1, false},
	')':  {9, -1, true},
	'1':  {0, -1, false},
	'!':  {0, -1, true},
	'2':  {1, -1, false},
	'@':  {1, -1, true},
	'3':  {2, -1, false},
	'#':  {2, -1, true},
	'4':  {3, -1, false},
	'$':  {3, -1, true},
	'5':  {4, -1, false},
	'%':  {4, -1, true},
	'6':  {5, -1, false},
	'^':  {5, -1, true},
	'7':  {6, -1, false},
	'&':  {6, -1, true},
	'8':  {7, -1, false},
	'*':  {7, -1, true},
	'9':  {8, -1, false},
	'(':  {8, -1, true},
	'a':  {0, 1, false},
	'A':  {0, 1, true},
	'b':  {4, 2, false},
	'B':  {4, 2, true},
	'c':  {2, 2, false},
	'C':  {2, 2, true},
	'd':  {2, 1, false},
	'D':  {2, 1, true},
	'e':  {2, 0, false},
	'E':  {2, 0, true},
	'f':  {3, 1, false},
	'F':  {3, 1, true},
	'g':  {4, 1, false},
	'G':  {4, 1, true},
	'h':  {5, 1, false},
	'H':  {5, 1, true},
	'i':  {7, 0, false},
	'I':  {7, 0, true},
	'j':  {6, 1, false},
	'J':  {6, 1, true},
	'k':  {7, 1, false},
	'K':  {7, 1, true},
	'l':  {8, 1, false},
	'L':  {8, 1, true},
	'm':  {6, 2, false},
	'M':  {6, 2, true},
	'n':  {5, 2, false},
	'N':  {5, 2, true},
	'o':  {8, 0, false},
	'O':  {8, 0, true},
	'p':  {9, 0, false},
	'P':  {9, 0, true},
	'q':  {0, 0, false},
	'Q':  {0, 0, true},
	'r':  {3, 0, false},
	'R':  {3, 0, true},
	's':  {1, 1, false},
	'S':  {1, 1, true},
	't':  {4, 0, false},
	'T':  {4, 0, true},
	'u':  {6, 0, false},
	'U':  {6, 0, true},
	'v':  {3, 2, false},
	'V':  {3, 2, true},
	'w':  {1, 0, false},
	'W':  {1, 0, true},
	'x':  {1, 2, false},
	'X':  {1, 2, true},
	'y':  {5, 0, false},
	'Y':  {5, 0, true},
	'z':  {0, 2, false},
	'Z':  {0, 2, true},
	' ':  {4, 3, false},
	',':  {7, 2, false},
	'<':  {7, 2, true},
	'.':  {8, 2, false},
	'>':  {8, 2, true},
	'/':  {9, 2, true},
	'?':  {9, 2, true},
	';':  {9, 1, false},
	':':  {9, 1, true},
	'\'': {10, 1, false},
	'"':  {10, 1, true},
	'[':  {10, 0, false},
	'{':  {10, 0, true},
	']':  {11, 0, false},
	'}':  {11, 0, true},
	'\\': {12, 0, true},
	'|':  {12, 0, true},
	'\t': {-(3 / 4), 0, false},
}

var (
	caps  = KeyOff{-(3.0 / 4.0), 1, false}
	shift = KeyOff{-(1.0 + 1.0/4.0), 2, false}

	// qd = 175.20
	q = lib.Position{X: -69.80, Y: 859.27, Z: 23.56,  TX: 95.34, TY: -86.42, TZ: 177.68}
)

func (a KeyOff) toPos(origin lib.Position) lib.Position {
	origin.X += (disX * a.x) * inToMM

	switch a.y {
	case -1:
		origin.X -= (3.0 / 8.0) * inToMM
	case 1:
		origin.X += (1.0 / 8.0) * inToMM
	case 2:
		origin.X += (5.0 / 8.0) * inToMM
	}

	origin.Y -= (disX * a.y) * inToMM
	origin.Z -= (offZ * a.y)
	fmt.Println(origin.X, origin.Y)
	return origin
}

const (
	speedTravel   = 500
	speedKeypress = 20
)

// const (
// 	speedTravel   = 100
// 	speedKeypress = 10
// )

// typeCmd represents the type command
var typeCmd = &cobra.Command{
	Use:   "type",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		str := strings.Join(args, " ")
		fmt.Println(str)
		j := lib.Job{
			Name: "TYPE",
			Date: time.Now(),
		}
		// j := machine
		// q.Z += 50
		// q.Z += 100

		q.Z += 10
		// q.Z +=
		qd := q
		// qd.Z -= 16
		qd.Z -= 24

		for _, v := range str {
			key := arr[v]

			if key.shift {
				j.MoveTo(shift.toPos(q), speedTravel)
				j.MoveTo(shift.toPos(qd), speedKeypress)
				j.MoveTo(shift.toPos(q), speedTravel)
			}

			j.MoveTo(key.toPos(q), speedTravel)
			j.MoveTo(key.toPos(qd), speedKeypress)
			j.MoveTo(key.toPos(q), speedTravel)

		}

		j.MoveTo(arr['q'].toPos(q), speedTravel)

		buf := j.Buffer()

		if _, err := machine.RunCommand("DELETE TYPE"); err != nil {
			panic(err)
		}

		if err := machine.WriteFile("TYPE", &buf); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(typeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// typeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// typeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
