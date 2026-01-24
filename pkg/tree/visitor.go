package tree

import (
	"context"

	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/berquerant/ndql/pkg/node"
	. "github.com/pingcap/tidb/pkg/parser/ast"
)

type TreeVisitor struct {
	ctx context.Context
}

func NewTreeVisitor(ctx context.Context) *TreeVisitor {
	return &TreeVisitor{
		ctx: ctx,
	}
}

func (v TreeVisitor) Visit(n Node) (NFunction, error) {
	switch n := n.(type) {
	case *SelectStmt:
		return v.VisitSelectStmt(n)
	default:
		return nil, v.notImplemented(n, "Visit")
	}
}

func (v TreeVisitor) VisitSelectStmt(n *SelectStmt) (NFunction, error) {
	var (
		fs = []NFunction{}
	)
	if x := n.From; x != nil {
		f, err := v.VisitTableRefsClause(x)
		if err != nil {
			return nil, v.newErr(err, n, "From")
		}
		fs = append(fs, f)
	}
	if x := n.Where; x != nil {
		f, err := v.VisitWhere(x)
		if err != nil {
			return nil, v.newErr(err, n, "Where")
		}
		fs = append(fs, f)
	}
	if x := n.Fields; x != nil {
		f, err := v.VisitFieldList(x)
		if err != nil {
			return nil, v.newErr(err, n, "FieldList")
		}
		fs = append(fs, f)
	}
	if len(fs) == 0 {
		return nil, v.notImplemented(n, "unknown SelectStmt")
	}
	r := fs[0]
	for _, f := range fs[1:] {
		nf, err := iterx.CombineFunction(r, f)
		if err != nil {
			return nil, v.newErr(err, n, "SelectStmt failed to combine function")
		}
		r = nf
	}
	return r, nil
}

func (v TreeVisitor) VisitFieldList(n *FieldList) (NFunction, error) {
	var (
		fs = make([]NFunction, len(n.Fields))
	)
	for i, f := range n.Fields {
		x, err := v.VisitSelectField(f)
		if err != nil {
			return nil, v.newErr(err, n, "FieldList[%d]", i)
		}
		fs[i] = x
	}
	if err := ValidateOnlyVariadicOrAllUnaryRet(fs...); err != nil {
		return nil, v.newErr(err, n, "FieldList should contain only variadic or all unary ret")
	}

	switch {
	case fs[0].RetArity() == iterx.Variadic:
		return fs[0], nil
	default:
		return iterx.NewMapFunction(func(x *N) (*N, error) {
			r := node.New()
			for i, f := range fs {
				a, err := f.CallAny(x)
				if err != nil {
					return nil, v.newErr(err, n, "function call[%d]", i)
				}
				for _, b := range a {
					r.Merge(b.Map)
				}
			}
			return r, nil
		}), nil
	}
}

func (v TreeVisitor) VisitSelectField(n *SelectField) (NFunction, error) {
	if n.WildCard != nil {
		return iterx.NewMapFunction(iterx.Identity[*N]), nil
	}
	f, err := v.VisitExpr(n.Expr)
	if err != nil {
		return nil, v.newErr(err, n, "SelectField Expr")
	}
	if x := n.AsName; x.O != "" {
		column := x.O
		return iterx.CombineFunction(f, iterx.NewMapFunction(MapNodeDataFunction(
			"ReplaceNodeColumn",
			func(k string, v ND) (string, ND, error) {
				key := KeyFromString(k)
				key.Column = column
				return key.String(), v, nil
			},
		)))
	}
	return iterx.CombineFunction(f, iterx.NewMapFunction(MapNodeDataFunction(
		"ReplaceDefaultNodeColumn",
		func(k string, v ND) (string, ND, error) {
			if k == NodeValueKey {
				return n.Text(), v, nil
			}
			return k, v, nil
		},
	)))
}
