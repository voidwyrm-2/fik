package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voidwyrm-2/fik/internal/filters"
)

var (
	listFlag_filters *[]string
	listFlag_only    *uint
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows a fanfiction in the storage.",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		swi := false

		lastInd := len(store.Fics) - 1

		filtered, filtersUsed, err := filters.Filter(store.Fics, *listFlag_filters)
		if err != nil {
			return err
		}

		for _, f := range filtersUsed {
			fmt.Printf("Using filter %s\n", f)
		}

		if len(filtersUsed) > 0 {
			fmt.Println()
		}

		for i, fic := range filtered {

			if swi {
				fmt.Print("\033[38:5:34m")
			} else {
				fmt.Print("\033[38:5:108m")
			}

			if fic.Favorite {
				fmt.Print("⭐️ ")
			}

			fmt.Println(fic.FormatSmall())

			if (i != lastInd && i+1 != int(*listFlag_only)) || len(filtersUsed) > 0 {
				fmt.Println()
			}

			if *listFlag_only > 0 && i+1 == int(*listFlag_only) {
				break
			}

			swi = !swi
		}

		fmt.Print("\033[0m")

		return nil
	},
}

func init() {
	listFlag_filters = listCmd.Flags().StringSliceP("filter", "f", []string{}, `The filter(s) to use when listing fictions.
Valid values are:
  favorites
  rating:<unrated|general|teen|mature|explicit>
  author:<name>`)
	listFlag_only = listCmd.Flags().UintP("only", "o", 0, "The amount of fictions to show.\nZero means it shows all fictions.")
	rootCmd.AddCommand(listCmd)
}
