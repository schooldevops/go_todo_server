package main

type user struct {
	ID        string `json:"id"`
	Birth     string `json:"birth,omitempty"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at,omitempty"`
}
