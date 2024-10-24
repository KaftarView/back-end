package seed

import (
	"first-project/src/bootstrap"
	"first-project/src/repository"

	"golang.org/x/crypto/bcrypt"
)

type RoleSeeder struct {
	userRepository *repository.UserRepository
	admin          *bootstrap.UserInfo
	moderator      *bootstrap.UserInfo
}

func NewRoleSeeder(userRepository *repository.UserRepository, admin, moderator *bootstrap.UserInfo) *RoleSeeder {
	return &RoleSeeder{
		userRepository: userRepository,
		admin:          admin,
		moderator:      moderator,
	}
}

func (rs *RoleSeeder) SeedRoles() {
	roleNames := []string{"Admin", "Moderator", "User"}
	for _, roleName := range roleNames {
		_, roleExist := rs.userRepository.FindRoleByName(roleName)
		if roleExist {
			continue
		}
		rs.userRepository.CreateNewRole(roleName)
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(rs.admin.Password), 14)
	if err != nil {
		panic(err)
	}
	adminUser := rs.userRepository.CreateNewUser("Admin", rs.admin.EmailAddress, string(bytes), "", true)
	adminRole, _ := rs.userRepository.FindRoleByName(roleNames[0])
	rs.userRepository.AssignRoleToUser(adminUser, adminRole)

	bytes, err = bcrypt.GenerateFromPassword([]byte(rs.moderator.Password), 14)
	if err != nil {
		panic(err)
	}
	moderatorUser := rs.userRepository.CreateNewUser("Moderator", rs.moderator.EmailAddress, string(bytes), "", true)
	moderatorRole, _ := rs.userRepository.FindRoleByName(roleNames[1])
	rs.userRepository.AssignRoleToUser(moderatorUser, moderatorRole)
}