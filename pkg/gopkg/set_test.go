package gopkg_test

import (
	"encoding/json"
	"testing"

	"github.com/berquerant/ndql/pkg/gopkg"
	"github.com/stretchr/testify/assert"
)

func TestDocumentSet(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		doc := []*gopkg.Document{
			{
				Path:  "k1",
				Text:  "text",
				Title: "t1",
			},
			{
				Path:  "k2",
				Text:  "text",
				Title: "k2",
			},
			{
				Path:  "k1.k11",
				Text:  "text",
				Title: "k1.k11",
			},
			{
				Path:  "k3.k31",
				Text:  "text",
				Title: "k3.k31",
			},
		}
		s := gopkg.NewDocumentSet(doc...)
		b, err := json.Marshal(s)
		if !assert.Nil(t, err) {
			return
		}
		var v gopkg.DocumentSet
		if !assert.Nil(t, json.Unmarshal(b, &v)) {
			return
		}
		assert.Equal(t, s, &v)
	})
}
