package app

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, container *Container) {
	api := app.Group("api/v1")

	roles := api.Group("/roles")
	roles.Post("/", container.RoleController.CreateRole)
	roles.Get("/", container.RoleController.GetAllRoles)
	roles.Get("/:id", container.RoleController.GetRoleById)
	roles.Put("/:id", container.RoleController.UpdateRole)
	roles.Delete("/:id", container.RoleController.DeleteRole)

	users := api.Group("/users")
	users.Post("/", container.UserController.CreateUser)
	users.Get("/", container.UserController.GetAllUsers)
	users.Get("/:id", container.UserController.GetUserById)
	users.Put("/:id", container.UserController.UpdateUser)
	users.Delete("/:id", container.UserController.DeleteUser)

	users.Get("/role/:roleName", container.UserController.GetUserRoleByRoleName)

	assignToRole := api.Group("assign-role")
	assignToRole.Post("/", container.UserController.AssignUserToRole)
	assignToRole.Get("/", container.UserController.GetAllUserRoles)
	assignToRole.Get("/:UserRoleID", container.UserController.GetUserRoleById)
	assignToRole.Put("/:id", container.UserController.EditAssignUserToRole)

	auth := api.Group("/auth")
	auth.Post("/login", container.AuthController.Login)

	upload := api.Group("/upload")
	upload.Post("/photo", container.UploadController.UploadPhoto)
}