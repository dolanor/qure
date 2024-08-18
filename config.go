package main

import "os"

type config struct {
	hostPort  string
	domain    string
	dbDirPath string
}

func loadConfig() config {
	hostPort := os.ExpandEnv("${QURE_HOST}:${QURE_PORT}")
	if hostPort == ":" {
		hostPort = ":4444"
	}

	domain := os.Getenv("QURE_DOMAIN")

	dbDirPath := os.Getenv("QURE_DB_DIR")

	return config{
		hostPort:  hostPort,
		domain:    domain,
		dbDirPath: dbDirPath,
	}
}
