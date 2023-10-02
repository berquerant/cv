package cv

import (
	"errors"
	"fmt"
)

type CV struct {
	Src       Type
	Dst       Type
	Delimiter rune
}

func New(src, dst Type, delimiter rune) *CV {
	return &CV{
		Src:       src,
		Dst:       dst,
		Delimiter: delimiter,
	}
}

var (
	ErrNoTranslator = errors.New("NoTranslator")
	ErrConvert      = errors.New("Convert")
	ErrInvert       = errors.New("Invert")
)

func (c *CV) Translate(data []byte) ([]byte, error) {
	v, err := c.Unmarshal(data)
	if err != nil {
		return nil, fmt.Errorf("%w: from %s", errors.Join(ErrConvert, err), c.Src)
	}
	u, err := c.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("%w: to %s", errors.Join(ErrInvert, err), c.Dst)
	}
	return u, nil
}

func (c *CV) translator() Translator {
	return NewAutoTranslator(c.Src, c.Dst, c.Delimiter)
}

func (c *CV) Unmarshal(data []byte) (any, error) {
	var a any
	if err := c.translator().Unmarshal(data, &a); err != nil {
		return c.unmarshalAdditional(data)
	}
	return a, nil
}

func (c *CV) unmarshalAdditional(data []byte) (any, error) {
	var a [][]string
	if err := c.translator().Unmarshal(data, &a); err != nil {
		return nil, err
	}
	return a, nil
}

func (c *CV) Marshal(v any) ([]byte, error) {
	return c.translator().Marshal(v)
}

type AutoTranslator struct {
	Src       Type
	Dst       Type
	Delimiter rune
}

func NewAutoTranslator(src, dst Type, delimiter rune) *AutoTranslator {
	return &AutoTranslator{
		Src:       src,
		Dst:       dst,
		Delimiter: delimiter,
	}
}

func (t *AutoTranslator) Marshal(v any) ([]byte, error) {
	translator, ok := t.Dst.Translator(WithDelimiter(t.Delimiter))
	if !ok {
		translator = NewJSONTranslator()
	}

	return translator.Marshal(v)
}

func (t *AutoTranslator) Unmarshal(data []byte, a any) error {
	var unmarshaler Unmarshaler
	if translator, ok := t.Src.Translator(WithDelimiter(t.Delimiter)); ok {
		unmarshaler = translator
	} else {
		unmarshaler = NewAutoUnmarshaler(t.Delimiter)
	}

	return unmarshaler.Unmarshal(data, a)
}

type AutoUnmarshaler struct {
	Delimiter rune
}

func NewAutoUnmarshaler(delimiter rune) *AutoUnmarshaler {
	return &AutoUnmarshaler{
		Delimiter: delimiter,
	}
}

func (a *AutoUnmarshaler) Unmarshal(data []byte, v any) error {
	for _, typ := range ListTypes() {
		translator, _ := typ.Translator(WithDelimiter(a.Delimiter))
		if err := translator.Unmarshal(data, v); err == nil {
			return nil
		}
	}
	return ErrNoTranslator
}
