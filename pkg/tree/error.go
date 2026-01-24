package tree

import (
	"errors"
	"fmt"
	"slices"

	"github.com/berquerant/ndql/pkg/errorx"
	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/pingcap/tidb/pkg/parser/ast"
)

var (
	ErrIgnore               = errors.New("Ignore")
	ErrNotImplmented        = errors.New("NotImplemented")
	ErrInvalidTree          = errors.New("InvalidTree")
	ErrInvalidKey           = errors.New("InvalidKey")
	ErrInvalidValue         = errors.New("InvalidValue")
	ErrInvalidArgument      = errors.New("InvalidArgument")
	ErrInvalidFunctionArity = errors.New("InvalidFunctionArity")
	ErrParseGenResult       = errors.New("ParseGenResult")
	ErrGenTemplate          = errors.New("GenTemplate")
)

func (TreeVisitor) newErr(err error, n ast.Node, msg string, v ...any) error {
	return fmt.Errorf("%w: %T, %s",
		errorx.WithVerbose(err, n),
		n,
		fmt.Sprintf(msg, v...),
	)
}

func (v TreeVisitor) notImplemented(n ast.Node, msg string, a ...any) error {
	return v.newErr(ErrNotImplmented, n, msg, a...)
}

func (v TreeVisitor) invalidTree(n ast.Node, msg string, a ...any) error {
	return v.newErr(ErrInvalidTree, n, msg, a...)
}

func (v TreeVisitor) invalidValue(n ast.Node, msg string, a ...any) error {
	return v.newErr(ErrInvalidValue, n, msg, a...)
}

func ValidateOnlyVariadicOrAllUnaryRet(fs ...NFunction) error {
	switch {
	case len(fs) == 0:
		return fmt.Errorf("%w: no functions", ErrInvalidFunctionArity)
	case slices.IndexFunc(fs, func(x NFunction) bool {
		return x.RetArity() == iterx.Variadic
	}) >= 0:
		if len(fs) != 1 {
			return fmt.Errorf("%w: multiple variadic ret", ErrInvalidFunctionArity)
		}
		return nil
	}

	return nil
}
