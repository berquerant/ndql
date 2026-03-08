package gopkg_test

import (
	"testing"

	"github.com/berquerant/ndql/pkg/gopkg"
	"github.com/stretchr/testify/assert"
)

func TestComment(t *testing.T) {
	t.Run("GetDocument", func(t *testing.T) {
		for _, tc := range []struct {
			title   string
			comment *gopkg.Comment
			want    *gopkg.Document
			err     error
		}{
			{
				title: "empty comment",
				comment: &gopkg.Comment{
					Text: "",
				},
				err: gopkg.ErrDocument,
			},
			{
				title: "empty path",
				comment: &gopkg.Comment{
					Text: "@path",
				},
				err: gopkg.ErrDocument,
			},
			{
				title: "no document",
				comment: &gopkg.Comment{
					Text: "@path somepath",
				},
				err: gopkg.ErrDocument,
			},
			{
				title: "path is after document",
				comment: &gopkg.Comment{
					Text: `@document
@path somepath`,
				},
				err: gopkg.ErrDocument,
			},
			{
				title: "empty document",
				comment: &gopkg.Comment{
					Text: `@path somepath
@document`,
				},
				err: gopkg.ErrDocument,
			},
			{
				title: "accept",
				comment: &gopkg.Comment{
					Text: `@path somepath
@title somtitle
@document
content`,
				},
				want: &gopkg.Document{
					Path:  "somepath",
					Text:  "content",
					Title: "somtitle",
				},
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				got, err := tc.comment.GetDocument()
				if tc.err != nil {
					assert.ErrorIs(t, err, tc.err)
					return
				}
				if !assert.Nil(t, err) {
					return
				}
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("GetAnnotation", func(t *testing.T) {
		for _, tc := range []struct {
			title   string
			comment *gopkg.Comment
			key     string
			want    *gopkg.Annotation
			ok      bool
		}{
			{
				title: "empty comment",
				comment: &gopkg.Comment{
					Text: "",
				},
				key: "key",
				ok:  false,
			},
			{
				title: "key not found",
				comment: &gopkg.Comment{
					Text: "document",
				},
				key: "key",
				ok:  false,
			},
			{
				title: "key is not at head of line",
				comment: &gopkg.Comment{
					Text: "some @key",
				},
				key: "key",
				ok:  false,
			},
			{
				title: "hit",
				comment: &gopkg.Comment{
					Text: "@key",
				},
				key: "key",
				want: &gopkg.Annotation{
					Key:   "key",
					Value: "",
					Linum: 0,
				},
				ok: true,
			},
			{
				title: "hit with value",
				comment: &gopkg.Comment{
					Text: "@key value",
				},
				key: "key",
				want: &gopkg.Annotation{
					Key:   "key",
					Value: "value",
					Linum: 0,
				},
				ok: true,
			},
			{
				title: "hit at line 2",
				comment: &gopkg.Comment{
					Text: `some
@key v1 v2`,
				},
				key: "key",
				want: &gopkg.Annotation{
					Key:   "key",
					Value: "v1 v2",
					Linum: 1,
				},
				ok: true,
			},
			{
				title: "hit at line 2 and 3",
				comment: &gopkg.Comment{
					Text: `some
@key v1 v2
@key v3`,
				},
				key: "key",
				want: &gopkg.Annotation{
					Key:   "key",
					Value: "v1 v2",
					Linum: 1,
				},
				ok: true,
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				got, ok := tc.comment.GetAnnotation(tc.key)
				if !tc.ok {
					assert.False(t, ok)
					return
				}
				if !assert.True(t, ok) {
					return
				}
				assert.Equal(t, tc.want, got)
			})
		}
	})
}
