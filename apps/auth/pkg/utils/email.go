package utils

import "github.com/owjoel/client-factpack/apps/auth/config"

// Checks email domain against allowed list, set through app environment variables
func IsAllowedDomain(domain string) bool {
	for _, d := range config.AllowedDoamins {
		if domain == d {
			return true
		}
	}
	return false
}