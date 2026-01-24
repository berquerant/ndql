package node

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand/v2"

	"github.com/berquerant/ndql/pkg/util"
)

func (v *Op) numberAsFloat() (Float, bool) {
	switch d := v.data.(type) {
	case Float:
		return d, true
	case Int:
		return Float(float64(d.Raw())), true
	default:
		return Float(0), false
	}
}

func (v *Op) newUnaryArgUnaryRetMathFunc(name string, f func(float64) float64) (*Op, error) {
	x, ok := v.numberAsFloat()
	if !ok {
		return nil, fmt.Errorf("%w: math func %s: invalid argument: %v", ErrUnavailable, name, v)
	}
	if y := f(x.Raw()); util.IsFinite(y) {
		return Float(y).AsOp(), nil
	}
	slog.Warn("math func: invalid argument, returned NULL", slog.String("name", name), slog.Float64("arg", x.Raw()))
	return NewNull().AsOp(), nil
}

func (v *Op) Abs() (*Op, error)  { return v.newUnaryArgUnaryRetMathFunc("abs", math.Abs) }
func (v *Op) Sqrt() (*Op, error) { return v.newUnaryArgUnaryRetMathFunc("sqrt", math.Sqrt) }

func (v *Op) Degrees() (*Op, error) { return v.newUnaryArgUnaryRetMathFunc("degrees", util.ToDegrees) }
func (v *Op) Radians() (*Op, error) { return v.newUnaryArgUnaryRetMathFunc("radians", util.ToRadians) }

func (v *Op) Acos() (*Op, error) { return v.newUnaryArgUnaryRetMathFunc("acos", math.Acos) }
func (v *Op) Asin() (*Op, error) { return v.newUnaryArgUnaryRetMathFunc("asin", math.Asin) }
func (v *Op) Atan() (*Op, error) { return v.newUnaryArgUnaryRetMathFunc("atan", math.Atan) }
func (v *Op) Cos() (*Op, error)  { return v.newUnaryArgUnaryRetMathFunc("cos", math.Cos) }
func (v *Op) Sin() (*Op, error)  { return v.newUnaryArgUnaryRetMathFunc("sin", math.Sin) }
func (v *Op) Tan() (*Op, error)  { return v.newUnaryArgUnaryRetMathFunc("tan", math.Tan) }
func (v *Op) Cot() (*Op, error) {
	return v.newUnaryArgUnaryRetMathFunc("cot", func(x float64) float64 {
		y := math.Tan(x)
		if util.IsNaN(y) || y == 0 {
			return math.NaN()
		}
		return 1 / y
	})
}

func (v *Op) Ln() (*Op, error)    { return v.newUnaryArgUnaryRetMathFunc("ln", math.Log) }
func (v *Op) Log2() (*Op, error)  { return v.newUnaryArgUnaryRetMathFunc("log2", math.Log2) }
func (v *Op) Log10() (*Op, error) { return v.newUnaryArgUnaryRetMathFunc("log10", math.Log10) }
func (v *Op) Exp() (*Op, error)   { return v.newUnaryArgUnaryRetMathFunc("exp", math.Exp) }

func (v *Op) Ceil() (*Op, error)  { return v.newUnaryArgUnaryRetMathFunc("ceil", math.Ceil) }
func (v *Op) Floor() (*Op, error) { return v.newUnaryArgUnaryRetMathFunc("floor", math.Floor) }
func (v *Op) Round() (*Op, error) { return v.newUnaryArgUnaryRetMathFunc("round", math.Round) }

func (v *Op) newBinaryArgUnaryRetMathFunc(name string, other *Op, f func(float64, float64) float64) (*Op, error) {
	x, ok := v.numberAsFloat()
	if !ok {
		return nil, fmt.Errorf("%w: math func %s: invalid argument left %v", ErrUnavailable, name, v)
	}
	y, ok := other.numberAsFloat()
	if !ok {
		return nil, fmt.Errorf("%w: math func %s: invalid argument right %v", ErrUnavailable, name, other)
	}
	if z := f(x.Raw(), y.Raw()); util.IsFinite(z) {
		return Float(z).AsOp(), nil
	}
	slog.Warn("math: func: invalid argument, returned NULL", slog.String("name", name), slog.Float64("left", x.Raw()), slog.Float64("right", y.Raw()))
	return NewNull().AsOp(), nil
}

func (v *Op) Atan2(other *Op) (*Op, error) {
	return v.newBinaryArgUnaryRetMathFunc("atan2", other, math.Atan2)
}

func (v *Op) Pow(other *Op) (*Op, error) {
	return v.newBinaryArgUnaryRetMathFunc("pow", other, math.Pow)
}

//
// const/no args
//

func E() *Op    { return Float(math.E).AsOp() }
func Pi() *Op   { return Float(math.Pi).AsOp() }
func Rand() *Op { return Float(rand.Float64()).AsOp() }
