package node

import (
	"crypto/sha256"
	"fmt"
	"slices"
	"strings"

	"github.com/berquerant/ndql/pkg/regexpx"
)

func (v *Op) Len() (*Op, error) {
	switch d := v.data.(type) {
	case String:
		return Int(int64(len([]rune(d.Raw())))).AsOp(), nil
	}
	return nil, unavailable("Len", v)
}

func (v *Op) Size() (*Op, error) {
	switch d := v.data.(type) {
	case String:
		return Int(int64(len(d.Raw()))).AsOp(), nil
	}
	return nil, unavailable("Size", v)
}

func (v *Op) Like(other *Op) (*Op, error) {
	if d, ok := v.String(); ok {
		if e, ok := other.String(); ok {
			expr, err := regexpx.Compile(regexpx.LikeToRegexpStringDefault(e.Raw()))
			if err != nil {
				return nil, withUnavailableErr(err, "Like", []*Op{v, other}, "invalid expr")
			}
			return Bool(expr.MatchString(d.Raw())).AsOp(), nil
		}
	}
	return nil, unavailable("Like", v, other)
}

func (v *Op) Regexp(other *Op) (*Op, error) {
	if d, ok := v.String(); ok {
		if e, ok := other.String(); ok {
			expr, err := regexpx.Compile(e.Raw())
			if err != nil {
				return nil, withUnavailableErr(err, "Regexp", []*Op{v, other}, "invalid expr")
			}
			return Bool(expr.MatchString(d.Raw())).AsOp(), nil
		}
	}
	return nil, unavailable("Regexp", v, other)
}

func (v *Op) RegexpCount(other *Op) (*Op, error) {
	if d, ok := v.String(); ok {
		if e, ok := other.String(); ok {
			expr, err := regexpx.Compile(e.Raw())
			if err != nil {
				return nil, withUnavailableErr(err, "RegexpCount", []*Op{v, other}, "invalid expr")
			}
			r := expr.FindAllStringSubmatchIndex(d.Raw(), -1)
			return Int(int64(len(r))).AsOp(), nil
		}
	}
	return nil, unavailable("RegexpCount", v, other)
}

func (v *Op) RegexpInstr(other *Op) (*Op, error) {
	if d, ok := v.String(); ok {
		if e, ok := other.String(); ok {
			expr, err := regexpx.Compile(e.Raw())
			if err != nil {
				return nil, withUnavailableErr(err, "RegexpInstr", []*Op{v, other}, "invalid expr")
			}
			r := expr.FindStringIndex(d.Raw())
			if len(r) == 0 {
				return Int(0).AsOp(), nil
			}
			return Int(int64(r[0] + 1)).AsOp(), nil
		}
	}
	return nil, unavailable("RegexpInstr", v, other)
}

func (v *Op) RegexpSubstr(other *Op) (*Op, error) {
	if d, ok := v.String(); ok {
		if e, ok := other.String(); ok {
			expr, err := regexpx.Compile(e.Raw())
			if err != nil {
				return nil, withUnavailableErr(err, "RegexpSubstr", []*Op{v, other}, "invalid expr")
			}
			r := expr.FindString(d.Raw())
			return String(r).AsOp(), nil
		}
	}
	return nil, unavailable("RegexpSubstr", v, other)
}

func (v *Op) RegexpReplace(pat, repl *Op) (*Op, error) {
	if d, ok := v.String(); ok {
		patStr, ok := pat.String()
		if !ok {
			return nil, withUnavailable("RegexpReplace", []*Op{v, pat, repl}, "no expr")
		}
		replStr, ok := repl.String()
		if !ok {
			return nil, withUnavailable("RegexpReplace", []*Op{v, pat, repl}, "no repl")
		}
		expr, err := regexpx.Compile(patStr.Raw())
		if err != nil {
			return nil, withUnavailableErr(err, "RegexpReplace", []*Op{v, pat, repl}, "invalid expr")
		}
		return String(expr.ReplaceAllString(d.Raw(), replStr.Raw())).AsOp(), nil
	}
	return nil, unavailable("RegexpReplace", v, pat, repl)
}

func (v *Op) Format(other ...*Op) (*Op, error) {
	switch d := v.data.(type) {
	case String:
		a := make([]any, len(other))
		for i, x := range other {
			a[i] = x.AsData().Any()
		}
		return String(fmt.Sprintf(d.Raw(), a...)).AsOp(), nil
	}
	return nil, unavailable("Format", append([]*Op{v}, other...)...)
}

func (v *Op) Lower() (*Op, error) {
	switch d := v.data.(type) {
	case String:
		return String(strings.ToLower(d.Raw())).AsOp(), nil
	default:
		return nil, unavailable("Lower", v)
	}
}

func (v *Op) Upper() (*Op, error) {
	switch d := v.data.(type) {
	case String:
		return String(strings.ToUpper(d.Raw())).AsOp(), nil
	default:
		return nil, unavailable("Upper", v)
	}
}

func (v *Op) Sha2() (*Op, error) {
	switch d := v.data.(type) {
	case String:
		return String(fmt.Sprintf("%x", sha256.Sum256([]byte(d.Raw())))).AsOp(), nil
	default:
		return nil, unavailable("Size", v)
	}
}

func (v *Op) ConcatWs(other ...*Op) (*Op, error) {
	switch d := v.data.(type) {
	case String:
		xs := make([]string, len(other))
		for i, x := range other {
			s, ok := x.String()
			if !ok {
				return nil, withUnavailable("ConcatWs", append([]*Op{v}, other...), "arg[%d] is not a string", i)
			}
			xs[i] = s.Raw()
		}
		return String(strings.Join(xs, d.Raw())).AsOp(), nil
	default:
		return nil, unavailable("ConcatWs", append([]*Op{v}, other...)...)
	}
}

func (v *Op) Instr(other *Op) (*Op, error) {
	switch d := v.data.(type) {
	case String:
		switch e := other.data.(type) {
		case String:
			return Int(int64(strings.Index(d.Raw(), e.Raw()) + 1)).AsOp(), nil
		}
	}
	return nil, unavailable("Instr", v, other)
}

func (v *Op) InstrCount(other *Op) (*Op, error) {
	switch d := v.data.(type) {
	case String:
		switch e := other.data.(type) {
		case String:
			return Int(int64(strings.Count(d.Raw(), e.Raw()))).AsOp(), nil
		}
	}
	return nil, unavailable("Instr", v, other)
}

const substrMaxLength = 2048

func (v *Op) Substr(other ...*Op) (*Op, error) {
	switch d := v.data.(type) {
	case String:
		var (
			length = substrMaxLength
		)
		if len(other) < 1 || len(other) > 2 {
			return nil, withUnavailable("Substr", append([]*Op{v}, other...), "arg len is not in [2, 3]")
		}
		i, ok := other[0].Int()
		if !ok {
			return nil, withUnavailable("Substr", append([]*Op{v}, other...), "arg[1] should be Int")
		}
		pos := int(i.Raw())
		if len(other) == 2 {
			i, ok := other[1].Int()
			if !ok {
				return nil, withUnavailable("Substr", append([]*Op{v}, other...), "arg[2] should be Int")
			}
			length = int(i.Raw())
		}
		s := []rune(d.Raw())
		negPos := pos < 0
		for pos < 0 {
			pos += len(s)
		}
		if pos == 0 {
			if negPos {
				return v, nil
			}
			if length < 1 {
				return String("").AsOp(), nil
			}
		}
		start := pos - 1
		if negPos {
			start++
		}
		end := min(len(s), start+length)
		r := s[start:end]
		return String(string(r)).AsOp(), nil
	}
	return nil, unavailable("Substr", append([]*Op{v}, other...)...)
}

func (v *Op) SubstrIndex(delim, count *Op) (*Op, error) {
	switch d := v.data.(type) {
	case String:
		delimStr, ok := delim.String()
		if !ok {
			return nil, withUnavailable("SubstrIndex", []*Op{v, delim, count}, "delim should be String")
		}
		countInt, ok := count.Int()
		if !ok {
			return nil, withUnavailable("Substr", []*Op{v, delim, count}, "count should be Int")
		}
		xs := strings.Split(d.Raw(), delimStr.Raw())
		countAbs := countInt.Raw()
		if countAbs < 0 {
			countAbs *= -1
		}
		if countInt.Raw() < 0 {
			slices.Reverse(xs)
		}
		end := min(len(xs), int(countAbs))
		rs := xs[:end]
		if countInt.Raw() < 0 {
			slices.Reverse(rs)
		}
		return String(strings.Join(rs, delimStr.Raw())).AsOp(), nil
	}
	return nil, unavailable("SubstrIndex", v, delim, count)
}

func (v *Op) Replace(from, to *Op) (*Op, error) {
	switch d := v.data.(type) {
	case String:
		fromStr, ok := from.String()
		if !ok {
			return nil, withUnavailable("Replace", []*Op{v, from, to}, "from should be String")
		}
		toStr, ok := to.String()
		if !ok {
			return nil, withUnavailable("Replace", []*Op{v, from, to}, "to should be String")
		}
		return String(strings.ReplaceAll(d.Raw(), fromStr.Raw(), toStr.Raw())).AsOp(), nil
	}
	return nil, unavailable("Replace", v, from, to)
}

func (v *Op) Trim(other ...*Op) (*Op, error) {
	switch d := v.data.(type) {
	case String:
		switch len(other) {
		case 0:
			return String(strings.TrimSpace(d.Raw())).AsOp(), nil
		case 1:
			e, ok := other[0].String()
			if !ok {
				return nil, withUnavailable("Trim", append([]*Op{v}, other...), "arg[1] should be String")
			}
			return String(strings.Trim(d.Raw(), e.Raw())).AsOp(), nil
		}
	}
	return nil, unavailable("Trim", append([]*Op{v}, other...)...)
}
