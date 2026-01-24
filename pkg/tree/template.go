package tree

import (
	"context"
	"os"
)

type GenTemplate interface {
	// Generate attributes for new nodes from an existing node.
	// The function must return a string in one of the following formats:
	//
	//   - An array of JSON objects
	//   - A single JSON object
	//   - An "equal pair" list
	//
	// The "equal pair" format is as follows:
	//
	//      key1=value11,key2=value12,...
	//      key1=value21,key2=value22,...
	//      ...
	//
	// This is equivalent to the following JSON structure:
	//
	//      [
	//        {"key1":"value11","key2":"value12",...},
	//        {"key1":"value21","key2":"value22",...},
	//        ...
	//      ]
	//
	// Each JSON object corresponds to a single node.
	// Note that nodes are not required to have the same set of keys.
	Generate(ctx context.Context, n *N) ([]byte, error)
}

func GenerateAndParse(ctx context.Context, n *N, g GenTemplate) ([]*N, error) {
	b, err := g.Generate(ctx, n)
	if err != nil {
		return nil, err
	}
	return ParseGenResult(b)
}

func NodeAsEnviron(n *N) []string {
	xs := []string{}
	for k, v := range n.Unwrap() {
		s, err := v.AsOp().AsString()
		if err != nil {
			continue
		}
		xs = append(xs, k+"="+s.Raw())
	}
	return append(os.Environ(), xs...)
}

func genTemplateEnvOr(key, v string) string {
	if x := os.Getenv(key); x != "" {
		return x
	}
	return v
}

func NodeAsStructuredMap(n *N) map[string]any {
	r := make(map[string]any)
	for k, v := range n.Unwrap() {
		value := v.Any()
		key := KeyFromString(k)
		if key.Table == "" {
			r[k] = value
			continue
		}
		if d, ok := r[key.Table]; ok {
			if dv, ok := d.(map[string]any); ok {
				dv[key.Column] = v
				continue
			}
		}
		r[key.Table] = map[string]any{
			key.Column: v,
		}
	}
	return r
}
