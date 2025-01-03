package seed

import (
	"first-project/src/bootstrap"
	"first-project/src/entities"
	"first-project/src/enums"
	repository_database_interfaces "first-project/src/repository/database/interfaces"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type roleSeeder struct {
	db             *gorm.DB
	userRepository repository_database_interfaces.UserRepository
	superAdmin     *bootstrap.AdminCredentials
}

func NewRoleSeeder(
	db *gorm.DB,
	userRepository repository_database_interfaces.UserRepository,
	superAdmin *bootstrap.AdminCredentials,
) *roleSeeder {
	return &roleSeeder{
		db:             db,
		userRepository: userRepository,
		superAdmin:     superAdmin,
	}
}

func (roleSeeder *roleSeeder) SeedRoles() {
	permissionTypes := enums.GetAllPermissionTypes()
	for _, permissionType := range permissionTypes {
		_, permissionExist := roleSeeder.userRepository.FindPermissionByType(roleSeeder.db, permissionType)
		if permissionExist {
			continue
		}
		roleSeeder.userRepository.CreateNewPermission(roleSeeder.db, permissionType, permissionType.Description())
	}

	rolePermissions := map[enums.RoleType][]enums.PermissionType{
		enums.SuperAdmin: {enums.All},
		enums.User:       {},
	}

	roleTypes := enums.GetAllRoleTypes()
	for _, roleType := range roleTypes {
		_, roleExist := roleSeeder.userRepository.FindRoleByType(roleSeeder.db, roleType.String())
		if roleExist {
			continue
		}
		role := roleSeeder.userRepository.CreateNewRole(roleSeeder.db, roleType.String())
		for _, permissionType := range rolePermissions[roleType] {
			permission, _ := roleSeeder.userRepository.FindPermissionByType(roleSeeder.db, permissionType)
			roleSeeder.userRepository.AssignPermissionToRole(roleSeeder.db, role, permission)
		}
	}

	_, superAdminExist := roleSeeder.userRepository.FindActiveOrVerifiedUserByUsername(roleSeeder.db, roleSeeder.superAdmin.Name)
	if !superAdminExist {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(roleSeeder.superAdmin.Password), 14)
		if err != nil {
			panic(err)
		}
		superAdmin := &entities.User{
			Name:     roleSeeder.superAdmin.Name,
			Email:    roleSeeder.superAdmin.EmailAddress,
			Password: string(hashedPassword),
			Token:    "",
			Verified: true,
		}
		if err := roleSeeder.userRepository.CreateNewUser(roleSeeder.db, superAdmin); err != nil {
			panic(err)
		}
		superAdminRole, _ := roleSeeder.userRepository.FindRoleByType(roleSeeder.db, enums.SuperAdmin.String())
		roleSeeder.userRepository.AssignRoleToUser(roleSeeder.db, superAdmin, superAdminRole)
	}
}
