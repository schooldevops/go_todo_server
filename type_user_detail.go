package main

type userDetail struct {
	ID         string `json:"id"`
	AvatarImg  string `json:"avatar_img,omitempty"`
	Category   string `json:"category,omitempty"`
	Nick       string `json:"nick,omitempty"`
	ModifiedAt string `json:"modified_at,omitempty"`
}
