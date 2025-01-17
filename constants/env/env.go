package env

import "os"

var (
	DB_USER       = os.Getenv("DB_USER")
	DB_PASSWORD   = os.Getenv("DB_PASSWORD")
	DB_URL        = os.Getenv("DB_URL")
	DB_NAME       = os.Getenv("DB_NAME")
	CLIENT_SECRET = os.Getenv("CLIENT_SECRET")
)
