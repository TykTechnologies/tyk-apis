package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"sigs.k8s.io/controller-tools/pkg/crd"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/yaml"

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
	for k := range parser.Types {
		for _, x := range g.Targets {
			name := strings.ToLower(k.Name)
			if name == strings.ToLower(x) {
				parser.NeedFlattenedSchemaFor(k)
				fullSchema := parser.FlattenedSchemata[k]
				schema := fullSchema.DeepCopy()
				schema.ID = id
				schema.Schema = scheme
				writeJSON(ctx, name+".json",
					schema)
			}
		}
	}
	return nil
}

func writeJSON(g *genall.GenerationContext, itemPath string, object proto.Message) error {
	out, err := g.Open(nil, itemPath)
	if err != nil {
		return err
	}
	defer out.Close()
	b, err := yaml.Marshal(object)
	if err != nil {
		return err
	}
	o := map[string]interface{}{}
	err = yaml.Unmarshal(b, &o)
	e := json.NewEncoder(out)
	e.SetIndent("", "  ")
	return e.Encode(object)

}
