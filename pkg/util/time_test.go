package util_test

import (
	"testing"
	"time"

	"github.com/berquerant/ndql/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestNewDuration(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		assert.Equal(t, time.Duration(1500000000), util.NewDuration(float64(1.5), time.Second))
	})
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, 2*time.Second, util.NewDuration(int64(2), time.Second))
	})
	t.Run("uint", func(t *testing.T) {
		assert.Equal(t, 2*time.Second, util.NewDuration(uint64(2), time.Second))
	})
	t.Run("uintprt", func(t *testing.T) {
		assert.Equal(t, 2*time.Second, util.NewDuration(uintptr(2), time.Second))
	})
}
