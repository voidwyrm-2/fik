package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voidwyrm-2/fik/internal/fic"
)

var favCmd = &cobra.Command{
	Use:   "fav",
	Short: "Favorite a fanfiction in the storage.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := fic.ParseId(args[0])
		if err != nil {
			return err
		}

		idx, ok := store.Ids[id]
		if !ok {
			return fmt.Errorf("No fiction with ID %d in storage", id)
		}

		fic := &store.Fics[idx]

		fic.Favorite = !fic.Favorite

		if fic.Favorite {
			fmt.Printf("Favorited `%s`\n", fic.Title)
		} else {
			fmt.Printf("Unfavorited `%s`\n", fic.Title)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(favCmd)
}
