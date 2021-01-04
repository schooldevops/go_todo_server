package main

type todos struct {
	ID              int64  `json:"id,omitempty"`
	UserID          int64  `json:"user_id,omitempty"`
	TITLE           string `json:"title,omitempty"`
	Priority        string `json:"priority,omitempty"`
	Status          string `json:"status,omitempty"`
	CompletionLevel int64  `json:"completion_level,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	ModifiedAt      string `json:"modified_at,omitempty"`
	DoneAt          string `json:"done_at,omitempty"`
}
