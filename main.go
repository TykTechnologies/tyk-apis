package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TykTechnologies/tyk-apis/open3"
	"github.com/getkin/kin-openapi/openapi3"
	"sigs.k8s.io/controller-tools/pkg/crd"
	"sigs.k8s.io/controller-tools/pkg/loader"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

const (
	id     = "http://jsonschema.net"
	scheme = "http://json-schema.org/draft-04/schema"
)

var (
	registry = &markers.Registry{}
)

func init() {
	genMarker := markers.Must(markers.MakeDefinition("tyk", markers.DescribesPackage, Generator{}))
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

//go:generate go run convert/main.go convert  paths=./open3/

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

var _ genall.Generator = Generator{}

// Generator generates api spec.
type Generator struct {
	Targets []string `json:"" marker:""`
}

func (g Generator) c() crd.Generator {
	return crd.Generator{}
}
func (g Generator) CheckFilter() loader.NodeFilter {
	return g.c().CheckFilter()
}

func (g Generator) RegisterMarkers(into *markers.Registry) error {
	for _, m := range open3.All {
		err := into.Register(m)
		if err != nil {
			return err
		}
	}
	return g.c().RegisterMarkers(into)
}
func (g Generator) Generate(ctx *genall.GenerationContext) error {
	parser := &crd.Parser{
		Collector:           ctx.Collector,
		Checker:             ctx.Checker,
		AllowDangerousTypes: true,
	}
	for _, root := range ctx.Roots {
		parser.NeedPackage(root)
	}

	for _, root := range ctx.Roots {
		markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {

		})
	}

	swagg := &openapi3.Swagger{}
	for _, root := range ctx.Roots {
		err := open3.Load(swagg, ctx.Collector, root)
		if err != nil {
			return err
		}
	}
	// add components
	return writeJSON(ctx, "schema.json", swagg)
}

func writeJSON(g *genall.GenerationContext, itemPath string, object *openapi3.Swagger) error {
	out, err := g.Open(nil, itemPath)
	if err != nil {
		return err
	}
	defer out.Close()
	e := json.NewEncoder(out)
	e.SetIndent("", "  ")
	return e.Encode(object)
}
