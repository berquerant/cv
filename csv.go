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
		return fmt.Errorf("%w: need *[][]string", ErrCSVTranslation)
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

func (t *CSVTranslator) Marshal(v any) ([]byte, error) {
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Slice {
		return nil, fmt.Errorf("%w: not a slice, %v", ErrCSVTranslation, typ.Kind())
	}
	val := reflect.ValueOf(v)

	var (
		buf bytes.Buffer
		w   = csv.NewWriter(&buf)
	)
	w.Comma = t.Delimiter

	for i := 0; i < val.Len(); i++ {
		row := val.Index(i)
		if row.Elem().Kind() != reflect.Slice {
			return nil, fmt.Errorf(
				"%w: row %d, %+v is not proper kind, %v",
				ErrCSVTranslation, i, row.Interface(), row.Kind())
		}

		elems := make([]string, row.Elem().Len())
		for j := 0; j < row.Elem().Len(); j++ {
			x := row.Elem().Index(j)
			switch x.Interface().(type) {
			case bool, int, int8, int16, int32, int64,
				uint, uint8, uint16, uint32, uint64,
				float32, float64, string:
				elems[j] = fmt.Sprint(x.Interface())
			default:
				return nil, fmt.Errorf(
					"%w: row %d, col %d, %+v is not proper kind, %v",
					ErrCSVTranslation, i, j, x, x.Kind())
			}
		}
		if err := w.Write(elems); err != nil {
			return nil, errors.Join(ErrCSVTranslation, err)
		}
	}

	w.Flush()
	return buf.Bytes(), nil
}
