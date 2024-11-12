package enums

type RoleType uint

const (
	SuperAdmin RoleType = iota + 1
	EventManager
	ContentManager
	Editor
	Moderator
	Viewer
	User
)

func GetAllRoleTypes() []RoleType {
	return []RoleType{
		SuperAdmin,
		EventManager,
		ContentManager,
		Editor,
		Moderator,
		Viewer,
		User,
	}
}
