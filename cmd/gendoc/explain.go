package main

import (
	"errors"
	"fmt"
	"slices"

	"github.com/berquerant/ndql/pkg/gopkg"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(explainCmd)
	explainCmd.Flags().StringP("key", "k", "", "key")
}

var explainCmd = &cobra.Command{
	Use: "explain",
	RunE: func(cmd *cobra.Command, args []string) error {
		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			return errors.New("key is required")
		}
		s, err := loadDocumentSet(cmd, args)
		if err != nil {
			return err
		}
		v, ok := s.Get(key)
		if !ok {
			return errors.New("Not Found")
		}
		fmt.Println(v)
		return nil
	},
}

func loadDocumentSet(cmd *cobra.Command, args []string) (*gopkg.DocumentSet, error) {
	docs, err := loadDocuments(cmd, args)
	if err != nil {
		return nil, err
	}

	return gopkg.NewDocumentSet(slices.Collect(docs)...), nil
}
