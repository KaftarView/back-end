package enums

type PermissionType uint

const (
	All PermissionType = iota + 1
	ManageUsers
	ManageRoles
	CreateEvent
	ManageEvent
	ReviewEvent
	ManageNews
	ManageJournal
	ManagePodcasts
	ModerateComments
	ViewReports
	CustomerSupport
)

func (p PermissionType) String() string {
	switch p {
	case All:
		return "All"
	case ManageUsers:
		return "ManageUsers"
	case ManageRoles:
		return "ManageRoles"
	case CreateEvent:
		return "CreateEvent"
	case ManageEvent:
		return "ManageEvent"
	case ReviewEvent:
		return "ReviewEvent"
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
	case CustomerSupport:
		return "CustomerSupport"
	}
	return ""
}

var permissionDescriptions = map[PermissionType]string{
	All:              "Grants all permissions.",
	ManageUsers:      "Allows managing users in the system.",
	ManageRoles:      "Allows managing roles and their permissions.",
	CreateEvent:      "Allows creating new events.",
	ManageEvent:      "Allows managing event details.",
	ReviewEvent:      "Allows reviewing and approving events.",
	ManageNews:       "Allows managing and creating news content.",
	ManageJournal:    "Allows managing and creating journal entries.",
	ManagePodcasts:   "Allows managing and creating podcasts.",
	ModerateComments: "Allows moderating comments.",
	ViewReports:      "Allows viewing system reports.",
	CustomerSupport:  "Allows communication with users",
}

func (p PermissionType) Description() string {
	if description, ok := permissionDescriptions[p]; ok {
		return description
	}
	return ""
}

func GetAllPermissionTypes() []PermissionType {
	return []PermissionType{
		All,
		ManageUsers,
		ManageRoles,
		CreateEvent,
		ManageEvent,
		ReviewEvent,
		ManageNews,
		ManageJournal,
		ModerateComments,
		ViewReports,
		CustomerSupport,
	}
}
