package enums

type PermissionType uint

const (
	View PermissionType = iota + 1
	Edit
	Delete
)

func GetAllPermissionTypes() []PermissionType {
	return []PermissionType{View, Edit, Delete}
}
