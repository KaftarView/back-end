package seed

import (
	"first-project/src/bootstrap"
	"first-project/src/enums"
	repository_database_interfaces "first-project/src/repository/database/interfaces"

	"golang.org/x/crypto/bcrypt"
)

type RoleSeeder struct {
	userRepository repository_database_interfaces.UserRepository
	admin          *bootstrap.UserInfo
}

func NewRoleSeeder(userRepository repository_database_interfaces.UserRepository, admin *bootstrap.UserInfo) *RoleSeeder {
	return &RoleSeeder{
		userRepository: userRepository,
		admin:          admin,
	}
}

func (rs *RoleSeeder) SeedRoles() {
	permissionTypes := enums.GetAllPermissionTypes()
	for _, permissionType := range permissionTypes {
		_, permissionExist := rs.userRepository.FindPermissionByType(permissionType)
		if permissionExist {
			continue
		}
		rs.userRepository.CreateNewPermission(permissionType)
	}

	roleTypes := enums.GetAllRoleTypes()
	for _, roleType := range roleTypes {
		_, roleExist := rs.userRepository.FindRoleByType(roleType.String())
		if roleExist {
			continue
		}
		rs.userRepository.CreateNewRole(roleType.String())
	}

	_, superAdminExist := rs.userRepository.FindActiveOrVerifiedUserByUsername("Admin")
	if !superAdminExist {
		bytes, err := bcrypt.GenerateFromPassword([]byte(rs.admin.Password), 14)
		if err != nil {
			panic(err)
		}
		superAdminUser := rs.userRepository.CreateNewUser("Admin", rs.admin.EmailAddress, string(bytes), "", true)
		superAdminRole, _ := rs.userRepository.FindRoleByType(enums.SuperAdmin.String())
		rs.userRepository.AssignRoleToUser(superAdminUser, superAdminRole)

		permission, _ := rs.userRepository.FindPermissionByType(enums.All)
		rs.userRepository.AssignPermissionToRole(superAdminRole, permission)
	}
}
