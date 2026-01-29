package auth

import (
	"slices"

	"github.com/goaperture/goaperture/v2/exception"
)

type Permission string
type Permissions []Permission

func (p *Permissions) Check(key string) bool {
	return slices.Contains(*p, Permission(key))
}

func (p *Permissions) CheckX(key string) {
	if !p.Check(key) {
		exception.NotAccess(key)
	}
}
