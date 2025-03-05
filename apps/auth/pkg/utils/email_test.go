package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/owjoel/client-factpack/apps/auth/config"
)

// Test: IsAllowedDomain (Success Case)
func TestIsAllowedDomain_Success(t *testing.T) {
	// Mock allowed domains
	config.AllowedDomains = []string{"example.com", "allowed.com"}

	t.Run("Allowed Domain", func(t *testing.T) {
		result := IsAllowedDomain("example.com")
		assert.True(t, result, "Expected domain to be allowed")
	})
}

// Test: IsAllowedDomain (Failure Case)
func TestIsAllowedDomain_Failure(t *testing.T) {
	// Mock allowed domains
	config.AllowedDomains = []string{"example.com", "allowed.com"}

	t.Run("Blocked Domain", func(t *testing.T) {
		result := IsAllowedDomain("blocked.com")
		assert.False(t, result, "Expected domain to be blocked")
	})
}
