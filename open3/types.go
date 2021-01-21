package open3

//+o:convert
type Contact struct {
	Name  string `marker:",optional" json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `marker:"url,optional" json:"url,omitempty" yaml:"url,omitempty"`
	Email string `marker:",optional" json:"email,omitempty" yaml:"email,omitempty"`
}

//+o:convert
type License struct {
	Name string `json:"name" yaml:"name"` // Required
	URL  string `marker:"url,optional" json:"url,omitempty" yaml:"url,omitempty"`
}
