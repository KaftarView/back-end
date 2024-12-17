package enums

type PermissionType uint

const (
	ManageUsers PermissionType = iota + 1
	ManageRoles
	CreateEvent
	ManageEvent
	EditEvent
	PublishEvent
	ManageNewsAndBlogs
	ManagePodcasts
	ModerateComments
	ViewReports
	All
)

func (p PermissionType) String() string {
	switch p {
	case ManageUsers:
		return "ManageUsers"
	case ManageRoles:
		return "ManageRoles"
	case CreateEvent:
		return "CreateEvent"
	case ManageEvent:
		return "ManageEvent"
	case EditEvent:
		return "EditEvent"
	case PublishEvent:
		return "PublishEvent"
	case ManageNewsAndBlogs:
		return "ManageNewsAndBlogs"
	case ManagePodcasts:
		return "ManagePodcasts"
	case ModerateComments:
		return "ModerateComments"
	case ViewReports:
		return "ViewReports"
	case All:
		return "All"
	}
	return ""
}

func GetAllPermissionTypes() []PermissionType {
	return []PermissionType{
		ManageUsers,
		ManageRoles,
		CreateEvent,
		ManageEvent,
		EditEvent,
		PublishEvent,
		ManageNewsAndBlogs,
		ModerateComments,
		ViewReports,
		All,
	}
}
