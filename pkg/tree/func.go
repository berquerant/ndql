package tree

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/berquerant/ndql/pkg/iox"
	"github.com/berquerant/ndql/pkg/iterx"
	"github.com/berquerant/ndql/pkg/logx"
	"github.com/berquerant/ndql/pkg/node"
	. "github.com/pingcap/tidb/pkg/parser/ast"
)

// Function calls except aggregations.
func (v TreeVisitor) VisitFuncCallExpr(n *FuncCallExpr) (NFunction, error) {
	name := n.FnName.L
	f, err := v.visitFuncCallExpr(n)
	if err != nil {
		return nil, v.newErr(err, n, "FuncCallExpr[%s]", name)
	}
	return f, nil
}

const (
	FuncArgMaxLen = 100
)

const (
	FuncGrep = "grep"
	FuncTmpl = "tmpl"
	FuncSh   = "sh"
	FuncLua  = "lua"
	FuncExpr = "expr"

	FuncToInt      = "to_int"
	FuncToFloat    = "to_float"
	FuncToBool     = "to_bool"
	FuncToString   = "to_string"
	FuncToTime     = "to_time"
	FuncToDuration = "to_duration"

	FuncLeast    = "least"
	FuncGreatest = "greatest"
	FuncCoalesce = "coalesce"

	FuncIf     = "if"
	FuncIfNull = "ifnull"
	FuncNullIf = "nullif"

	FuncAbs     = "abs"
	FuncSqrt    = "sqrt"
	FuncDegrees = "degrees"
	FuncRadians = "radians"
	FuncAcos    = "acos"
	FuncAsin    = "asin"
	FuncAtan    = "atan"
	FuncCos     = "cos"
	FuncSin     = "sin"
	FuncTan     = "tan"
	FuncCot     = "cot"
	FuncLn      = "ln"
	FuncLog2    = "log2"
	FuncLog10   = "log10"
	FuncExp     = "exp"
	FuncCeil    = "ceil"
	FuncFloor   = "floor"
	FuncRound   = "round"
	FuncAtan2   = "atan2"
	FuncPow     = "pow"
	FuncE       = "e"
	FuncPi      = "pi"
	FuncRand    = "rand"

	FuncLen           = "len"
	FuncSize          = "size"
	FuncRegexpCount   = "regexp_count"
	FuncRegexpInstr   = "regexp_instr"
	FuncRegexpSubstr  = "regexp_substr"
	FuncRegexpReplace = "regexp_replace"
	FuncRegexpLike    = "regexp_like"
	FuncFormat        = "format"
	FuncLower         = "lower"
	FuncUpper         = "upper"
	FuncSha2          = "sha2"
	FuncConcatWs      = "concat_ws"
	FuncInstr         = "instr"
	FuncInstrCount    = "instr_count"
	FuncSubstr        = "substr"
	FuncSubstrIndex   = "substr_index"
	FuncReplace       = "replace"
	FuncTrim          = "trim"

	FuncStrToTime  = "strtotime"
	FuncTimeFormat = "timeformat"
	FuncYear       = "year"
	FuncMonth      = "month"
	FuncDay        = "day"
	FuncHour       = "hour"
	FuncMinute     = "minute"
	FuncSecond     = "second"
	FuncDayOfWeek  = "dayofweek"
	FuncDayOfYear  = "dayofyear"
	FuncNewTime    = "newtime"
	FuncSleep      = "sleep"
	FuncNow        = "now"

	FuncDir       = "dir"
	FuncBasename  = "basename"
	FuncExtension = "extension"
	FuncAbsPath   = "abspath"
	FuncRelPath   = "relpath"

	FuncInverse = "inverse"
	FuncEnvOr   = "envor"
	FuncEnv     = "env"
)

func (v TreeVisitor) visitFuncCallExpr(n *FuncCallExpr) (NFunction, error) {
	args := n.Args
	switch n.FnName.L {
	case FuncExpr:
		return v.funcCallExpr(args)
	case FuncLua:
		return v.funcCallLua(args)
	case FuncGrep:
		return v.funcCallGrep(args)
	case FuncSh:
		return v.funcCallSh(args)
	case FuncTmpl:
		return v.funcCallTmpl(args)
	case FuncToInt:
		return v.funcCallToInt(args)
	case FuncToFloat:
		return v.funcCallToFloat(args)
	case FuncToBool:
		return v.funcCallToBool(args)
	case FuncToString:
		return v.funcCallToString(args)
	case FuncToTime:
		return v.funcCallToTime(args)
	case FuncToDuration:
		return v.funcCallToDuration(args)
	case FuncLeast:
		return v.funcCallLeast(args)
	case FuncGreatest:
		return v.funcCallGreatest(args)
	case FuncCoalesce:
		return v.funcCallCoalesce(args)
	case FuncIf:
		return v.funcCallIf(args)
	case FuncIfNull:
		return v.funcCallIfNull(args)
	case FuncNullIf:
		return v.funcCallNullIf(args)
	case FuncInverse:
		return v.funcCallInverse(args)
	case FuncAbs:
		return v.funcCallAbs(args)
	case FuncSqrt:
		return v.funcCallSqrt(args)
	case FuncDegrees:
		return v.funcCallDegrees(args)
	case FuncRadians:
		return v.funcCallRadians(args)
	case FuncAcos:
		return v.funcCallAcos(args)
	case FuncAsin:
		return v.funcCallAsin(args)
	case FuncAtan:
		return v.funcCallAtan(args)
	case FuncCos:
		return v.funcCallCos(args)
	case FuncSin:
		return v.funcCallSin(args)
	case FuncTan:
		return v.funcCallTan(args)
	case FuncCot:
		return v.funcCallCot(args)
	case FuncLn:
		return v.funcCallLn(args)
	case FuncLog2:
		return v.funcCallLog2(args)
	case FuncLog10:
		return v.funcCallLog10(args)
	case FuncExp:
		return v.funcCallExp(args)
	case FuncCeil:
		return v.funcCallCeil(args)
	case FuncFloor:
		return v.funcCallFloor(args)
	case FuncRound:
		return v.funcCallRound(args)
	case FuncAtan2:
		return v.funcCallAtan2(args)
	case FuncPow:
		return v.funcCallPow(args)
	case FuncE:
		return v.funcCallE(args)
	case FuncPi:
		return v.funcCallPi(args)
	case FuncRand:
		return v.funcCallRand(args)
	case FuncLen:
		return v.funcCallLen(args)
	case FuncSize:
		return v.funcCallSize(args)
	case FuncRegexpCount:
		return v.funcCallRegexpCount(args)
	case FuncRegexpInstr:
		return v.funcCallRegexpInstr(args)
	case FuncRegexpSubstr:
		return v.funcCallRegexpSubstr(args)
	case FuncRegexpReplace:
		return v.funcCallRegexpReplace(args)
	case FuncRegexpLike:
		return v.funcCallRegexpLike(args)
	case FuncFormat:
		return v.funcCallFormat(args)
	case FuncLower:
		return v.funcCallLower(args)
	case FuncUpper:
		return v.funcCallUpper(args)
	case FuncSha2:
		return v.funcCallSha2(args)
	case FuncConcatWs:
		return v.funcCallConcatWs(args)
	case FuncInstr:
		return v.funcCallInstr(args)
	case FuncInstrCount:
		return v.funcCallInstrCount(args)
	case FuncSubstr:
		return v.funcCallSubstr(args)
	case FuncSubstrIndex:
		return v.funcCallSubstrIndex(args)
	case FuncReplace:
		return v.funcCallReplace(args)
	case FuncTrim:
		return v.funcCallTrim(args)
	case FuncStrToTime:
		return v.funcCallStrToTime(args)
	case FuncTimeFormat:
		return v.funcCallTimeFormat(args)
	case FuncYear:
		return v.funcCallYear(args)
	case FuncMonth:
		return v.funcCallMonth(args)
	case FuncDay:
		return v.funcCallDay(args)
	case FuncHour:
		return v.funcCallHour(args)
	case FuncMinute:
		return v.funcCallMinute(args)
	case FuncSecond:
		return v.funcCallSecond(args)
	case FuncDayOfWeek:
		return v.funcCallDayOfWeek(args)
	case FuncDayOfYear:
		return v.funcCallDayOfYear(args)
	case FuncNewTime:
		return v.funcCallNewTime(args)
	case FuncSleep:
		return v.funcCallSleep(args)
	case FuncEnvOr:
		return v.funcCallEnvOr(args)
	case FuncEnv:
		return v.funcCallEnv(args)
	case FuncNow:
		return v.funcCallNow(args)
	case FuncDir:
		return v.funcCallDir(args)
	case FuncBasename:
		return v.funcCallBasename(args)
	case FuncExtension:
		return v.funcCallExtension(args)
	case FuncAbsPath:
		return v.funcCallAbsPath(args)
	case FuncRelPath:
		return v.funcCallRelPath(args)
	default:
		return nil, ErrNotImplmented
	}
}

func (TreeVisitor) assertFuncCallArgLen(args []ExprNode, minLen, maxLen int) error {
	if len(args) >= minLen && len(args) <= maxLen {
		return nil
	}
	return fmt.Errorf("%w: argLen should be in [%d, %d] but got %d", ErrInvalidValue, minLen, maxLen, len(args))
}

func (v TreeVisitor) evalFuncCallArgs(args ...ExprNode) ([]NFunction, error) {
	r := make([]NFunction, len(args))
	for i, x := range args {
		f, err := v.VisitExpr(x)
		if err != nil {
			return nil, fmt.Errorf("%w: evalFuncCallArgs[%d]", err, i)
		}
		r[i] = f
	}
	return r, nil
}

func (v TreeVisitor) newNullaryArgUnaryRetFunction(args []ExprNode, name string, f func() (ND, error)) (NFunction, error) {
	if err := v.assertFuncCallArgLen(args, 0, 0); err != nil {
		return nil, err
	}
	return ReturnContainerValue(name, func(_ ND) (ND, error) {
		return f()
	}), nil
}

func (v TreeVisitor) newUnaryArgUnaryRetFunction(args []ExprNode, name string, f func(ND) (ND, error)) (NFunction, error) {
	if err := v.assertFuncCallArgLen(args, 1, 1); err != nil {
		return nil, err
	}
	a, err := v.evalFuncCallArgs(args...)
	if err != nil {
		return nil, err
	}
	return iterx.CombineFunction(a[0], ReturnContainerValue(name, f))
}

func (v TreeVisitor) newVariadicArgUnaryRetFunction(args []ExprNode, name string, minLen, maxLen int, f func(...ND) (ND, error)) (NFunction, error) {
	if err := v.assertFuncCallArgLen(args, minLen, maxLen); err != nil {
		return nil, err
	}
	a, err := v.evalFuncCallArgs(args...)
	if err != nil {
		return nil, err
	}
	if err := ValidateOnlyVariadicOrAllUnaryRet(a...); err != nil {
		return nil, err
	}
	switch {
	case a[0].RetArity() == iterx.Variadic:
		if minLen <= 1 {
			return iterx.CombineFunction(a[0], ReturnContainerValue(name, func(x ND) (ND, error) {
				return f(x)
			}))
		}
		return nil, fmt.Errorf("%w: variadic argument mismatched with [%d, %d]", ErrInvalidArgument, minLen, maxLen)
	default:
		return iterx.NewMapFunction(func(x *N) (*N, error) {
			xs := make([]ND, len(a))
			for i, b := range a {
				c, err := b.CallAny(x)
				if err != nil {
					return nil, fmt.Errorf("%w: function call arg[%d]", err, i)
				}
				_, d, ok := AsValueContainer(c[0]).GetFirstValue()
				if !ok {
					return nil, fmt.Errorf("%w: function call arg[%d] contain no value", ErrInvalidValue, i)
				}
				xs[i] = d
			}
			rd, err := f(xs...)
			if err != nil {
				return nil, err
			}
			r := AsValueContainer(node.New())
			r.SetContainerValue(rd)
			return r.N, nil
		}), nil
	}
}

func (v TreeVisitor) newGeneratorFunction(args []ExprNode, name string, minLen, maxLen int, f func(...ND) (GenTemplate, error)) (NFunction, error) {
	if err := v.assertFuncCallArgLen(args, minLen, maxLen); err != nil {
		return nil, err
	}
	a, err := v.evalFuncCallArgs(args...)
	if err != nil {
		return nil, err
	}
	if err := ValidateOnlyVariadicOrAllUnaryRet(a...); err != nil {
		return nil, err
	}
	switch {
	case a[0].RetArity() == iterx.Variadic:
		if !(minLen <= 1) {
			return nil, fmt.Errorf("%w: generator function variadic argument mismatched with [%d, %d]", ErrInvalidArgument, minLen, maxLen)
		}
		return iterx.NewFanoutFunction(func(x *N) ([]*N, error) {
			aa, err := a[0].CallAny(x)
			if err != nil {
				return nil, fmt.Errorf("%w: generator function call arg[0]", err)
			}
			r := []*N{}
			for _, b := range aa {
				_, d, ok := AsValueContainer(b).GetFirstValue()
				if !ok {
					slog.Warn("generator function call arg[0] contain no value", logx.Err(ErrInvalidValue))
					continue
				}
				g, err := f(d)
				if err != nil {
					slog.Warn("generator function failed to create template", logx.Err(err))
					continue
				}
				rv, err := GenerateAndParse(v.ctx, x, g)
				if err != nil {
					slog.Warn("generator function failed to generate", logx.Err(err))
					continue
				}
				for _, d := range rv {
					y := x.Clone()
					y.Merge(d.Map) // append data to original rows
					z := node.New()
					z.Map = y
					r = append(r, z)
				}
			}
			return r, nil
		}), nil
	default:
		return iterx.NewFanoutFunction(func(x *N) ([]*N, error) {
			xs := make([]ND, len(a))
			for i, b := range a {
				c, err := b.CallAny(x)
				if err != nil {
					return nil, fmt.Errorf("%w: generator function call arg[%d]", err, i)
				}
				_, d, ok := AsValueContainer(c[0]).GetFirstValue()
				if !ok {
					return nil, fmt.Errorf("%w: generator function call arg[%d] contain no value", ErrInvalidValue, i)
				}
				xs[i] = d
			}
			g, err := f(xs...)
			if err != nil {
				return nil, fmt.Errorf("%w: generator function failed to create template", err)
			}
			rv, err := GenerateAndParse(v.ctx, x, g)
			if err != nil {
				return nil, fmt.Errorf("%w: generator function failed to generate", err)
			}
			r := make([]*N, len(rv))
			for i, d := range rv {
				y := x.Clone()
				y.Merge(d.Map) // append data to original rows
				z := node.New()
				z.Map = y
				r[i] = z
			}
			return r, nil
		}), nil
	}
}

// @file: read the file, otherwise use the string as it is.
func (v TreeVisitor) readFileOrString(d ND) ([]byte, error) {
	s, ok := d.AsOp().String()
	if !ok {
		return nil, fmt.Errorf("%w: readFileOrString requires String: %v", ErrInvalidArgument, d)
	}
	src, err := iox.NewFileOrStringSource(s.Raw())
	if err != nil {
		return nil, fmt.Errorf("%w: readFileOrString failed to read", errors.Join(ErrInvalidArgument, err))
	}
	defer func() {
		_ = src.AsReadCloser().Close()
	}()
	b, err := src.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w: readFileOrString failed to read", errors.Join(ErrInvalidArgument, err))
	}
	return b, nil
}

//
// generator
//

func (v TreeVisitor) funcCallExpr(args []ExprNode) (NFunction, error) {
	return v.newGeneratorFunction(args, FuncExpr, 1, 1, func(x ...ND) (GenTemplate, error) {
		b, err := v.readFileOrString(x[0])
		if err != nil {
			return nil, fmt.Errorf("%w: failed to create expr template", err)
		}
		return NewExprGenTemplate(string(b)), nil
	})
}
func (v TreeVisitor) funcCallLua(args []ExprNode) (NFunction, error) {
	return v.newGeneratorFunction(args, FuncLua, 2, 2, func(x ...ND) (GenTemplate, error) {
		b, err := v.readFileOrString(x[0])
		if err != nil {
			return nil, fmt.Errorf("%w: failed to create lua template", err)
		}
		e, ok := x[1].AsOp().String()
		if !ok {
			return nil, fmt.Errorf("%w: lua template requires entrypoint", ErrInvalidArgument)
		}
		return NewLuaGenTemplate(string(b), e.Raw()), nil
	})
}
func (v TreeVisitor) funcCallGrep(args []ExprNode) (NFunction, error) {
	return v.newGeneratorFunction(args, FuncGrep, 2, 2, func(x ...ND) (GenTemplate, error) {
		expr, ok := x[0].AsOp().String()
		if !ok {
			return nil, fmt.Errorf("%w: grep template requires expr", ErrInvalidArgument)
		}
		tmpl, ok := x[1].AsOp().String()
		if !ok {
			return nil, fmt.Errorf("%w: grep template requires tmpl", ErrInvalidArgument)
		}
		return NewRegexpGenTemplate(expr.Raw(), tmpl.Raw()), nil
	})
}
func (v TreeVisitor) funcCallSh(args []ExprNode) (NFunction, error) {
	return v.newGeneratorFunction(args, FuncSh, 1, 1, func(x ...ND) (GenTemplate, error) {
		b, err := v.readFileOrString(x[0])
		if err != nil {
			return nil, fmt.Errorf("%w: failed to create sh template", err)
		}
		return NewShellGenTemplate(string(b)), nil
	})
}
func (v TreeVisitor) funcCallTmpl(args []ExprNode) (NFunction, error) {
	return v.newGeneratorFunction(args, FuncTmpl, 1, 1, func(x ...ND) (GenTemplate, error) {
		b, err := v.readFileOrString(x[0])
		if err != nil {
			return nil, fmt.Errorf("%w: failed to create tmpl template", err)
		}
		return NewStringGenTemplate(string(b)), nil
	})
}

//
// cast
//

func (v TreeVisitor) funcCallToInt(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncToInt, func(x ND) (ND, error) { return x.AsOp().AsInt() })
}
func (v TreeVisitor) funcCallToFloat(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncToFloat, func(x ND) (ND, error) { return x.AsOp().AsFloat() })
}
func (v TreeVisitor) funcCallToBool(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncToBool, func(x ND) (ND, error) { return x.AsOp().AsBool() })
}
func (v TreeVisitor) funcCallToString(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncToString, func(x ND) (ND, error) { return x.AsOp().AsString() })
}
func (v TreeVisitor) funcCallToTime(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncToTime, func(x ND) (ND, error) { return x.AsOp().AsTime() })
}
func (v TreeVisitor) funcCallToDuration(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncToDuration, func(x ND) (ND, error) { return x.AsOp().AsDuration() })
}

//
// common
//

func (v TreeVisitor) funcCallLeast(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncLeast, 1, FuncArgMaxLen,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].Least(x[1:]...), nil
		}))
}
func (v TreeVisitor) funcCallGreatest(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncGreatest, 1, FuncArgMaxLen,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].Greatest(x[1:]...), nil
		}))
}
func (v TreeVisitor) funcCallCoalesce(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncCoalesce, 1, FuncArgMaxLen,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].Coalesce(x[1:]...), nil
		}))
}

//
// control
//

func (v TreeVisitor) funcCallIf(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncIf, 3, 3,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].If(x[1], x[2]), nil
		}))
}
func (v TreeVisitor) funcCallIfNull(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncIfNull, 2, 2,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].IfNull(x[1]), nil
		}))
}
func (v TreeVisitor) funcCallNullIf(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncNullIf, 2, 2,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].NullIf(x[1]), nil
		}))
}

//
// math
//

func (v TreeVisitor) funcCallAbs(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncAbs, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Abs() }))
}
func (v TreeVisitor) funcCallSqrt(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncSqrt, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Sqrt() }))
}
func (v TreeVisitor) funcCallDegrees(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncDegrees, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Degrees() }))
}
func (v TreeVisitor) funcCallRadians(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncRadians, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Radians() }))
}
func (v TreeVisitor) funcCallAcos(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncAcos, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Acos() }))
}
func (v TreeVisitor) funcCallAsin(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncAsin, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Asin() }))
}
func (v TreeVisitor) funcCallAtan(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncAtan, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Atan() }))
}
func (v TreeVisitor) funcCallCos(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncCos, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Cos() }))
}
func (v TreeVisitor) funcCallSin(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncSin, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Sin() }))
}
func (v TreeVisitor) funcCallTan(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncTan, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Tan() }))
}
func (v TreeVisitor) funcCallCot(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncCot, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Cot() }))
}
func (v TreeVisitor) funcCallLn(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncLn, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Ln() }))
}
func (v TreeVisitor) funcCallLog2(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncLog2, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Log2() }))
}
func (v TreeVisitor) funcCallLog10(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncLog10, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Log10() }))
}
func (v TreeVisitor) funcCallExp(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncExp, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Exp() }))
}
func (v TreeVisitor) funcCallCeil(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncCeil, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Ceil() }))
}
func (v TreeVisitor) funcCallFloor(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncFloor, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Floor() }))
}
func (v TreeVisitor) funcCallRound(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncRound, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Round() }))
}
func (v TreeVisitor) funcCallAtan2(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncAtan2, 2, 2,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].Atan2(x[1])
		}),
	)
}
func (v TreeVisitor) funcCallPow(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncPow, 2, 2,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].Pow(x[1])
		}),
	)
}
func (v TreeVisitor) funcCallE(args []ExprNode) (NFunction, error) {
	return v.newNullaryArgUnaryRetFunction(args, FuncE, func() (ND, error) {
		return node.E().AsData(), nil
	})
}
func (v TreeVisitor) funcCallPi(args []ExprNode) (NFunction, error) {
	return v.newNullaryArgUnaryRetFunction(args, FuncPi, func() (ND, error) {
		return node.Pi().AsData(), nil
	})
}
func (v TreeVisitor) funcCallRand(args []ExprNode) (NFunction, error) {
	return v.newNullaryArgUnaryRetFunction(args, FuncRand, func() (ND, error) {
		return node.Rand().AsData(), nil
	})
}

//
// string
//

func (v TreeVisitor) funcCallLen(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncLen, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Len() }))
}
func (v TreeVisitor) funcCallSize(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncSize, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Size() }))
}
func (v TreeVisitor) funcCallFormat(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncFormat, 1, FuncArgMaxLen,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].Format(x[1:]...)
		}),
	)
}
func (v TreeVisitor) funcCallLower(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncLower, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Lower() }))
}
func (v TreeVisitor) funcCallUpper(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncUpper, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Upper() }))
}
func (v TreeVisitor) funcCallSha2(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncSha2, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Sha2() }))
}
func (v TreeVisitor) funcCallRegexpCount(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncRegexpCount, 2, 2,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].RegexpCount(x[1])
		}),
	)
}
func (v TreeVisitor) funcCallRegexpInstr(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncRegexpInstr, 2, 2,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].RegexpInstr(x[1])
		}),
	)
}
func (v TreeVisitor) funcCallRegexpSubstr(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncRegexpSubstr, 2, 2,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].RegexpSubstr(x[1])
		}),
	)
}
func (v TreeVisitor) funcCallRegexpReplace(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncRegexpReplace, 3, 3,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].RegexpReplace(x[1], x[2])
		}),
	)
}
func (v TreeVisitor) funcCallRegexpLike(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncRegexpLike, 2, 2,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].Regexp(x[1])
		}),
	)
}
func (v TreeVisitor) funcCallConcatWs(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncConcatWs, 1, FuncArgMaxLen,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].ConcatWs(x[1:]...)
		}),
	)
}
func (v TreeVisitor) funcCallInstr(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncInstr, 2, 2,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].Instr(x[1])
		}),
	)
}
func (v TreeVisitor) funcCallInstrCount(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncInstrCount, 2, 2,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].InstrCount(x[1])
		}),
	)
}
func (v TreeVisitor) funcCallSubstr(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncSubstr, 2, 3,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].Substr(x[1:]...)
		}),
	)
}
func (v TreeVisitor) funcCallSubstrIndex(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncSubstrIndex, 3, 3,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].SubstrIndex(x[1], x[2])
		}),
	)
}
func (v TreeVisitor) funcCallReplace(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncReplace, 3, 3,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].Replace(x[1], x[2])
		}),
	)
}
func (v TreeVisitor) funcCallTrim(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncTrim, 1, 2,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].Trim(x[1:]...)
		}),
	)
}

//
// time
//

func (v TreeVisitor) funcCallStrToTime(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncStrToTime, 2, 2,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].StrToTime(x[1])
		}),
	)
}
func (v TreeVisitor) funcCallTimeFormat(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncTimeFormat, 2, 2,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].TimeFormat(x[1])
		}),
	)
}
func (v TreeVisitor) funcCallYear(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncYear, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Year() }))
}
func (v TreeVisitor) funcCallMonth(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncMonth, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Month() }))
}
func (v TreeVisitor) funcCallDay(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncDay, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Day() }))
}
func (v TreeVisitor) funcCallHour(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncHour, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Hour() }))
}
func (v TreeVisitor) funcCallMinute(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncMinute, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Minute() }))
}
func (v TreeVisitor) funcCallSecond(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncSecond, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Second() }))
}
func (v TreeVisitor) funcCallDayOfWeek(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncDayOfWeek, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.DayOfWeek() }))
}
func (v TreeVisitor) funcCallDayOfYear(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncDayOfYear, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.DayOfYear() }))
}
func (v TreeVisitor) funcCallNewTime(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncNewTime, 1, 6,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].NewTime(x[1:]...)
		}),
	)
}
func (v TreeVisitor) funcCallSleep(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncSleep, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Sleep() }))
}
func (v TreeVisitor) funcCallNow(args []ExprNode) (NFunction, error) {
	return v.newNullaryArgUnaryRetFunction(args, FuncNow, func() (ND, error) { return node.Now().AsData(), nil })
}

//
// path
//

func (v TreeVisitor) funcCallDir(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncDir, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Dir() }))
}
func (v TreeVisitor) funcCallBasename(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncBasename, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Basename() }))
}
func (v TreeVisitor) funcCallExtension(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncExtension, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Extension() }))
}
func (v TreeVisitor) funcCallAbsPath(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncAbsPath, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.AbsPath() }))
}
func (v TreeVisitor) funcCallRelPath(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncRelPath, 2, 2,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].RelPath(x[1])
		}),
	)
}

//
// etc
//

func (v TreeVisitor) funcCallInverse(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncInverse, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Inverse() }))
}
func (v TreeVisitor) funcCallEnvOr(args []ExprNode) (NFunction, error) {
	return v.newVariadicArgUnaryRetFunction(args, FuncEnvOr, 2, 2,
		AsVariadicArgUnaryRetNodeDataFunction(func(x ...*OP) (*OP, error) {
			return x[0].EnvOr(x[1])
		}),
	)
}
func (v TreeVisitor) funcCallEnv(args []ExprNode) (NFunction, error) {
	return v.newUnaryArgUnaryRetFunction(args, FuncEnv, AsUnaryArgUnaryRetNodeDataFunction(func(x *OP) (*OP, error) { return x.Env() }))
}
