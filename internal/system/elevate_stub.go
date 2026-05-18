//go:build !windows

package system

func isAdmin() bool {
	return false
}

func ensureAdmin() bool {
	return false
}
