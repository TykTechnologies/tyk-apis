// +build !ignore_autogenerated

// Code generated by helpgen. DO NOT EDIT.
package open3

import (
	"github.com/getkin/kin-openapi/openapi3"
)

func (a *Contact) Convert() *openapi3.Contact {
	if a == nil {
		return nil
	}
	return &openapi3.Contact{
		Name:  a.Name,
		URL:   a.URL,
		Email: a.Email,
	}
}

func (a *License) Convert() *openapi3.License {
	if a == nil {
		return nil
	}
	return &openapi3.License{
		Name: a.Name,
		URL:  a.URL,
	}
}
