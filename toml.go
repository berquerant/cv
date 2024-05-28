package cv

import "github.com/pelletier/go-toml/v2"

type TOMLTranslator struct{}

var (
	_ Translator = &TOMLTranslator{}
)

func NewTOMLTranslator() *TOMLTranslator {
	return &TOMLTranslator{}
}

func (t *TOMLTranslator) Marshal(v any) ([]byte, error) {
	return toml.Marshal(v)
}

func (t *TOMLTranslator) Unmarshal(data []byte, v any) error {
	return toml.Unmarshal(data, v)
}
