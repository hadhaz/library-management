package handlers

import (
	"app/server/domain"
	"app/server/services"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"strconv"
)

func GetAuthors(service services.AuthorsService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authors, err := service.GetAuthors(c.UserContext())
		if err != nil {
			slog.Error("Error getting authors: ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(domain.AuthorResponse{
			Authors: authors,
		})
	}
}

func GetAuthorByID(service services.AuthorsService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		paramID := c.Params("id")
		id, err := strconv.Atoi(paramID)

		author, err := service.GetAuthor(c.UserContext(), id)
		if err != nil {

		}

		return c.JSON(author)
	}
}

func sendError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(domain.ErrorResponse{
		Error: message,
	})
}
