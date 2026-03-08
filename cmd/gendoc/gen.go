package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.AddCommand(genJSONCmd, genFilesCmd)
	genFilesCmd.Flags().String("root", "", "root directory")
}

var genCmd = &cobra.Command{
	Use: "gen",
}

var genJSONCmd = &cobra.Command{
	Use: "json",
	RunE: func(cmd *cobra.Command, args []string) error {
		docs, err := loadDocumentSet(cmd, args)
		if err != nil {
			return err
		}
		b, err := json.Marshal(docs)
		if err != nil {
			return err
		}
		_, err = fmt.Printf("%s\n", b)
		return err
	},
}

var genFilesCmd = &cobra.Command{
	Use: "files",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, _ := cmd.Flags().GetString("root")
		if root == "" {
			return errors.New("root is required")
		}
		docs, err := loadDocumentSet(cmd, args)
		if err != nil {
			return err
		}
		return docs.IntoFiles(root)
	},
}
