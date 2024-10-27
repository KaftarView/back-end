package seed

import (
	"first-project/src/bootstrap"
	"first-project/src/enums"
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
	rolesType := enums.GetAllRoleTypes()
	for _, roleType := range rolesType {
		_, roleExist := rs.userRepository.FindRoleByType(roleType)
		if roleExist {
			continue
		}
		rs.userRepository.CreateNewRole(roleType)
	}

	_, adminExist := rs.userRepository.FindByUsernameAndVerified("Admin", true)
	if !adminExist {
		bytes, err := bcrypt.GenerateFromPassword([]byte(rs.admin.Password), 14)
		if err != nil {
			panic(err)
		}
		adminUser := rs.userRepository.CreateNewUser("Admin", rs.admin.EmailAddress, string(bytes), "", true)
		adminRole, _ := rs.userRepository.FindRoleByType(enums.Admin)
		rs.userRepository.AssignRoleToUser(adminUser, adminRole)
	}

	_, moderatorExist := rs.userRepository.FindByUsernameAndVerified("Moderator", true)
	if !moderatorExist {
		bytes, err := bcrypt.GenerateFromPassword([]byte(rs.moderator.Password), 14)
		if err != nil {
			panic(err)
		}
		moderatorUser := rs.userRepository.CreateNewUser("Moderator", rs.moderator.EmailAddress, string(bytes), "", true)
		moderatorRole, _ := rs.userRepository.FindRoleByType(enums.Moderator)
		rs.userRepository.AssignRoleToUser(moderatorUser, moderatorRole)
	}
}
