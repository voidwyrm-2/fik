package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/voidwyrm-2/fik/internal/fic"
)

var (
	cleanFlag_title, cleanFlag_author, cleanFlag_summary, cleanFlag_fandoms *bool
	cleanFlag_relationships, cleanFlag_characters, cleanFlag_tags           *bool
	cleanFlag_all                                                           *bool
)

func clean(s *string) {
	*s = strings.TrimSpace(*s)
}

func cleanItems(items []string) {
	for i := range items {
		clean(&items[i])
	}
}

func cleanFic(f *fic.Fic) {
	if *cleanFlag_title {
		clean(&f.Title)
	}

	if *cleanFlag_author {
		clean(&f.Author)
	}

	if *cleanFlag_summary {
		clean(&f.Summary)
	}

	if *cleanFlag_fandoms {
		cleanItems(f.Fandoms)
	}

	if *cleanFlag_relationships {
		cleanItems(f.Relationships)
	}

	if *cleanFlag_characters {
		cleanItems(f.Characters)
	}

	if *cleanFlag_tags {
		cleanItems(f.Tags)
	}
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans the textual parts of fictions in the storage, e.g. the title or summary.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if *cleanFlag_all {
			for i := range store.Fics {
				f := &store.Fics[i]

				cleanFic(f)

				fmt.Printf("Cleaned `%s` (ID %d)\n", f.Title, f.Id)
			}
		} else {
			for _, a := range args {
				id, err := fic.ParseId(a)
				if err != nil {
					return err
				}

				idx, ok := store.Ids[id]
				if !ok {
					return fmt.Errorf("No fiction with ID %d in storage", id)
				}

				f := &store.Fics[idx]

				cleanFic(f)

				fmt.Printf("Cleaned `%s` (ID %d)\n", f.Title, f.Id)
			}
		}

		return nil
	},
}

func init() {
	cleanFlag_title = cleanCmd.Flags().Bool("title", false, "Clean the fiction title.")
	cleanFlag_author = cleanCmd.Flags().Bool("author", false, "Clean the fiction author.")
	cleanFlag_summary = cleanCmd.Flags().Bool("summary", false, "Clean the fiction summary.")
	cleanFlag_fandoms = cleanCmd.Flags().Bool("fandoms", false, "Clean the fiction fandoms.")
	cleanFlag_relationships = cleanCmd.Flags().BoolP("relationships", "r", false, "Clean the fiction relationships.")
	cleanFlag_characters = cleanCmd.Flags().BoolP("characters", "c", false, "Clean the fiction characters.")
	cleanFlag_tags = cleanCmd.Flags().Bool("tags", false, "Clean the fiction tags.")

	cleanFlag_all = cleanCmd.Flags().Bool("all", false, "Clean all stored fictions.")

	rootCmd.AddCommand(cleanCmd)
}
