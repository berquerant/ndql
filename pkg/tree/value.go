package tree

import (
	"fmt"
	"time"

	"github.com/berquerant/ndql/pkg/errorx"
	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/berquerant/ndql/pkg/node"
	. "github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pingcap/tidb/pkg/parser/mysql"
	"github.com/pingcap/tidb/pkg/types"
	driver "github.com/pingcap/tidb/pkg/types/parser_driver"
)

func (v TreeVisitor) VisitValueExpr(n ValueExpr) (NFunction, error) {
	switch n := n.(type) {
	case *driver.ValueExpr:
		x, err := v.VisitValueExprDriver(n)
		if err != nil {
			return nil, v.newErr(err, n, "ValueExprDriver")
		}
		return ReturnContainerValue("ValueExpr", func(_ ND) (ND, error) {
			return x, nil
		}), nil
	default:
		return nil, v.notImplemented(n, "unknown ValueExpr")
	}
}

func (v TreeVisitor) VisitValueExprDriver(n *driver.ValueExpr) (ND, error) {
	switch n.Kind() {
	case types.KindNull:
		return node.NewNull(), nil
	case types.KindInt64:
		if n.Type.GetFlag()&mysql.IsBooleanFlag != 0 {
			return node.Bool(n.GetInt64() > 0), nil
		}
		return node.Int(n.GetInt64()), nil
	case types.KindUint64:
		return node.Int(int64(n.GetUint64())), nil
	case types.KindFloat32, types.KindFloat64:
		return node.Float(n.GetFloat64()), nil
	case types.KindString:
		return node.String(n.GetString()), nil
	case types.KindMysqlDecimal:
		f, err := n.GetMysqlDecimal().ToFloat64()
		if err != nil {
			return nil, v.newErr(err, n, "convert mysql decimal to float64")
		}
		return node.Float(f), nil
	case types.KindMysqlDuration:
		return node.Duration(n.GetMysqlDuration().Duration), nil
	case types.KindMysqlTime:
		t, err := n.GetMysqlTime().GoTime(time.UTC)
		if err != nil {
			return nil, v.newErr(err, n, "convert mysql time to go time")
		}
		return node.Time(t), nil
	default:
		return nil, v.notImplemented(n, "unknown ValueExpr driver(kind=%d)", n.Kind())
	}
}

// A single value as a Node.
type ValueContainer struct {
	*N
}

func AsValueContainer(n *N) *ValueContainer {
	if n == nil {
		n = node.New()
	}
	return &ValueContainer{n}
}

func (c *ValueContainer) SetContainerValue(v ND)        { c.Set(NodeValueKey, v) }
func (c *ValueContainer) GetContainerValue() (ND, bool) { return c.Get(NodeValueKey) }

// GetFirstValue unwraps the value.
// Returns the key, value and exists value or not.
func (c *ValueContainer) GetFirstValue() (string, ND, bool) {
	keys := c.Keys()
	if len(keys) == 0 {
		return "", nil, false
	}
	v, ok := c.Get(keys[0])
	return keys[0], v, ok
}

// Create a function that converts a ValueContainer.
func ReturnContainerValue(name string, f func(ND) (ND, error)) NFunction {
	return iterx.NewMapFunction(func(n *N) (*N, error) {
		key, value, ok := AsValueContainer(n).GetFirstValue()
		if !ok {
			return nil, fmt.Errorf("%w: %s no upstream value", errorx.WithVerbose(ErrInvalidValue, n), name)
		}
		got, err := f(value)
		if err != nil {
			return nil, fmt.Errorf("%w: %s failed, key=%s", errorx.WithVerbose(err, n), name, key)
		}
		r := AsValueContainer(node.New())
		r.SetContainerValue(got)
		return r.N, nil
	})
}

// Extract a single value from Expr.
func (v TreeVisitor) visitValueExpr(n ExprNode) (ND, error) {
	switch n := n.(type) {
	case ValueExpr:
		f, err := v.VisitValueExpr(n)
		if err != nil {
			return nil, err
		}
		dummy := AsValueContainer(node.New())
		dummy.SetContainerValue(node.String("dummy"))
		d, err := f.CallAny(dummy.N)
		if err != nil {
			return nil, v.newErr(err, n, "cannot call value expr function")
		}
		if len(d) == 0 {
			return nil, v.invalidTree(n, "cannot get value from value expr function")
		}
		e, ok := AsValueContainer(d[0]).GetContainerValue()
		if !ok {
			return nil, v.invalidValue(n, "cannot get container value from value expr")
		}
		return e, nil
	default:
		return nil, v.notImplemented(n, "want ValueExpr")
	}
}
