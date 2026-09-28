package routes

import (
	"github.com/RomeyKung/social/controllers"
	validation "github.com/RomeyKung/social/validation"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	//auth
	app.Post("/user/signup", validation.ValidateCreateUser, controllers.Register)
	app.Post("/user/signin", validation.ValidateLoginUser, controllers.Login)
}
