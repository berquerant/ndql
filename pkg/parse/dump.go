package parse

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/berquerant/ndql/pkg/logx"
	ast "github.com/pingcap/tidb/pkg/parser/ast"
)

type DumpMode int

const (
	DumpModeUnknown DumpMode = iota
	DumpModeText
	DumpModeVerbose
)

// Dump displays AST as text.
func Dump(w io.Writer, n ast.Node, indentString string, mode DumpMode) {
	n.Accept(NewTreeDumper(w, indentString, mode))
}

type TreeDumper struct {
	w            io.Writer
	nodeIndex    int
	indentLevel  int
	indentString string
	mode         DumpMode
}

func NewTreeDumper(w io.Writer, indentString string, mode DumpMode) *TreeDumper {
	return &TreeDumper{
		w:            w,
		indentString: indentString,
		mode:         mode,
	}
}

func (d *TreeDumper) print(n ast.Node) {
	indent := strings.Repeat(d.indentString, d.indentLevel)
	output := fmt.Sprintf("%05d:%03d:%s:%T:",
		d.nodeIndex, d.indentLevel, indent, n,
	)
	switch d.mode {
	case DumpModeText:
		output += fmt.Sprintf("%v", n)
	case DumpModeVerbose:
		b, err := json.Marshal(n)
		if err != nil {
			slog.Error("TreeDumper",
				slog.Int("index", d.nodeIndex),
				slog.Int("indent", d.indentLevel),
				logx.Err(err),
			)
			return
		}
		output += string(b)
	}
	_, _ = fmt.Fprintln(d.w, output)
}

func (d TreeDumper) trace(n ast.Node, msg string) {
	logx.Trace(msg, slog.Int("index", d.nodeIndex), slog.Int("indent", d.indentLevel))
}

func (d *TreeDumper) Enter(n ast.Node) (ast.Node, bool) {
	d.trace(n, "Enter")
	d.print(n)
	d.nodeIndex++
	d.indentLevel++
	return n, false
}

func (d *TreeDumper) Leave(n ast.Node) (ast.Node, bool) {
	d.trace(n, "Leave")
	d.indentLevel--
	return n, true
}
