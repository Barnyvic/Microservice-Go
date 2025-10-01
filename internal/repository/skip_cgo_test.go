//go:build !cgo
// +build !cgo

package repository

import (
	"testing"
)

// Skip repository tests when CGO is not available
// This is a workaround for environments where CGO is disabled
// The repository functionality is tested through integration tests

func TestSkipRepositoryTests_CGORequired(t *testing.T) {
	t.Skip("Repository tests require CGO to be enabled for SQLite driver. " +
		"To run repository tests, ensure CGO_ENABLED=1 and a C compiler is available. " +
		"Repository functionality is validated through service layer tests.")
}
