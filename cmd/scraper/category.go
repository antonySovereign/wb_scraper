package main

type Category struct {
	ID          int        `json:"id"`
	Parent      int        `json:"parent,omitempty"`
	Name        string     `json:"name"`
	Url         string     `json:"url"`
	Shard       string     `json:"shard,omitempty"`
	Query       string     `json:"query,omitempty"`
	SearchQuery string     `json:"searchQuery,omitempty"`
	Childs      []Category `json:"childs,omitempty"`
}
