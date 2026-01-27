package main

import (
	"fmt"

	"github.com/berquerant/ndql/pkg/config"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query QUERY [PATH]",
	Short: "Run query",
	Long: fmt.Sprintf(`Run query.

## QUERY, PATH
%s

## Examples

List all files and directories under the dir:

    ndql query 'select *' dir

List all file paths under the dir:

    ndql query 'select path where not is_dir' dir

Add a "count" column for the number of times 'keyword' appears in the file:

    ndql query 'select sh("echo count=$(grep -c keyword $(get path))") where not is_dir'

Roughly list the Go func in the file except tests:

    ndql query 'select path, func from (select grep("func (?P<name>[^(]+)", "func=$name") where not is_dir and path not like "%%_test.go")' dir

Add the line number where each func is defined as 'line':

    ndql query 'select path, func from (select grep("func (?P<name>[^(]+)", "func=$name") where not is_dir and path not like "%%_test.go")' dir | ndql query 'select sh("echo line=$(grep -n \"func $(get func)\" $(get path) | cut -d: -f1)")' -i@-

Roughly lists the names and values of variables exported in bash format. The search is restricted only to files tracked by git:

    git ls-files | ndql query 'select grep("export (?P<name>[^=]+)=(?P<value>.+)", "name=$name,value=$value") where not is_dir' @-

Extracts metadata from mp3 and m4a files using ffprobe:

    ndql query 'select sh("ffprobe -v error -hide_banner -show_entries format -of json=c=1 \"$(get path)\" | jq .format.tags -c") where not is_dir and extension(path) in (".mp3", ".m4a")' dir
`, config.DescribeSourceUsage()),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMain(cmd, args, config.ModeQuery)
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
	initFlags(queryCmd)
}
