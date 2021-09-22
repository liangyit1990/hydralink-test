package config

import (
	"os"
	"strings"
)

// IsDevEnvironment checks if APPENV is set to DEV
func IsDevEnvironment() bool {
	return strings.EqualFold(os.Getenv("APPENV"), "dev")
}
