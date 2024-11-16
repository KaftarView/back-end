package seed

import (
	"first-project/src/bootstrap"
	"first-project/src/enums"
	repository_database "first-project/src/repository/database"

	"golang.org/x/crypto/bcrypt"
)

type RoleSeeder struct {
	userRepository *repository_database.UserRepository
	admin          *bootstrap.UserInfo
	moderator      *bootstrap.UserInfo
}

func NewRoleSeeder(userRepository *repository_database.UserRepository, admin, moderator *bootstrap.UserInfo) *RoleSeeder {
	return &RoleSeeder{
		userRepository: userRepository,
		admin:          admin,
		moderator:      moderator,
	}
}

func (rs *RoleSeeder) SeedRoles() {
	permissionTypes := enums.GetAllPermissionTypes()
	for _, permissionType := range permissionTypes {
		_, roleExist := rs.userRepository.FindPermissionByType(permissionType)
		if roleExist {
			continue
		}
		rs.userRepository.CreateNewPermission(permissionType)
	}

	roleTypes := enums.GetAllRoleTypes()
	for _, roleType := range roleTypes {
		_, roleExist := rs.userRepository.FindRoleByType(roleType)
		if roleExist {
			continue
		}
		rs.userRepository.CreateNewRole(roleType)
	}

	_, superAdminExist := rs.userRepository.FindActiveOrVerifiedUserByUsername("Admin")
	if !superAdminExist {
		bytes, err := bcrypt.GenerateFromPassword([]byte(rs.admin.Password), 14)
		if err != nil {
			panic(err)
		}
		superAdminUser := rs.userRepository.CreateNewUser("Admin", rs.admin.EmailAddress, string(bytes), "", true)
		superAdminRole, _ := rs.userRepository.FindRoleByType(enums.SuperAdmin)
		rs.userRepository.AssignRoleToUser(superAdminUser, superAdminRole)

		permission, _ := rs.userRepository.FindPermissionByType(enums.ManageUsers)
		rs.userRepository.AssignPermissionToRole(superAdminRole, permission)
		permission, _ = rs.userRepository.FindPermissionByType(enums.ManageRoles)
		rs.userRepository.AssignPermissionToRole(superAdminRole, permission)
		permission, _ = rs.userRepository.FindPermissionByType(enums.CreateEvent)
		rs.userRepository.AssignPermissionToRole(superAdminRole, permission)
		permission, _ = rs.userRepository.FindPermissionByType(enums.EditEvent)
		rs.userRepository.AssignPermissionToRole(superAdminRole, permission)
		permission, _ = rs.userRepository.FindPermissionByType(enums.PublishEvent)
		rs.userRepository.AssignPermissionToRole(superAdminRole, permission)
		permission, _ = rs.userRepository.FindPermissionByType(enums.ManageNewsAndBlogs)
		rs.userRepository.AssignPermissionToRole(superAdminRole, permission)
		permission, _ = rs.userRepository.FindPermissionByType(enums.ModerateComments)
		rs.userRepository.AssignPermissionToRole(superAdminRole, permission)
		permission, _ = rs.userRepository.FindPermissionByType(enums.ViewReports)
		rs.userRepository.AssignPermissionToRole(superAdminRole, permission)
	}
}
