package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows a fanfiction in the storage.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		swi := false

		lastInd := len(store.Fics) - 1

		for i, fic := range store.Fics {
			if swi {
				fmt.Print("\033[38:5:34m")
			} else {
				fmt.Print("\033[38:5:108m")
			}

			fmt.Println(fic.FormatSmall())

			if i == lastInd {
				fmt.Print("\033[0m")
			} else {
				fmt.Println()
			}

			swi = !swi
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
