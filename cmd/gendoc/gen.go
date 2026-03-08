package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.AddCommand(genJSONCmd)
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
