package node_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/berquerant/ndql/pkg/node"
	"github.com/berquerant/ndql/pkg/util"
	"github.com/stretchr/testify/assert"
)

type unaryOpTestcase struct {
	v    node.Data
	want node.Data
	err  error
}

func newFailedUnaryOpTestcases(v ...node.Data) []*unaryOpTestcase {
	xs := make([]*unaryOpTestcase, len(v))
	for i, x := range v {
		xs[i] = &unaryOpTestcase{
			v: x,
		}
	}
	return xs
}

func (tc unaryOpTestcase) title() string { return tc.v.Display() }

func runUnaryOpTest(t *testing.T, f func(node.Data) (node.Data, error), cases ...[]*unaryOpTestcase) {
	for _, tcs := range cases {
		for _, tc := range tcs {
			t.Run(tc.title(), func(t *testing.T) {
				got, err := f(tc.v)
				if tc.want == nil {
					assert.NotNil(t, err, "should be failed")
					if tc.err != nil {
						assert.ErrorIs(t, err, tc.err)
					}
					return
				}
				assert.Nil(t, err, "should be succeeded")
				assert.Equal(t, tc.want, got)
			})
		}
	}
}

type binaryOpTestacase struct {
	left, right node.Data
	want        node.Data
	err         error
}

func (tc binaryOpTestacase) title() string {
	return fmt.Sprintf("%s_%s", tc.left.Display(), tc.right.Display())
}

func newFailedBinaryOpTestcases(leftRight ...node.Data) []*binaryOpTestacase {
	var (
		i  int
		xs []*binaryOpTestacase
	)
	for i < len(leftRight) {
		left := leftRight[i]
		i++
		if i < len(leftRight) {
			right := leftRight[i]
			i++
			xs = append(xs, &binaryOpTestacase{
				left:  left,
				right: right,
			})
		}
		break
	}
	return xs
}

func runBinaryOpTest(t *testing.T, f func(node.Data, node.Data) (node.Data, error), cases ...[]*binaryOpTestacase) {
	for _, tcs := range cases {
		for _, tc := range tcs {
			t.Run(tc.title(), func(t *testing.T) {
				got, err := f(tc.left, tc.right)
				if tc.want == nil {
					assert.NotNil(t, err, "should be failed")
					if tc.err != nil {
						assert.ErrorIs(t, err, tc.err)
					}
					return
				}
				assert.Nil(t, err, "should be succeeded")
				assert.Equal(t, tc.want, got)
			})
		}
	}
}

type failedBinaryOpTestcaseBuilder struct {
	*failedTestcaseSeed
	acc []*binaryOpTestacase
}

func newFailedBinaryOpTestcaseBuilder(seed *failedTestcaseSeed) *failedBinaryOpTestcaseBuilder {
	return &failedBinaryOpTestcaseBuilder{
		failedTestcaseSeed: seed,
	}
}

func defaultFailedBinaryOpTestcaseBuilder() *failedBinaryOpTestcaseBuilder {
	return newFailedBinaryOpTestcaseBuilder(defaultFailedTestcaseSeed())
}

func (bb *failedBinaryOpTestcaseBuilder) build() []*binaryOpTestacase {
	d := map[string]bool{}
	r := []*binaryOpTestacase{}
	for _, x := range bb.acc {
		k := x.title()
		if _, ok := d[k]; ok {
			continue
		}
		d[k] = true
		r = append(r, x)
	}
	return r
}

// Add a test case.
func (bb *failedBinaryOpTestcaseBuilder) add(left, right node.Data) *failedBinaryOpTestcaseBuilder {
	bb.acc = append(bb.acc, &binaryOpTestacase{
		left:  left,
		right: right,
	})
	return bb
}

// Add all pairs from the given data list to the test cases.
func (bb *failedBinaryOpTestcaseBuilder) perm(v ...node.Data) *failedBinaryOpTestcaseBuilder {
	for _, xs := range util.Perm2(v...) {
		bb.add(xs[0], xs[1])
	}
	return bb
}

// Add all pairs includes null to the test cases.
func (bb *failedBinaryOpTestcaseBuilder) nullPerm() *failedBinaryOpTestcaseBuilder {
	return bb.pairPerm(bb.n(), bb.all()...)
}

// Add all pairs of the left and each right to the test cases.
func (bb *failedBinaryOpTestcaseBuilder) pairPerm(left node.Data, right ...node.Data) *failedBinaryOpTestcaseBuilder {
	for _, v := range right {
		bb.perm(left, v)
	}
	return bb
}

func (bb *failedBinaryOpTestcaseBuilder) pairPermExcept(left node.Data, exceptRight ...node.Data) *failedBinaryOpTestcaseBuilder {
	v := slices.DeleteFunc(bb.all(), func(x node.Data) bool {
		return slices.Contains(exceptRight, x)
	})
	return bb.pairPerm(left, v...)
}

type failedTestcaseSeed struct {
	Null     node.Null
	Float    node.Float
	Int      node.Int
	String   node.String
	Bool     node.Bool
	Time     node.Time
	Duration node.Duration
}

func (s failedTestcaseSeed) all() []node.Data {
	return []node.Data{
		s.Null,
		s.Float,
		s.Int,
		s.String,
		s.Bool,
		s.Time,
		s.Duration,
	}
}

func (s failedTestcaseSeed) except(v ...node.Data) []node.Data {
	return slices.DeleteFunc(s.all(), func(x node.Data) bool {
		return slices.Contains(v, x)
	})
}

func defaultFailedTestcaseSeed() *failedTestcaseSeed {
	d := node.Default()
	return &failedTestcaseSeed{
		Null:     d.Null(),
		Float:    d.Float(),
		Int:      d.Int(),
		String:   d.String(),
		Time:     d.Time(),
		Duration: d.Duration(),
	}
}

func (s failedTestcaseSeed) n() node.Null     { return s.Null }
func (s failedTestcaseSeed) f() node.Float    { return s.Float }
func (s failedTestcaseSeed) i() node.Int      { return s.Int }
func (e failedTestcaseSeed) s() node.String   { return e.String }
func (s failedTestcaseSeed) b() node.Bool     { return s.Bool }
func (s failedTestcaseSeed) t() node.Time     { return s.Time }
func (s failedTestcaseSeed) d() node.Duration { return s.Duration }

type variadicOpTestacase struct {
	v    []node.Data
	want node.Data
	err  error
}

func (tc variadicOpTestacase) title() string {
	xs := make([]string, len(tc.v))
	for i, x := range tc.v {
		xs[i] = x.Display()
	}
	return strings.Join(xs, "_")
}

func runVariadicOpTest(t *testing.T, f func(...node.Data) (node.Data, error), cases ...[]*variadicOpTestacase) {
	for _, tcs := range cases {
		for _, tc := range tcs {
			t.Run(tc.title(), func(t *testing.T) {
				got, err := f(tc.v...)
				if tc.want == nil {
					assert.NotNil(t, err, "should be failed")
					if tc.err != nil {
						assert.ErrorIs(t, err, tc.err)
					}
					return
				}
				assert.Nil(t, err, "should be succeeded")
				assert.Equal(t, tc.want, got)
			})
		}
	}
}

func DataListToOpList(v ...node.Data) []*node.Op {
	xs := make([]*node.Op, len(v))
	for i, x := range v {
		xs[i] = x.AsOp()
	}
	return xs
}

func OpListToDataList(v ...*node.Op) []node.Data {
	xs := make([]node.Data, len(v))
	for i, x := range v {
		xs[i] = x.AsData()
	}
	return xs
}
