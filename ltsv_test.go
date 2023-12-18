package cv_test

import (
	"strings"
	"testing"

	"github.com/berquerant/cv"
	"github.com/stretchr/testify/assert"
)

func TestLTSVTranslator(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		for _, tc := range []struct {
			title string
			data  any
			want  []byte
		}{
			{
				title: "empty",
				data:  []map[string]any{},
				want:  nil,
			},
			{
				title: "1 row",
				data: []map[string]any{
					{
						"a": "b",
					},
				},
				want: []byte(`a:b
`),
			},
			{
				title: "2 rows",
				data: []map[string]any{
					{
						"a1": "b1",
						"c1": "d1",
					},
					{
						"a2": "b2",
						"c2": "d2",
					},
				},
				want: []byte(strings.Join([]string{
					"a1:b1\tc1:d1",
					"a2:b2\tc2:d2",
				}, "\n") + "\n"),
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				translator := cv.NewLTSVTranslator(':')
				got, err := translator.Marshal(tc.data)
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got, "%s != %s", tc.want, got)
			})
		}
	})

	t.Run("Unmarshal", func(t *testing.T) {
		for _, tc := range []struct {
			title string
			data  []byte
			want  []map[string]string
		}{
			{
				title: "nil",
				data:  nil,
				want:  nil,
			},
			{
				title: "1 row",
				data:  []byte(`a:b`),
				want: []map[string]string{
					{
						"a": "b",
					},
				},
			},
			{
				title: "2 rows",
				data: []byte(strings.Join([]string{
					"a1:b1\tc1:d1",
					"a2:b2\tc2:d2",
				}, "\n") + "\n"),
				want: []map[string]string{
					{
						"a1": "b1",
						"c1": "d1",
					},
					{
						"a2": "b2",
						"c2": "d2",
					},
				},
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				translator := cv.NewLTSVTranslator(':')
				var v []map[string]string
				assert.Nil(t, translator.Unmarshal(tc.data, &v))
				assert.Equal(t, tc.want, v)
			})
		}
	})
}
