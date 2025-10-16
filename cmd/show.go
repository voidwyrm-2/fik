package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voidwyrm-2/fik/internal/fic"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows a fanfiction in the storage.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := fic.ParseId(args[0])
		if err != nil {
			return err
		}

		if idx, ok := store.Ids[id]; ok {
			fmt.Println(&store.Fics[idx])
		} else {
			return fmt.Errorf("No fiction with ID %d in storage", id)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
