package node

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/berquerant/ndql/pkg/iox"
	"github.com/berquerant/ndql/pkg/mapx"
	"github.com/berquerant/ndql/pkg/util"
)

// Node is a set of Data; what we call 'row' in SQL.
type Node struct {
	*mapx.Map[string, Data]
}

func New() *Node {
	return &Node{mapx.NewMap[string, Data](nil)}
}

func FromMap(v map[string]Data) *Node { return &Node{mapx.NewMap(v)} }

func FromWalkerEntry(v *iox.WalkerEntry) *Node {
	return &Node{
		mapx.NewMap(map[string]Data{
			KeyPath:    String(v.Path),
			KeySize:    Int(v.Size),
			KeyIsDir:   Bool(v.IsDir),
			KeyModTime: Time(v.ModTime),
			KeyMode:    String(v.Mode.String()),
		}),
	}
}

var ErrInvalidNode = errors.New("InvalidNode")

// Validate returns an error if the node does not have all the builtin keys.
func (n *Node) Validate() error {
	switch {
	case !util.OK(n.Get(KeyPath)):
		return fmt.Errorf("%w: no %s", ErrInvalidNode, KeyPath)
	case !util.OK(n.Get(KeySize)):
		return fmt.Errorf("%w: no %s", ErrInvalidNode, KeySize)
	case !util.OK(n.Get(KeyIsDir)):
		return fmt.Errorf("%w: no %s", ErrInvalidNode, KeyIsDir)
	case !util.OK(n.Get(KeyModTime)):
		return fmt.Errorf("%w: no %s", ErrInvalidNode, KeyModTime)
	case !util.OK(n.Get(KeyMode)):
		return fmt.Errorf("%w: no %s", ErrInvalidNode, KeyMode)
	default:
		switch {
		case !util.OK(util.MustOK(n.Get(KeyPath)).AsOp().String()):
			return fmt.Errorf("%w: invalid %s", ErrInvalidNode, KeyPath)
		case !util.OK(util.MustOK(n.Get(KeySize)).AsOp().Int()):
			return fmt.Errorf("%w: invalid %s", ErrInvalidNode, KeySize)
		case !util.OK(util.MustOK(n.Get(KeyIsDir)).AsOp().Bool()):
			return fmt.Errorf("%w: invalid %s", ErrInvalidNode, KeyIsDir)
		case !util.OK(util.MustOK(n.Get(KeyModTime)).AsOp().Time()):
			return fmt.Errorf("%w: invalid %s", ErrInvalidNode, KeyModTime)
		case !util.OK(util.MustOK(n.Get(KeyMode)).AsOp().String()):
			return fmt.Errorf("%w: invalid %s", ErrInvalidNode, KeyMode)
		default:
			return nil
		}
	}
}

func (n *Node) GetPath() String {
	v, _ := n.Get(KeyPath)
	return v.(String)
}

func (n *Node) GetSize() Int {
	v, _ := n.Get(KeySize)
	return v.(Int)
}

func (n *Node) GetIsDir() Bool {
	v, _ := n.Get(KeyIsDir)
	return v.(Bool)
}

func (n *Node) GetModTime() Time {
	v, _ := n.Get(KeyModTime)
	return v.(Time)
}

func (n *Node) GetMode() String {
	v, _ := n.Get(KeyMode)
	return v.(String)
}

func (n *Node) MarshalJSON() ([]byte, error) {
	d := make(map[string]*Op, n.Len())
	for _, k := range n.Keys() {
		v, _ := n.Get(k)
		d[k] = v.AsOp()
	}
	return json.Marshal(d)
}

func (n *Node) UnmarshalJSON(data []byte) error {
	d := map[string]*Op{}
	if err := json.Unmarshal(data, &d); err != nil {
		return nil
	}
	m := make(map[string]Data, len(d))
	for k, v := range d {
		if v == nil {
			v = NewNull().AsOp()
		}
		m[k] = v.AsData()
	}
	r := New()
	*r.Map = *mapx.NewMap(m)
	*n = *r
	return nil
}
