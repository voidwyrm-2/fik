package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ficsCmd = &cobra.Command{
	Use:   "fics",
	Short: "Shows the amount of fictions in storage.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("There are %d fictions in storage\n", len(store.Ids))
	},
}

func init() {
	rootCmd.AddCommand(ficsCmd)
}
