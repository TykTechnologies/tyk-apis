package open3

import "github.com/getkin/kin-openapi/openapi3"

//+o:convert
type Contact struct {
	Name  string `marker:",optional" json:"name,omitempty"`
	URL   string `marker:"url,optional" json:"url,omitempty"`
	Email string `marker:",optional" json:"email,omitempty"`
}

//+o:convert
type License struct {
	Name string `json:"name" yaml:"name"` // Required
	URL  string `marker:"url,optional" json:"url,omitempty" yaml:"url,omitempty"`
}

type Tag struct {
	Name         string            `json:"name,omitempty" `
	Description  string            `json:"description,omitempty"`
	ExternalDocs map[string]string `json:"externalDocs,omitempty"`
}

func (t Tag) Convert() *openapi3.Tag {
	n := &openapi3.Tag{Name: t.Name, Description: t.Description}
	if len(t.ExternalDocs) > 0 {
		n.ExternalDocs = &openapi3.ExternalDocs{
			Description: t.ExternalDocs["description"],
			URL:         t.ExternalDocs["url"],
		}
	}
	return n
}
