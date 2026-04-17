package admin

// Permissions represents a set of permissions that a user or role has.
type Permissions []string

func (p Permissions) Has(permission string) bool {
	for _, perm := range p {
		if perm == "*" || perm == permission {
			return true
		}
	}

	return false
}
