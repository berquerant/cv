// Code generated by "goconfig -field Delimiter rune -option -configOption Option -output interface_config_generated.go"; DO NOT EDIT.

package cv

type ConfigItem[T any] struct {
	modified     bool
	value        T
	defaultValue T
}

func (s *ConfigItem[T]) Set(value T) {
	s.modified = true
	s.value = value
}
func (s *ConfigItem[T]) Get() T {
	if s.modified {
		return s.value
	}
	return s.defaultValue
}
func (s *ConfigItem[T]) Default() T {
	return s.defaultValue
}
func (s *ConfigItem[T]) IsModified() bool {
	return s.modified
}
func NewConfigItem[T any](defaultValue T) *ConfigItem[T] {
	return &ConfigItem[T]{
		defaultValue: defaultValue,
	}
}

type Config struct {
	Delimiter *ConfigItem[rune]
}
type ConfigBuilder struct {
	delimiter rune
}

func (s *ConfigBuilder) Delimiter(v rune) *ConfigBuilder {
	s.delimiter = v
	return s
}
func (s *ConfigBuilder) Build() *Config {
	return &Config{
		Delimiter: NewConfigItem(s.delimiter),
	}
}

func NewConfigBuilder() *ConfigBuilder { return &ConfigBuilder{} }
func (s *Config) Apply(opt ...Option) {
	for _, x := range opt {
		x(s)
	}
}

type Option func(*Config)

func WithDelimiter(v rune) Option {
	return func(c *Config) {
		c.Delimiter.Set(v)
	}
}
