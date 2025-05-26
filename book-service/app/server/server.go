package server

import (
	"context"

	"app/datasources"
	"app/server/handlers"
	"app/server/services"

	"github.com/gofiber/fiber/v2"
)

// NewServer creates a new Fiber app and sets up the routes
func NewServer(ctx context.Context, dataSources *datasources.DataSources) *fiber.App {
	app := fiber.New()
	apiRoutes := app.Group("/api")

	apiRoutes.Get("/status", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	apiRoutes.Get("/v1/books", handlers.GetBooks(services.NewBooksService(dataSources.DB)))
	apiRoutes.Get("/v1/books:id", handlers.GetBook(services.NewBooksService(dataSources.DB)))
	apiRoutes.Post("/v1/books", handlers.AddBook(services.NewBooksService(dataSources.DB)))
	apiRoutes.Delete("/v1/books/:id", handlers.DeleteBook(services.NewBooksService(dataSources.DB)))
	apiRoutes.Put("/v1/books", handlers.UpdateBook(services.NewBooksService(dataSources.DB)))

	return app
}
