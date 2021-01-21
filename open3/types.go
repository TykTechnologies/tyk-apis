package open3

//+o:convert
type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

//+o:convert
type License struct {
	Name string `json:"name" yaml:"name"` // Required
	URL  string `marker:"url,omitempty" json:"url,omitempty" yaml:"url,omitempty"`
}
