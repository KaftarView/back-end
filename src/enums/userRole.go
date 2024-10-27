package enums

type RoleType uint

const (
	User RoleType = iota + 1
	Moderator
	Admin
)

func GetAllRoleTypes() []RoleType {
	return []RoleType{User, Moderator, Admin}
}
