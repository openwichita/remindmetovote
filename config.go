package main

import "os"

// Config holds the application's configuration parameters.
type Config struct {
	DB   string
	Bind string
}

// NewConfig creates a Config from the environment.
func NewConfig() (Config, error) {
	db := os.Getenv("DATABASE_URL")
	if db == "" {
		// Fall back to local defaults
		db = "postgres://root@localhost/remindmetovote?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		// Fall back to local defaults
		port = "3927"
	}
	port = ":" + port

	return Config{
		DB:   db,
		Bind: port,
	}, nil
}
