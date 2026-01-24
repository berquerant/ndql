package node_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/berquerant/ndql/pkg/node"
	"github.com/berquerant/ndql/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	v := node.FromMap(map[string]node.Data{
		"n": node.NewNull(),
		"f": node.Float(2.5),
		"i": node.Int(1),
		"b": node.Bool(true),
		"s": node.String("str"),
		"t": node.Time(util.Must(time.Parse(time.DateTime, "2026-01-02 03:04:05"))),
		"d": node.Duration(util.Must(time.ParseDuration("2s"))),
	})
	s := `{"n":null,"f":2.5,"i":1,"b":true,"s":"str","t":"2026-01-02 03:04:05","d":"2s"}`
	var wantMap map[string]any
	util.FailOnError(json.Unmarshal([]byte(s), &wantMap))

	t.Run("marshal", func(t *testing.T) {
		b, err := json.Marshal(v)
		if !assert.Nil(t, err, fmt.Sprintf("%v", err)) {
			return
		}
		var gotMap map[string]any
		if !assert.Nil(t, json.Unmarshal(b, &gotMap)) {
			return
		}
		assert.Equal(t, wantMap, gotMap)
	})
	t.Run("unmarshal", func(t *testing.T) {
		n := node.New()
		if err := json.Unmarshal([]byte(s), n); !assert.Nil(t, err, fmt.Sprintf("%v", err)) {
			return
		}
		assert.Equal(t, v, n)
	})
}
