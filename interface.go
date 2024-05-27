package cv

import "slices"

type Unmarshaler interface {
	Unmarshal(data []byte, v any) error
}

type Marshaler interface {
	Marshal(v any) ([]byte, error)
}

type Translator interface {
	Marshaler
	Unmarshaler
}

type Type int

const (
	Tunknown Type = iota
	Tjson
	Tyaml
	Ttoml
	Tcsv
	Tltsv
)

func (t Type) String() string {
	switch t {
	case Tjson:
		return "json"
	case Tyaml:
		return "yaml"
	case Ttoml:
		return "toml"
	case Tcsv:
		return "csv"
	case Tltsv:
		return "ltsv"
	default:
		return "unknown"
	}
}

type StringTypeTuple struct {
	Keys  []string
	Value Type
}

type StringTypes struct{}

func (StringTypes) Tuples() []*StringTypeTuple {
	return []*StringTypeTuple{
		{
			Keys:  []string{"json", "j"},
			Value: Tjson,
		},
		{
			Keys:  []string{"yaml", "yml", "y"},
			Value: Tyaml,
		},
		{
			Keys:  []string{"toml", "t"},
			Value: Ttoml,
		},
		{
			Keys:  []string{"csv", "c"},
			Value: Tcsv,
		},
		{
			Keys:  []string{"ltsv", "l"},
			Value: Tltsv,
		},
	}
}

func (t StringTypes) Get(s string) Type {
	for _, x := range t.Tuples() {
		if slices.Contains(x.Keys, s) {
			return x.Value
		}
	}
	return Tunknown
}

//go:generate go run github.com/berquerant/goconfig@v0.3.0 -field "Delimiter rune" -option -configOption Option -output interface_config_generated.go

func (t Type) Translator(opt ...Option) (Translator, bool) {
	config := NewConfigBuilder().Delimiter(',').Build()
	config.Apply(opt...)

	switch t {
	case Tjson:
		return NewJSONTranslator(), true
	case Tyaml:
		return NewYAMLTranslator(), true
	case Ttoml:
		return NewTOMLTranslator(), true
	case Tcsv:
		return NewCSVTranslator(config.Delimiter.Get()), true
	case Tltsv:
		return NewLTSVTranslator(':'), true
	default:
		return nil, false
	}
}
