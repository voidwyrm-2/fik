package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/voidwyrm-2/fik/internal/fic"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a fanfiction to the storage.",
	Run: func(cmd *cobra.Command, args []string) {
		for _, a := range args {
			id, err := fic.ParseId(a)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				continue
			}

			if idx, ok := store.Ids[id]; ok {
				fmt.Printf("Fiction `%s` is already in the storage\n", store.Fics[idx].Title)
				continue
			}

			fic, err := fic.GetFicFromId(id)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				continue
			}

			store.Ids[id] = uint32(len(store.Fics))
			store.Fics = append(store.Fics, fic)

			fmt.Printf("Added the fiction `%s` to the storage\n", fic.Title)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
