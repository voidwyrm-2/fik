package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voidwyrm-2/fik/internal/fic"
)

var addchapCmd = &cobra.Command{
	Use:   "setchap",
	Short: "Sets a current chapter on a specified fanfiction in the storage.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := fic.ParseId(args[0])
		if err != nil {
			return err
		}

		chapterId, err := fic.ParseId(args[1])
		if err != nil {
			return err
		} else if chapterId == 0 {
			return fmt.Errorf("'%d' is not a valid AO3 fiction ID", chapterId)
		}

		idx, ok := store.Ids[id]
		if !ok {
			return fmt.Errorf("No fiction with ID %d in storage", id)
		}

		fic := &store.Fics[idx]

		fic.ChapterInfo.Id = chapterId

		err = fic.GetCurrentChapterInfo()
		if err != nil {
			return err
		}

		if fic.ChapterInfo.Id != 0 {
			fmt.Printf("The current chapter of the fiction `%s` was set to chapter %d, `%s`, id %d in storage\n", fic.Title, fic.ChapterInfo.Num, fic.ChapterInfo.Title, fic.ChapterInfo.Id)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addchapCmd)
}
