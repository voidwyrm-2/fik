package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/voidwyrm-2/fik/internal/fic"
)

var (
	addFlag_file                 *string
	addFlag_first, addFlag_force *bool
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a fanfiction to the storage.",
	Long: `Adds a fanfiction to the storage.
Use the format '[id],[chapter id]' to additional give a fiction a current chapter.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var entries []string

		if len(*addFlag_file) > 0 {
			data, err := os.ReadFile(*addFlag_file)
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

			id, chapterId, err := fic.ParseFicEntry(entry)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				continue
			}

			if idx, ok := store.Ids[id]; ok && !*addFlag_force {
				fmt.Fprintf(os.Stderr, "Fiction `%s` is already in the storage\n", store.Fics[idx].Title)
				continue
			}

			fic, err := fic.GetFicFromId(id, chapterId, *addFlag_first)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				continue
			}

			store.Ids[id] = uint32(len(store.Fics))
			store.Fics = append(store.Fics, fic)

			if fic.ChapterInfo.Id != 0 {
				fmt.Printf("Added the fiction `%s` (with current chapter %d, `%s`, id %d) to the storage\n", fic.Title, fic.ChapterInfo.Num, fic.ChapterInfo.Title, fic.ChapterInfo.Id)
			} else {
				fmt.Printf("Added the fiction `%s` to the storage\n", fic.Title)
			}
		}

		return nil
	},
}

func init() {
	addFlag_file = addCmd.Flags().StringP("file", "f", "", "Read the fictions from a file.")
	addFlag_first = addCmd.Flags().Bool("first", false, "Set the current chapter of a fiction to its first chapter if a current chapter isn't specified.")
	addFlag_force = addCmd.Flags().Bool("force", false, "Force the addition of the fictions, even if they're already in storage")
	rootCmd.AddCommand(addCmd)
}
