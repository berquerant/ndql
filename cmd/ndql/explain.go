package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/berquerant/ndql/pkg/gopkg"
	"github.com/spf13/cobra"
)

//go:embed docs.json
var docs []byte

var documentSet gopkg.DocumentSet

func init() {
	if err := json.Unmarshal(docs, &documentSet); err != nil {
		panic(fmt.Errorf("%w: failed to load docs.json", err))
	}
	rootCmd.AddCommand(explainCmd)
}

var explainCmd = &cobra.Command{
	Use:   "explain KEY",
	Short: "Describe ndql resources.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println(explainShowKeys())
			return nil
		}
		d, ok := documentSet.Get(args[0])
		if !ok {
			fmt.Println(explainShowKeys())
			return fmt.Errorf("%s is not found", args[0])
		}
		fmt.Println(d)
		return nil
	},
}

func explainShowKeys() string {
	keys := documentSet.Keys()
	xs := make([]string, len(keys))
	for i, k := range keys {
		xs[i] = "- " + k
	}
	return fmt.Sprintf("Available keys:\n%s", strings.Join(xs, "\n"))
}
