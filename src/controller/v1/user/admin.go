package controller_v1_user

import (
	"first-project/src/application"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type AdminUserController struct {
	constants   *bootstrap.Constants
	userService *application.UserService
}

func NewAdminUserController(
	constants *bootstrap.Constants,
	userService *application.UserService,
) *AdminUserController {
	return &AdminUserController{
		constants:   constants,
		userService: userService,
	}
}

func (adminUserController *AdminUserController) GetRolesList(c *gin.Context) {
	rolesList := adminUserController.userService.GetRolesList()
	controller.Response(c, 200, "", rolesList)
}

func (adminUserController *AdminUserController) CreateRole(c *gin.Context) {
	type createRolesParams struct {
		Permissions []string `json:"permissions" validate:"required"`
		RoleName    string   `json:"role" validate:"required"`
	}
	param := controller.Validated[createRolesParams](c, &adminUserController.constants.Context)
	role := adminUserController.userService.CreateNewRole(param.RoleName)
	adminUserController.userService.AssignPermissionsToRole(role.ID, param.Permissions)

	trans := controller.GetTranslator(c, adminUserController.constants.Context.Translator)
	message, _ := trans.T("successMessage.createRole")
	controller.Response(c, 200, message, nil)
}

func (adminUserController *AdminUserController) GetRoleOwners(c *gin.Context) {
	type getRoleParams struct {
		RoleID uint `uri:"roleID" validate:"required"`
	}
	param := controller.Validated[getRoleParams](c, &adminUserController.constants.Context)
	roleOwners := adminUserController.userService.GetRoleOwners(param.RoleID)
	controller.Response(c, 200, "", roleOwners)
}

func (adminUserController *AdminUserController) DeleteRole(c *gin.Context) {
	type deleteRoleParams struct {
		RoleID uint `uri:"roleID" validate:"required"`
	}
	param := controller.Validated[deleteRoleParams](c, &adminUserController.constants.Context)
	adminUserController.userService.DeleteRole(param.RoleID)

	trans := controller.GetTranslator(c, adminUserController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteRole")
	controller.Response(c, 200, message, nil)
}

func (adminUserController *AdminUserController) UpdateRole(c *gin.Context) {
	type updateRolesParams struct {
		Permissions []string `json:"permissions" validate:"required"`
		RoleID      uint     `uri:"roleID" validate:"required"`
	}
	param := controller.Validated[updateRolesParams](c, &adminUserController.constants.Context)
	adminUserController.userService.AssignPermissionsToRole(param.RoleID, param.Permissions)

	trans := controller.GetTranslator(c, adminUserController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updateRole")
	controller.Response(c, 200, message, nil)
}

func (adminUserController *AdminUserController) DeleteRolePermission(c *gin.Context) {
	type deleteRolePermissionParams struct {
		RoleID       uint `uri:"roleID" validate:"required"`
		PermissionID uint `uri:"permissionID" validate:"required"`
	}
	param := controller.Validated[deleteRolePermissionParams](c, &adminUserController.constants.Context)
	adminUserController.userService.DeleteRolePermission(param.RoleID, param.PermissionID)

	trans := controller.GetTranslator(c, adminUserController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteRolePermission")
	controller.Response(c, 200, message, nil)
}

func (adminUserController *AdminUserController) GetPermissionsList(c *gin.Context) {
	permissionsList := adminUserController.userService.GetPermissionsList()
	controller.Response(c, 200, "", permissionsList)
}

func (adminUserController *AdminUserController) UpdateUserRoles(c *gin.Context) {
	type userRolesParams struct {
		Roles []string `json:"roles" validate:"required"`
		Email string   `json:"email" validate:"required"`
	}
	param := controller.Validated[userRolesParams](c, &adminUserController.constants.Context)
	adminUserController.userService.UpdateUserRoles(param.Email, param.Roles)

	trans := controller.GetTranslator(c, adminUserController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updateUserRole")
	controller.Response(c, 200, message, nil)
}

func (adminUserController *AdminUserController) DeleteUserRole(c *gin.Context) {
	type deleteUserRolesParams struct {
		Email  string `json:"email" validate:"required"`
		RoleID uint   `uri:"roleID" validate:"required"`
	}
	param := controller.Validated[deleteUserRolesParams](c, &adminUserController.constants.Context)
	adminUserController.userService.DeleteUserRole(param.Email, param.RoleID)

	trans := controller.GetTranslator(c, adminUserController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteUserRole")
	controller.Response(c, 200, message, nil)
}
