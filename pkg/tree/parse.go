package tree

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/berquerant/ndql/pkg/node"
)

// Parse the result of GenTemplate.Generate.
func ParseGenResult(b []byte) ([]*N, error) {
	r, err := parseGenResult(b)
	if err != nil {
		return nil, errors.Join(ErrParseGenResult, err)
	}
	return r, nil
}

func parseGenResult(b []byte) ([]*N, error) {
	switch {
	case bytes.HasPrefix(b, []byte("[")):
		return parseGenResultMultipleNodes(b)
	case bytes.HasPrefix(b, []byte("{")):
		return parseGenResultSingleNode(b)
	default:
		return parseGenResultEqualPairs(b)
	}
}

func parseGenResultMultipleNodes(b []byte) ([]*N, error) {
	ns := []*N{}
	if err := json.Unmarshal(b, &ns); err != nil {
		return nil, err
	}
	return ns, nil
}

func parseGenResultSingleNode(b []byte) ([]*N, error) {
	n := node.New()
	if err := json.Unmarshal(b, n); err != nil {
		return nil, err
	}
	return []*N{n}, nil
}

// Parse
//
//	k1=v1,k2=v2,...
//	k1=v2
//
// into
//
//	[{"k1":"v1", "k2": "v2"}, {"k1": "v2"}]
func parseGenResultEqualPairs(b []byte) ([]*N, error) {
	r := []*N{}
	sc := bufio.NewScanner(bytes.NewBuffer(b))
	for sc.Scan() {
		n := node.New()
		for _, p := range strings.Split(sc.Text(), ",") {
			ss := strings.SplitN(p, "=", 2)
			if len(ss) == 2 {
				n.Set(ss[0], node.String(ss[1]))
			}
		}
		r = append(r, n)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return r, nil
}
