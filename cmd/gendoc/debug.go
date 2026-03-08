package main

import (
	"encoding/json"
	"fmt"
	"iter"
	"log/slog"

	"github.com/berquerant/ndql/pkg/gopkg"
	"github.com/berquerant/ndql/pkg/logx"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.Flags().Bool("doc", true, "extract documents")
	debugCmd.Flags().Bool("docstring", false, "extract document as a string")
}

var debugCmd = &cobra.Command{
	Use: "debug",
	RunE: func(cmd *cobra.Command, args []string) error {
		if doc, _ := cmd.Flags().GetBool("doc"); doc {
			xs, err := loadDocuments(cmd, args)
			if err != nil {
				return err
			}
			str, _ := cmd.Flags().GetBool("docstring")
			for x := range xs {
				if str {
					fmt.Print(x.String())
					continue
				}
				b, _ := json.Marshal(x)
				fmt.Printf("%s\n", b)
			}
			return nil
		}

		xs, err := loadComments(cmd, args)
		if err != nil {
			return err
		}
		for x := range xs {
			b, _ := json.Marshal(x)
			fmt.Printf("%s\n", b)
		}
		return nil
	},
}

func loadComments(cmd *cobra.Command, args []string) (iter.Seq[*gopkg.Comment], error) {
	if len(args) == 0 {
		args = []string{"./..."}
	}
	loader := gopkg.NewLoader()
	if err := loader.Load(cmd.Context(), args...); err != nil {
		return nil, err
	}
	return loader.Comments(), nil
}

func loadDocuments(cmd *cobra.Command, args []string) (iter.Seq[*gopkg.Document], error) {
	comments, err := loadComments(cmd, args)
	if err != nil {
		return nil, err
	}

	return func(yield func(*gopkg.Document) bool) {
		for x := range comments {
			d, err := x.GetDocument()
			if err != nil {
				slog.Debug("failed to get document", logx.Err(err))
				continue
			}
			if !yield(d) {
				return
			}
		}
	}, nil
}
