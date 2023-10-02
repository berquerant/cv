package cv

import "gopkg.in/yaml.v3"

type YAMLTranslator struct{}

func NewYAMLTranslator() *YAMLTranslator {
	return &YAMLTranslator{}
}

func (t *YAMLTranslator) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (t *YAMLTranslator) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}
