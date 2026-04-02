package auth

import (
	"errors"
	"slices"
)

var (
	Admin   = "admin"
	Manager = "manager"
	Exec    = "exec"
)

func AuthorizeUser(role string, allowedRoles ...string) error {
	if slices.Contains(allowedRoles, role) {
		return nil
	}
	return errors.New("user not authorized")
}
