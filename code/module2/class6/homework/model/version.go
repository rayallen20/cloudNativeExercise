package model

import "os"

func GetVersion() string {
	return os.Getenv("VERSION")
}
