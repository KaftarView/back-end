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

	_, adminExist := rs.userRepository.FindActiveOrVerifiedUserByUsername("Admin")
	if !adminExist {
		bytes, err := bcrypt.GenerateFromPassword([]byte(rs.admin.Password), 14)
		if err != nil {
			panic(err)
		}
		adminUser := rs.userRepository.CreateNewUser("Admin", rs.admin.EmailAddress, string(bytes), "", true)
		adminRole, _ := rs.userRepository.FindRoleByType(enums.Admin)
		rs.userRepository.AssignRoleToUser(adminUser, adminRole)
		viewPermission, _ := rs.userRepository.FindPermissionByType(enums.View)
		deletePermission, _ := rs.userRepository.FindPermissionByType(enums.Delete)
		editPermission, _ := rs.userRepository.FindPermissionByType(enums.Edit)
		rs.userRepository.AssignPermissionToRole(adminRole, viewPermission)
		rs.userRepository.AssignPermissionToRole(adminRole, deletePermission)
		rs.userRepository.AssignPermissionToRole(adminRole, editPermission)
	}

	_, moderatorExist := rs.userRepository.FindActiveOrVerifiedUserByUsername("Moderator")
	if !moderatorExist {
		bytes, err := bcrypt.GenerateFromPassword([]byte(rs.moderator.Password), 14)
		if err != nil {
			panic(err)
		}
		moderatorUser := rs.userRepository.CreateNewUser("Moderator", rs.moderator.EmailAddress, string(bytes), "", true)
		moderatorRole, _ := rs.userRepository.FindRoleByType(enums.Moderator)
		rs.userRepository.AssignRoleToUser(moderatorUser, moderatorRole)
		viewPermission, _ := rs.userRepository.FindPermissionByType(enums.View)
		editPermission, _ := rs.userRepository.FindPermissionByType(enums.Edit)
		rs.userRepository.AssignPermissionToRole(moderatorRole, viewPermission)
		rs.userRepository.AssignPermissionToRole(moderatorRole, editPermission)
	}
}
