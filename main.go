package main

import (
	"rest-api/database"
	"rest-api/routes"

	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("hello world!")
	})

	app.Post("/api/users", routes.CreateUser)
	app.Get("/api/users", routes.GetAllUsers)
	app.Get("/api/users/:id", routes.GetUserById)
	app.Put("/api/users/:id", routes.UpdateUser)
	app.Delete("/api/users/:id", routes.DeleteUser)

	app.Post("/api/products", routes.CreateProduct)
	app.Get("/api/products", routes.GetProducts)
	app.Get("/api/products/:id", routes.GetProductById)
	app.Put("/api/products/:id", routes.UpdateProduct)
	app.Delete("/api/products/:id", routes.DeleteProductById)

	app.Post("/api/orders", routes.CreateOrder)
	app.Get("/api/orders", routes.GetOrders)
	app.Get("/api/orders/:id", routes.GetOrderById)
	app.Put("/api/orders/:id", routes.UpdateOrder)
	app.Delete("/api/orders/:id", routes.DeleteOrder)
}

func main() {

	database.ConnectDB()
	app := fiber.New()

	setupRoutes(app)

	app.Listen(":3000")
}
