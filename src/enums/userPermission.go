package enums

type PermissionType uint

const (
	ManageUsers PermissionType = iota + 1
	ManageRoles
	CreateEvent
	EditEvent
	PublishEvent
	ManageNewsAndBlogs
	ModerateComments
	ViewReports
)

func GetAllPermissionTypes() []PermissionType {
	return []PermissionType{
		ManageUsers,
		ManageRoles,
		CreateEvent,
		EditEvent,
		PublishEvent,
		ManageNewsAndBlogs,
		ModerateComments,
		ViewReports,
	}
}
