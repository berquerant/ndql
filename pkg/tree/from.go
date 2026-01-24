package tree

import (
	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/berquerant/ndql/pkg/node"
	. "github.com/pingcap/tidb/pkg/parser/ast"
)

//
// FROM (SELECT ...) [AS ...]
//

func (v TreeVisitor) VisitTableRefsClause(n *TableRefsClause) (NFunction, error) {
	return v.VisitJoin(n.TableRefs)
}

func (v TreeVisitor) VisitJoin(n *Join) (NFunction, error) {
	switch x := n.Left.(type) {
	case *TableSource:
		return v.VisitTableSource(x)
	default:
		return nil, v.notImplemented(n, "unknown Join.Left")
	}
}

func (v TreeVisitor) VisitTableSource(n *TableSource) (NFunction, error) {
	switch x := n.Source.(type) {
	case *SelectStmt:
		f, err := v.VisitSelectStmt(x)
		if err != nil {
			return nil, v.newErr(err, n, "TableSource")
		}
		if s := n.AsName; s.O != "" {
			tableName := s.O
			return iterx.CombineFunction(f, iterx.NewMapFunction(MapNodeDataFunction(
				"ReplaceNodeTable",
				func(k string, v ND) (string, ND, error) {
					key := KeyFromString(k)
					if key.Table == "" && node.IsBuiltinKey(key.Column) {
						return k, v, nil
					}
					key.Table = tableName
					return key.String(), v, nil
				},
			)))
		}
		return f, nil
	default:
		return nil, v.notImplemented(n, "unknown Source")
	}
}
