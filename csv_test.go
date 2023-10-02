package cv_test

import (
	"testing"

	"github.com/berquerant/cv"
	"github.com/stretchr/testify/assert"
)

func TestCSVTranslator(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		for _, tc := range []struct {
			title string
			data  any
			want  []byte
		}{
			{
				title: "empty",
				data:  [][]string{},
				want:  nil,
			},
			{
				title: "ints",
				data: [][]int{
					[]int{1, 2, 3},
				},
				want: []byte(`1,2,3
`),
			},
			{
				title: "strings",
				data: [][]string{
					[]string{"a", "b", "c"},
					[]string{"d", "e", "f"},
				},
				want: []byte(`a,b,c
d,e,f
`),
			},
			{
				title: "anys",
				data: [][]any{
					[]any{true, 1, 1.1},
					[]any{false, uint(10), "100"},
				},
				want: []byte(`true,1,1.1
false,10,100
`),
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				translator := cv.NewCSVTranslator(',')
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
			want  [][]string
		}{
			{
				title: "nil",
				data:  nil,
				want:  nil,
			},
			{
				title: "string",
				data:  []byte(`a,b,c`),
				want: [][]string{
					[]string{"a", "b", "c"},
				},
			},
			{
				title: "ints",
				data: []byte(`1,2,3
4,5,6`),
				want: [][]string{
					[]string{"1", "2", "3"},
					[]string{"4", "5", "6"},
				},
			},
			{
				title: "anys",
				data: []byte(`true,1,1.1
false,10,100
`),
				want: [][]string{
					[]string{"true", "1", "1.1"},
					[]string{"false", "10", "100"},
				},
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				translator := cv.NewCSVTranslator(',')
				var v [][]string
				assert.Nil(t, translator.Unmarshal(tc.data, &v))
				assert.Equal(t, tc.want, v)
			})
		}
	})
}
