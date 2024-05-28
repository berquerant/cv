package cv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"reflect"
)

type CSVTranslator struct {
	Delimiter rune
}

var (
	_ Translator = &CSVTranslator{}
)

func NewCSVTranslator(delimiter rune) *CSVTranslator {
	return &CSVTranslator{
		Delimiter: delimiter,
	}
}

var (
	ErrCSVTranslation = errors.New("CSVTranslation")
)

func (t *CSVTranslator) Unmarshal(data []byte, v any) error {
	result, ok := v.(*[][]string)
	if !ok {
		return fmt.Errorf("%w: need *[][]string, %T", ErrCSVTranslation, v)
	}

	buf := bytes.NewBuffer(data)
	r := csv.NewReader(buf)
	r.Comma = t.Delimiter

	rows, err := r.ReadAll()
	if err != nil {
		return errors.Join(ErrCSVTranslation, err)
	}
	*result = rows
	return nil
}

func (t *CSVTranslator) marshalTrySliceToNestedSlice(v any) any {
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Slice {
		return v
	}
	val := reflect.ValueOf(v)

	if val.Len() == 0 {
		return v
	}
	if val.Index(0).Kind() != reflect.Slice {
		return []any{v}
	}
	return v
}

// Marshal [][]any or []any into csv.
func (t *CSVTranslator) Marshal(v any) ([]byte, error) {
	v = t.marshalTrySliceToNestedSlice(v)
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Slice {
		return nil, fmt.Errorf("%w: needs [][]any or []any, %T, %v", ErrCSVTranslation, v, typ.Kind())
	}
	val := reflect.ValueOf(v)

	var (
		buf bytes.Buffer
		w   = csv.NewWriter(&buf)
	)
	w.Comma = t.Delimiter

	for i := 0; i < val.Len(); i++ {
		row := val.Index(i)
		if row.Kind() != reflect.Slice {
			return nil, fmt.Errorf(
				"%w: row %d, %+v but need slice, %v %T",
				ErrCSVTranslation, i, row.Interface(), row.Kind(), row.Interface())
		}

		elems := make([]string, row.Len())
		for j := 0; j < row.Len(); j++ {
			x := row.Index(j)
			if !IsConvertibleToString(x.Interface()) {
				return nil, fmt.Errorf(
					"%w: row %d, col %d, %+v is not proper kind, %v %T",
					ErrCSVTranslation, i, j, x, x.Kind(), x.Interface())
			}
			elems[j] = fmt.Sprint(x.Interface())
		}
		if err := w.Write(elems); err != nil {
			return nil, errors.Join(ErrCSVTranslation, err)
		}
	}

	w.Flush()
	return buf.Bytes(), nil
}
