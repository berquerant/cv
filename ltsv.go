package cv

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type LTSVTranslator struct {
	Delimiter rune
}

func NewLTSVTranslator(delimiter rune) *LTSVTranslator {
	return &LTSVTranslator{
		Delimiter: delimiter,
	}
}

var (
	ErrLTSVTranslation = errors.New("LTSVTranslation")
)

func (t *LTSVTranslator) Unmarshal(data []byte, v any) error {
	result, ok := v.(*[]map[string]string)
	if !ok {
		return fmt.Errorf("%w: need *[]map[string]string", ErrLTSVTranslation)
	}

	var (
		buf = bytes.NewBuffer(data)
		sc  = bufio.NewScanner(buf)
		sep = string(t.Delimiter)
		acc []map[string]string
	)

	for sc.Scan() {
		line := sc.Text()
		dict := map[string]string{}
		for _, elem := range strings.Split(line, "\t") {
			xs := strings.SplitN(elem, sep, 2)
			if len(xs) == 2 {
				dict[xs[0]] = xs[1]
			}
		}
		acc = append(acc, dict)
	}
	if err := sc.Err(); err != nil {
		return errors.Join(ErrLTSVTranslation, err)
	}

	*result = acc
	return nil
}

func (t *LTSVTranslator) Marshal(v any) ([]byte, error) {
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Slice {
		return nil, fmt.Errorf("%w: not a slice, %v", ErrLTSVTranslation, typ.Kind())
	}
	val := reflect.ValueOf(v)

	var (
		sep = string(t.Delimiter)
		buf bytes.Buffer
		w   = csv.NewWriter(&buf)
	)
	w.Comma = '\t'

	for i := 0; i < val.Len(); i++ {
		row := val.Index(i)
		rowDict, ok := row.Interface().(map[string]any)
		if !ok {
			return nil, fmt.Errorf(
				"%w: row %d, %+v is not proper kind",
				ErrLTSVTranslation, i, row.Interface())
		}

		var (
			elems = make([]string, len(rowDict))
			j     int
		)
		for k, v := range rowDict {
			switch v.(type) {
			case bool, int, int8, int16, int32, int64,
				uint, uint8, uint16, uint32, uint64,
				float32, float64, string:
				elems[j] = fmt.Sprintf("%s%s%v", k, sep, v)
			default:
				return nil, fmt.Errorf(
					"%w: row %d, col %d, value %+v is not proper kind",
					ErrLTSVTranslation, i, j, v)
			}
			j++
		}

		sort.Strings(elems)
		if err := w.Write(elems); err != nil {
			return nil, errors.Join(ErrLTSVTranslation, err)
		}
	}

	w.Flush()
	return buf.Bytes(), nil
}
