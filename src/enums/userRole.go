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

func (r RoleType) String() string {
	switch r {
	case SuperAdmin:
		return "SuperAdmin"
	case EventManager:
		return "EventManager"
	case ContentManager:
		return "ContentManager"
	case Editor:
		return "Editor"
	case Moderator:
		return "Moderator"
	case Viewer:
		return "Viewer"
	case User:
		return "User"
	}
	return ""
}

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
