package routes

import (
	"main/controllers"
	"main/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	auth := app.Group("/auth")
	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
	auth.Post("/logout", controllers.Logout)

	api := app.Group("/api")
	api.Use(middleware.Auth)
	api.Get("/user", controllers.User)
}
