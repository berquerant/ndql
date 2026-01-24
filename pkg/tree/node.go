package tree

import (
	"fmt"
	"iter"
	"log/slog"
	"strings"

	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/berquerant/ndql/pkg/logx"
	"github.com/berquerant/ndql/pkg/node"
)

type (
	N                  = node.Node
	ND                 = node.Data
	OP                 = node.Op
	NIter              = iter.Seq[*N]
	NFunction          = iterx.Function[*N]
	NMapFunction       = iterx.MapFunction[*N]
	NReduceFunction    = iterx.ReduceFunction[*N]
	NFanoutFunction    = iterx.FanoutFunction[*N]
	NMultiMapFunction  = iterx.MultiMapFunction[*N]
	NDIter             = iter.Seq[ND]
	NDFunction         = iterx.Function[ND]
	NDMapFunction      = iterx.MapFunction[ND]
	NDReduceFunction   = iterx.ReduceFunction[ND]
	NDFanoutFunction   = iterx.FanoutFunction[ND]
	NDMultiMapFunction = iterx.MultiMapFunction[ND]
)

const (
	InputTableName    = "input"
	TableKeySeparator = "___"
	NodeValueKey      = "___value___"
)

type Key struct {
	Table  string
	Column string
}

func NewKey(table, column string) *Key {
	return &Key{
		Table:  table,
		Column: column,
	}
}

func KeyFromString(s string) *Key {
	xs := strings.SplitN(s, TableKeySeparator, 2)
	switch len(xs) {
	case 1: // column
		return NewKey("", s)
	case 2: // table___column
		return NewKey(xs[0], xs[1])
	default:
		panic(fmt.Errorf("%w: cannot create key from %s", ErrInvalidKey, s))
	}
}

func (k Key) String() string {
	if k.Table == "" {
		return k.Column
	}
	return k.Table + TableKeySeparator + k.Column
}

func KeyFromName(s string) *Key {
	xs := strings.SplitN(s, ".", 2)
	if len(xs) == 2 { // table.column
		return NewKey(xs[0], xs[1])
	}
	return NewKey("", s)
}

func (k Key) Name() string {
	if k.Table == "" {
		return k.Column
	}
	return k.Table + "." + k.Column
}

func (k Key) Get(n *N) (*N, bool) {
	r := node.New()
	switch {
	case k.Table != "": // like 'table.key'
		if v, ok := n.Get(k.String()); ok {
			r.Set(k.String(), v)
			return r, true
		}
	default: // like 'key'
		// fallback to other tables except default
		for _, s := range n.Keys() {
			t := KeyFromString(s)
			if t.Table == "" {
				// ignore default table
				continue
			}
			if t.Column == k.Column {
				v, _ := n.Get(s)
				r.Set(s, v)
				return r, true
			}
		}
		// fallback to default table
		if v, ok := n.Get(k.Column); ok {
			r.Set(k.Column, v)
			return r, true
		}
	}

	return nil, false
}

func (k Key) NFunction() NFunction {
	return iterx.NewMapFunction(func(n *N) (*N, error) {
		r, _ := k.Get(n)
		return r, nil
	})
}

type (
	NodeDataMapper = func(string, ND) (string, ND, error)
	NodeOpMapper   = func(string, *OP) (string, *OP, error)
)

func NodeDataMapperAsOP(f NodeDataMapper) NodeOpMapper {
	return func(k string, v *OP) (string, *OP, error) {
		k2, v2, err := f(k, v.AsData())
		if err != nil {
			return "", nil, err
		}
		return k2, v2.AsOp(), nil
	}
}

func NodeOpMapperAsData(f NodeOpMapper) NodeDataMapper {
	return func(k string, v ND) (string, ND, error) {
		k2, v2, err := f(k, v.AsOp())
		if err != nil {
			return "", nil, err
		}
		return k2, v2.AsData(), nil
	}
}

func MapNodeDataFunction(name string, f NodeDataMapper) func(*N) (*N, error) {
	return func(n *N) (*N, error) {
		if n == nil {
			return nil, ErrIgnore
		}
		x := node.New()
		for k, v := range n.Unwrap() {
			k2, v2, err := f(k, v)
			if err != nil {
				return nil, fmt.Errorf("%w: MapNodeData %s key=%s value=%v", err, name, k, v)
			}
			x.Set(k2, v2)
		}
		logx.Trace("MapNodeData", slog.String("name", name), logx.JSON("from", n), logx.JSON("to", x))
		return x, nil
	}
}

func AsUnaryArgUnaryRetNodeDataFunction(f func(*OP) (*OP, error)) func(ND) (ND, error) {
	return func(x ND) (ND, error) {
		r, err := f(x.AsOp())
		if err != nil {
			return nil, err
		}
		return r.AsData(), nil
	}
}

func AsVariadicArgUnaryRetNodeDataFunction(f func(...*OP) (*OP, error)) func(...ND) (ND, error) {
	return func(x ...ND) (ND, error) {
		ds := make([]*OP, len(x))
		for i, a := range x {
			ds[i] = a.AsOp()
		}
		r, err := f(ds...)
		if err != nil {
			return nil, err
		}
		return r.AsData(), nil
	}
}
