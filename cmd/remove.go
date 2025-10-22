package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/voidwyrm-2/fik/internal/fic"
)

var removeFlag_file *string

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes a fanfiction from the storage.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var entries []string

		if len(*removeFlag_file) > 0 {
			data, err := os.ReadFile(*removeFlag_file)
			if err != nil {
				return err
			}

			entries = strings.Fields(string(data))
		} else {
			entries = args
		}

		for _, entry := range entries {
			if len(entry) == 0 {
				continue
			}

			id, _, err := fic.ParseFicEntry(entry)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				continue
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
				fmt.Fprintf(os.Stderr, "No fiction with ID %d in storage\n", id)
			}
		}

		return nil
	},
}

func init() {
	removeFlag_file = removeCmd.Flags().StringP("file", "f", "", "Read the fictions from a file.")
	rootCmd.AddCommand(removeCmd)
}
