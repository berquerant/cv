package cv

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
	default:
		return "unknown"
	}
}

func TypeFromString(s string) Type {
	switch s {
	case "json", "j":
		return Tjson
	case "yaml", "yml", "y":
		return Tyaml
	case "toml", "t":
		return Ttoml
	case "csv", "c":
		return Tcsv
	default:
		return Tunknown
	}
}

func ListTypes() []Type {
	return []Type{
		Tjson,
		Tyaml,
		Ttoml,
		Tcsv,
	}
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
	default:
		return nil, false
	}
}
