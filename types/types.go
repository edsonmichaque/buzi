package types

import (
	"fmt"
)

type Metadata struct {
	BaseURL  string `json:"base_url,omitempty" yaml:"base_url,omitempty"`
	Name     string `json:"name,omitempty" yaml:"name,omitempty"`
	FullName string `json:"full_name,omitempty" yaml:"full_name,omitempty"`
	UID      string `json:"uid,omitempty" yaml:"uid,omitempty"`
}

type Manifest struct {
	Metadata   Metadata               `json:"metadata" yaml:"metadata"`
	Operations map[string]Operation   `json:"operations,omitempty" yaml:"operations,omitempty"`
	Types      map[string]*Definition `json:"types,omitempty" yaml:"types,omitempty"`
	Params     map[string]string      `json:"-" yaml:"-"`
}

type Operation struct {
	Name   string        `json:"name" yaml:"name"`
	Http   *Http         `json:"http,omitempty" yaml:"http,omitempty"`
	Input  *Definition   `json:"input,omitempty" yaml:"input,omitempty"`
	Output *Definition   `json:"output,omitempty" yaml:"output,omitempty"`
	Errors []*Definition `json:"errors,omitempty" yaml:"errors,omitempty"`
}

type DefinitionValue struct {
	Type       string                 `json:"type" yaml:"type"`
	Format     string                 `json:"format" yaml:"format"`
	Properties map[string]*Definition `json:"properties,omitempty" yaml:"properties,omitempty"`
	Items      *Definition            `json:"items,omitempty" yaml:"items,omitempty"`
	Error      *Error                 `json:"error,omitempty" yaml:"error,omitempty"`
	Exception  bool                   `json:"exception,omitempty" yaml:"exception,omitempty"`
	Sensitive  bool                   `json:"sensitive,omitempty" yaml:"sensitive,omitempty"`
	Fault      bool                   `json:"fault,omitempty" yaml:"fault,omitempty"`
	Pattern    string                 `json:"pattern,omitempty" yaml:"pattern,omitempty"`
}

type Error struct {
	StatusCode int `json:"status_code" yaml:"status_code"`
}

type Reference struct {
	Ref string `json:"ref" yaml:"ref"`
}

type Definition struct {
	Value *DefinitionValue `json:",inline" yaml:",inline"`
	Ref   *Reference       `json:",inline" yaml:",inline"`
}

type Http struct {
	Method        string        `json:"method" yaml:"method"`
	RequestURI    string        `json:"request_uri" yaml:"request_uri"`
	Authorization Authorization `json:"auth" yaml:"auth"`
}

type Authorization struct {
	Basic *struct {
		User     string
		Password string
	}
	Bearer bool
}

type BasicAuth struct{}

type File struct {
	Path    string
	Content []byte
}

func (f File) String() string {
	return fmt.Sprintf("Path: %s, Content: %s", f.Path, string(f.Content))
}
