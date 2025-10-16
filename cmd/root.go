package cmd

import (
	_ "embed"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/voidwyrm-2/fik/internal/fic"
)

var version string

type Store struct {
	Ids  map[fic.Id]uint32
	Fics []fic.Fic
}

var store = Store{
	Ids:  map[fic.Id]uint32{},
	Fics: []fic.Fic{},
}

var rootCmd = &cobra.Command{
	Use:   "fik",
	Short: "A tool for keeping a local catalog of AO3 fanfiction.",
}

func Execute(ver string) error {
	version = ver

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	storePath := filepath.Join(home, ".fikstore.json")

	if _, err = os.Stat(storePath); err == nil {
		err = func() error {
			f, err := os.Open(storePath)
			if err != nil {
				return err
			}

			defer f.Close()

			d := json.NewDecoder(f)

			return d.Decode(&store)
		}()
		if err != nil {
			return err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	defer func() {
		f, err := os.Create(storePath)
		if err != nil {
			return
		}

		defer f.Close()

		d := json.NewEncoder(f)

		d.Encode(&store)
	}()

	if err = rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	return nil
}
