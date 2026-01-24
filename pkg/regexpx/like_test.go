package regexpx_test

import (
	"testing"

	"github.com/berquerant/ndql/pkg/regexpx"
	"github.com/stretchr/testify/assert"
)

func TestLikeToRegexpString(t *testing.T) {
	for _, tc := range []struct {
		s      string
		escape rune
		want   string
	}{
		{
			s:      "abc",
			escape: '|',
			want:   "abc",
		},
		{
			s:      "abc_",
			escape: '|',
			want:   "abc.",
		},
		{
			s:      "abc|_",
			escape: '|',
			want:   "abc_",
		},
		{
			s:      "abc%",
			escape: '|',
			want:   "abc.*",
		},
		{
			s:      "abc|%",
			escape: '|',
			want:   "abc%",
		},
	} {
		t.Run(tc.s, func(t *testing.T) {
			assert.Equal(t, tc.want, regexpx.LikeToRegexpString(tc.s, tc.escape))
		})
	}
}
