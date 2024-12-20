package enums

type RoleType uint

const (
	SuperAdmin RoleType = iota + 1
	User
)

func (r RoleType) String() string {
	switch r {
	case SuperAdmin:
		return "SuperAdmin"
	case User:
		return "User"
	}
	return ""
}

func GetAllRoleTypes() []RoleType {
	return []RoleType{
		SuperAdmin,
		User,
	}
}
