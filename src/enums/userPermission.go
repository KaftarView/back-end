package enums

type PermissionType uint

const (
	ManageUsers PermissionType = iota + 1
	ManageRoles
	CreateEvent
	ManageEvent
	EditEvent
	PublishEvent
	ManageNews
	ManageJournal
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
	case ManageNews:
		return "ManageNews"
	case ManageJournal:
		return "ManageJournal"
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
		ManageNews,
		ManageJournal,
		ModerateComments,
		ViewReports,
		All,
	}
}
