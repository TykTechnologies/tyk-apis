package open3

import (
	"github.com/getkin/kin-openapi/openapi3"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

var All []*markers.Definition

func init() {
	All = append(All,
		markers.Must(
			markers.MakeDefinition("o:title", markers.DescribesPackage, ""),
		),
		markers.Must(
			markers.MakeDefinition("o:description", markers.DescribesPackage, ""),
		),
		markers.Must(
			markers.MakeDefinition("o:termsOfService", markers.DescribesPackage, ""),
		),
		markers.Must(
			markers.MakeDefinition("o:contact", markers.DescribesPackage, Contact{}),
		),
		markers.Must(
			markers.MakeDefinition("o:license", markers.DescribesPackage, License{}),
		),
		markers.Must(
			markers.MakeDefinition("o:version", markers.DescribesPackage, ""),
		),
		markers.Must(
			markers.MakeDefinition("o:model", markers.DescribesType, struct{}{}),
		),
	)
}

func IsModel(info *markers.TypeInfo) bool {
	return info.Markers.Get("o:model") != nil
}

var marks = []string{
	"o:title", "o:description", "o:termsOfService", "o:contact", "o:license", "o:version",
}

func Load(swagg *openapi3.Swagger, coll *markers.Collector, pkg *loader.Package) error {
	m, err := markers.PackageMarkers(coll, pkg)
	if err != nil {
		return err
	}
	for _, v := range marks {
		switch v {
		case "o:title", "o:description", "o:termsOfService", "o:contact", "o:license", "o:version":
			if err := loadInfo(m, swagg, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func loadInfo(m markers.MarkerValues, swagg *openapi3.Swagger, key string) error {
	if swagg.Info == nil {
		swagg.Info = &openapi3.Info{}
	}
	switch key {
	case "o:title":
		swagg.Info.Title = getString(m, key)
	case "o:description":
		swagg.Info.Description = getString(m, key)
	case "o:termsOfService":
		swagg.Info.TermsOfService = getString(m, key)
	case "o:contact":
		if x := m.Get(key); x != nil {
			o := x.(Contact)
			swagg.Info.Contact = o.Convert()
		}
	case "o:license":
		if x := m.Get(key); x != nil {
			o := x.(License)
			swagg.Info.License = o.Convert()
		}
	case "o:version":
		swagg.Info.Version = getString(m, key)
	}
	return nil
}

func getString(m markers.MarkerValues, key string) string {
	if v := m.Get(key); v != nil {
		s, _ := v.(string)
		return s
	}
	return ""
}
