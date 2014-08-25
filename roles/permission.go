package roles

type PermissionMode uint32

const (
	Read PermissionMode = 1 << (32 - 1 - iota)
	Update
	Create
	Delete
	CURD
)

type Permission struct {
	role       *Role
	allowRoles map[PermissionMode][]string
	denyRoles  map[PermissionMode][]string
}

func hasSameElem(vs1 []string, vs2 []string) bool {
	for _, v1 := range vs1 {
		for _, v2 := range vs2 {
			if v1 == v2 {
				return true
			}
		}
	}
	return false
}

func (permission *Permission) HasPermission(mode PermissionMode, roles ...string) bool {
	if len(permission.denyRoles) != 0 {
		if denyRoles := permission.denyRoles[mode]; denyRoles != nil {
			if hasSameElem(denyRoles, roles) {
				return false
			}
		}
	}

	if len(permission.allowRoles) != 0 {
		if allowRoles := permission.allowRoles[mode]; allowRoles != nil {
			if hasSameElem(allowRoles, roles) {
				return true
			}
		}
	} else if len(permission.denyRoles) != 0 {
		return true
	}

	return false
}

func (permission *Permission) Allow(mode PermissionMode, roles ...string) *Permission {
	if mode == CURD {
		return permission.Allow(Create, roles...).Allow(Update, roles...).Allow(Read, roles...).Allow(Delete, roles...)
	}

	if permission.allowRoles[mode] == nil {
		permission.allowRoles[mode] = []string{}
	}
	permission.allowRoles[mode] = append(permission.allowRoles[mode], roles...)
	return permission
}

func (permission *Permission) Deny(mode PermissionMode, roles ...string) *Permission {
	if mode == CURD {
		return permission.Deny(Create, roles...).Deny(Update, roles...).Deny(Read, roles...).Deny(Delete, roles...)
	}

	if permission.denyRoles[mode] == nil {
		permission.denyRoles[mode] = []string{}
	}
	permission.denyRoles[mode] = append(permission.denyRoles[mode], roles...)
	return permission
}
