package cv

import "encoding/json"

type JSONTranslator struct{}

var (
	_ Translator = &JSONTranslator{}
)

func NewJSONTranslator() *JSONTranslator {
	return &JSONTranslator{}
}

func (t *JSONTranslator) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (t *JSONTranslator) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
