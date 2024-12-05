package controller_v1_private

import (
	"first-project/src/application"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	constants   *bootstrap.Constants
	userService *application.UserService
}

func NewRoleController(constants *bootstrap.Constants, userService *application.UserService) *RoleController {
	return &RoleController{
		constants:   constants,
		userService: userService,
	}
}

func (roleController *RoleController) CreateRole(c *gin.Context) {
	type createRolesParams struct {
		Permissions []string `json:"permissions" validate:"required"`
		RoleName    string   `json:"role" validate:"required"`
	}
	param := controller.Validated[createRolesParams](c, &roleController.constants.Context)
	role := roleController.userService.CreateNewRole(param.RoleName)
	roleController.userService.AssignPermissionsToRole(role.ID, param.Permissions)

	trans := controller.GetTranslator(c, roleController.constants.Context.Translator)
	message, _ := trans.T("successMessage.createRole")
	controller.Response(c, 200, message, nil)
}

func (roleController *RoleController) UpdateRole(c *gin.Context) {
	type updateRolesParams struct {
		Permissions []string `json:"permissions" validate:"required"`
		RoleID      uint     `uri:"roleID" validate:"required"`
	}
	param := controller.Validated[updateRolesParams](c, &roleController.constants.Context)
	roleController.userService.AssignPermissionsToRole(param.RoleID, param.Permissions)

	trans := controller.GetTranslator(c, roleController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updateRole")
	controller.Response(c, 200, message, nil)
}

func (roleController *RoleController) UpdateUserRoles(c *gin.Context) {
	type userRolesParams struct {
		Roles  []string `json:"roles" validate:"required"`
		UserID uint     `uri:"userID" validate:"required"`
	}
	param := controller.Validated[userRolesParams](c, &roleController.constants.Context)
	roleController.userService.UpdateUserRoles(param.UserID, param.Roles)

	trans := controller.GetTranslator(c, roleController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updateUserRole")
	controller.Response(c, 200, message, nil)
}
