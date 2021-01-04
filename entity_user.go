package main

import "time"

// User entity
type User struct {
	ID        string `grom:"primaryKey"`
	Birth     string
	Name      string
	CreatedAt time.Time
}
