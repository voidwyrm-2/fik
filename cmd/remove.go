package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voidwyrm-2/fik/internal/fic"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes a fanfiction from the storage.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := fic.ParseId(args[0])
		if err != nil {
			return err
		}

		if idx, ok := store.Ids[id]; ok {
			delete(store.Ids, id)

			oldFic := store.Fics[idx]

			newFics := make([]fic.Fic, 0, len(store.Fics)-1)

			newFics = append(newFics, store.Fics[:idx]...)

			newFics = append(newFics, store.Fics[idx+1:]...)

			store.Fics = newFics

			fmt.Printf("Fiction `%s` (ID %d) from storage\n", oldFic.Title, id)
		} else {
			return fmt.Errorf("No fiction with ID %d in storage", id)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
