package util_test

import (
	"errors"
	"testing"

	"github.com/berquerant/ndql/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestNoError(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		assert.True(t, util.NoError(1, nil))
	})
	t.Run("fail", func(t *testing.T) {
		assert.False(t, util.NoError(1, errors.New("Error")))
	})
}

func TestOK(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		assert.True(t, util.OK(1, true))
	})
	t.Run("false", func(t *testing.T) {
		assert.False(t, util.OK(1, false))
	})
}
