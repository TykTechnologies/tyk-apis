package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"io"
	"os"
	"sort"
	"strings"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

const (
	markerName = "o:convert"
	header     = `
// +build !ignore_autogenerated
%[2]s
// Code generated by helpgen. DO NOT EDIT.
package %[1]s
import (
	"github.com/getkin/kin-openapi/openapi3"
)
`
	funcBlock = `
func (a *%[1]s) Convert() *openapi3.%[1]s {
	if a==nil{
		return nil
	}
	%[2]s
}
`
)

var _ genall.Generator = Generator{}

type Generator struct {
	HeaderFile string `marker:",optional"`
	Year       string `marker:",optional"`
}

func (Generator) RegisterMarkers(reg *markers.Registry) error {
	defn := markers.Must(markers.MakeDefinition(markerName, markers.DescribesType, struct{}{}))
	if err := reg.Register(defn); err != nil {
		return err
	}
	return nil
}

func (g Generator) Generate(ctx *genall.GenerationContext) error {
	var headerText string
	if g.HeaderFile != "" {
		headerBytes, err := ctx.ReadFile(g.HeaderFile)
		if err != nil {
			return err
		}
		headerText = string(headerBytes)
	}
	headerText = strings.ReplaceAll(headerText, " YEAR", g.Year)
	for _, root := range ctx.Roots {
		byType := make(map[string]string)
		if err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
			if m := info.Markers.Get(markerName); m == nil {
				return
			}
			outContent := new(bytes.Buffer)
			fmt.Fprintf(outContent, "return &openapi3.%s {\n", info.Name)
			for _, f := range info.Fields {
				if IsStruct(&f) {
					fmt.Fprintf(outContent, "	%s: a.%s.Convert(),\n", f.Name, f.Name)
				} else {
					fmt.Fprintf(outContent, "	%s: a.%s,\n", f.Name, f.Name)
				}
			}
			fmt.Fprint(outContent, "}")
			b := fmt.Sprintf(funcBlock, info.Name, outContent.String())
			byType[info.Name] = b
		}); err != nil {
			return err
		}

		if len(byType) == 0 {
			continue
		}

		// ensure a stable output order
		typeNames := make([]string, 0, len(byType))
		for typ := range byType {
			typeNames = append(typeNames, typ)
		}
		sort.Strings(typeNames)

		outContent := new(bytes.Buffer)
		fmt.Fprintf(outContent, header, root.Name, headerText)

		for _, typ := range typeNames {
			fmt.Fprintln(outContent, byType[typ])
		}

		outBytes := outContent.Bytes()
		if formatted, err := format.Source(outBytes); err != nil {
			root.AddError(err)
		} else {
			outBytes = formatted
		}

		outputFile, err := ctx.Open(root, "zz_generated.covert.go")
		if err != nil {
			root.AddError(err)
			continue
		}
		defer outputFile.Close()
		n, err := outputFile.Write(outBytes)
		if err != nil {
			root.AddError(err)
			continue
		}
		if n < len(outBytes) {
			root.AddError(io.ErrShortWrite)
		}
	}
	return nil
}

func IsStruct(f *markers.FieldInfo) bool {
	switch f.RawField.Type.(type) {
	case *ast.StarExpr:
		return true
	default:
		return false
	}
}

var (
	registry = &markers.Registry{}
)

func init() {
	genMarker := markers.Must(markers.MakeDefinition("convert", markers.DescribesPackage, Generator{}))
	if err := registry.Register(genMarker); err != nil {
		panic(err)
	}
	ruleMarker := markers.Must(markers.MakeDefinition("output:dir", markers.DescribesPackage, genall.OutputToDirectory("config")))
	if err := registry.Register(ruleMarker); err != nil {
		panic(err)
	}
	if err := genall.RegisterOptionsMarkers(registry); err != nil {
		panic(err)
	}
}

func main() {
	runtime, err := genall.FromOptions(registry, os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if hadErrs := runtime.Run(); hadErrs {
		fmt.Fprintln(os.Stderr, "not entirely successful")
		os.Exit(1)
	}
}