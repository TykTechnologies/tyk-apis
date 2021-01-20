package test

type API struct {
	Name  string `json:"name"`
	Embed Embed  `json:"embed"`
}

type Embed struct {
	Value int `json:"value"`
}
