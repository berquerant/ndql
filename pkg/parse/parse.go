package parse

import (
	"errors"

	tidb "github.com/pingcap/tidb/pkg/parser"
	ast "github.com/pingcap/tidb/pkg/parser/ast"
)

type Parser interface {
	Parse(s string) (*Result, error)
}

type Node = ast.StmtNode

type Result struct {
	Nodes []Node
	Warns []error
}

const (
	charset   = "utf8mb4"
	collation = "utf8mb4_0900_ai_ci"
)

type SQLParser struct {
	parser *tidb.Parser
}

func NewSQLParser() *SQLParser {
	return &SQLParser{
		parser: tidb.New(),
	}
}

var _ Parser = &SQLParser{}

var ErrParse = errors.New("ParseError")

func (p *SQLParser) Parse(s string) (*Result, error) {
	nodes, warns, err := p.parser.Parse(s, charset, collation)
	if err != nil {
		return nil, errors.Join(ErrParse, err)
	}
	return &Result{
		Nodes: nodes,
		Warns: warns,
	}, nil
}
