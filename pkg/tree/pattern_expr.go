package tree

import (
	"fmt"

	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/berquerant/ndql/pkg/node"
	. "github.com/pingcap/tidb/pkg/parser/ast"
)

// REGEXP
func (v TreeVisitor) VisitPatternRegexpExpr(n *PatternRegexpExpr) (NFunction, error) {
	f, err := v.VisitExpr(n.Expr)
	if err != nil {
		return nil, v.newErr(err, n, "PatternRegexpExpr expr")
	}
	g, err := v.visitValueExpr(n.Pattern)
	if err != nil {
		return nil, v.newErr(err, n, "PatternRegexpExpr pattern")
	}
	return iterx.CombineFunction(f, ReturnContainerValue("PatternRegexpExpr", func(x ND) (ND, error) {
		r, err := x.AsOp().Regexp(g.AsOp())
		if err != nil {
			return nil, fmt.Errorf("%w: eval", err)
		}
		if n.Not {
			if r, err = r.Not(); err != nil {
				return nil, fmt.Errorf("%w: not", err)
			}
		}
		return r.AsData(), nil
	}))
}

// LIKE
func (v TreeVisitor) VisitParrernLikeExpr(n *PatternLikeOrIlikeExpr) (NFunction, error) {
	f, err := v.VisitExpr(n.Expr)
	if err != nil {
		return nil, v.newErr(err, n, "PatternLike expr")
	}
	g, err := v.visitValueExpr(n.Pattern)
	if err != nil {
		return nil, v.newErr(err, n, "PatternLike pattern")
	}
	return iterx.CombineFunction(f, ReturnContainerValue("PatternLikeExpr", func(x ND) (ND, error) {
		r, err := x.AsOp().Like(g.AsOp())
		if err != nil {
			return nil, fmt.Errorf("%w: eval", err)
		}
		if n.Not {
			if r, err = r.Not(); err != nil {
				return nil, fmt.Errorf("%w: not", err)
			}
		}
		return r.AsData(), nil
	}))
}

// BETWEEN ... AND ...
func (v TreeVisitor) VisitPatternInExpr(n *PatternInExpr) (NFunction, error) {
	f, err := v.VisitExpr(n.Expr)
	if err != nil {
		return nil, v.newErr(err, n, "PatternInExpr expr")
	}
	list := make([]*OP, len(n.List))
	for i, x := range n.List {
		d, err := v.visitValueExpr(x)
		if err != nil {
			return nil, v.newErr(err, n, "PatternInExpr list[%d]", i)
		}
		list[i] = d.AsOp()
	}
	return iterx.CombineFunction(f, ReturnContainerValue("PatternInExpr", func(x ND) (ND, error) {
		return node.Bool(x.AsOp().In(list...) != n.Not), nil
	}))
}

func (v TreeVisitor) VisitBetweenExpr(n *BetweenExpr) (NFunction, error) {
	f, err := v.VisitExpr(n.Expr)
	if err != nil {
		return nil, v.newErr(err, n, "BetweenExpr expr")
	}
	left, err := v.visitValueExpr(n.Left)
	if err != nil {
		return nil, v.newErr(err, n, "BetweenExpr left")
	}
	right, err := v.visitValueExpr(n.Right)
	if err != nil {
		return nil, v.newErr(err, n, "BetweenExpr right")
	}
	return iterx.CombineFunction(f, ReturnContainerValue("BetweenExpr", func(x ND) (ND, error) {
		b, err := x.AsOp().Between(left.AsOp(), right.AsOp())
		if err != nil {
			return nil, err
		}
		return node.Bool(b != n.Not), nil
	}))
}
