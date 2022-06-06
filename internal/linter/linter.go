package linter

type Config struct {
	Info  Info  `json:"info,omitempty" yaml:"info,omitempty"`
	Paths Paths `json:"operation,omitempty" yaml:"operation,omitempty"`
}

type Info struct {
	License       *License `json:"license,omitempty" yaml:"license,omitempty"`
	Description   *Text    `json:"description,omitempty" yaml:"description,omitempty"`
	TermOfService *Text    `json:"term_of_service,omitempty" yaml:"term_of_service,omitempty"`
	Title         *Text    `json:"title,omitempty" yaml:"title,omitempty"`
}

type Paths struct {
	Summary     Text      `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description Text      `json:"description,omitempty" yaml:"description,omitempty"`
	Operation   Operation `json:"operation,omitempty" yaml:"operation,omitempty"`
}

type License struct {
	Required  bool     `json:"required,omitempty" yaml:"required,omitempty"`
	AllowList []string `json:"allow_list,omitempty" yaml:"allow_list,omitempty"`
	DenyList  []string `json:"deny_list,omitempty" yaml:"deny_list,omitempty"`
}

type Operation struct {
	OperationID Text `json:"operation_id,omitempty" yaml:"operation_id,omitempty"`
	Tags        Tags `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary     Text `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description Text `json:"description,omitempty" yaml:"description,omitempty"`
}

type Required struct {
	Required bool `json:"required,omitempty" yaml:"required,omitempty"`
}

type Text struct {
	Required  bool     `json:"required,omitempty" yaml:"required,omitempty"`
	Pattern   string   `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MinLength int      `json:"min_length,omitempty" yaml:"min_length,omitempty"`
	MaxLength int      `json:"max_length,omitempty" yaml:"max_length,omitempty"`
	DenyList  []string `json:"deny_list,omitempty" yaml:"deny_list,omitempty"`
}

type Tags struct {
	Text  `yaml:",inline"`
	Array `yaml:",inline"`
}

type Array struct {
	MinItems int `json:"min_items,omitempty" yaml:"min_items,omitempty"`
	MaxItems int `json:"max_items,omitempty" yaml:"max_items,omitempty"`
}

type Rule struct {
	Name          string         `json:"name,omitempty" yaml:"name,omitempty"`
	Object        string         `json:"object,omitempty" yaml:"objejct,omitempty"`
	Description   string         `json:"description,omitempty" yaml:"description,omitempty"`
	Contains      *Containts     `json:"containts,omitempty" yaml:"containts,omitempty"`
	DoesntContain *DoesntContain `json:"doesntContain,omitempty" yaml:"doesntContain,omitempty"`
	MinLength     int            `json:"min_length,omitempty" yaml:"min_length,omitempty"`
	MaxLength     int            `json:"max_length,omitempty" yaml:"max_length,omitempty"`
}

type Containts struct {
}

type DoesntContain struct {
}
