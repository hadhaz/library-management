package handlers

import (
	"log/slog"
	"strconv"
	"time"

	"app/server/domain"
	"app/server/services"

	"github.com/gofiber/fiber/v2"
)

// GetBooks returns a handler function that retrieves all books
func GetBooks(service services.BooksService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		books, err := service.GetBooks(c.UserContext())
		if err != nil {
			slog.Error("GetBooks failed", "error", err)
			return sendError(c, fiber.StatusInternalServerError, "internal error")
		}

		return c.JSON(domain.BooksResponse{
			Books: books,
		})
	}
}

// GetBook
func GetBook(service services.BooksService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		paramID := c.Params("id")
		id, err := strconv.Atoi(paramID)

		book, err := service.GetBook(c.UserContext(), id)
		if err != nil {
			slog.Error("GetBook failed", "error", err)
			return sendError(c, fiber.StatusInternalServerError, "internal error")
		}

		return c.JSON(book)
	}
}

// AddBook returns a handler function that adds a book
func AddBook(service services.BooksService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var book domain.Book
		if err := c.BodyParser(&book); err != nil {
			slog.Warn("AddBook request parsing failed", "error", err)
			return sendError(c, fiber.StatusBadRequest, "invalid request")
		}

		if book.Title == "" {
			return sendError(c, fiber.StatusBadRequest, "title is required")
		}

		if book.PublishDate.IsZero() {
			return sendError(c, fiber.StatusBadRequest, "published date is required")
		}
		if book.PublishDate.After(time.Now()) {
			return sendError(c, fiber.StatusBadRequest, "published date cannot be in the future")
		}

		err := service.SaveBook(c.UserContext(), book)
		if err != nil {
			slog.Error("AddBook failed", "error", err)
			return sendError(c, fiber.StatusInternalServerError, "internal error")
		}
		return c.SendStatus(fiber.StatusCreated)
	}
}

// DeleteBook
func DeleteBook(service services.BooksService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var book domain.Book
		if err := c.BodyParser(&book); err != nil {
			slog.Warn("DeleteBook request parsing failed", "error", err)
			return sendError(c, fiber.StatusBadRequest, "invalid request")
		}

		err := service.DeleteBook(c.UserContext(), book.ID)
		if err != nil {
			slog.Error("DeleteBook failed", "error", err)
			return sendError(c, fiber.StatusInternalServerError, "internal error")
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// UpdateBook
func UpdateBook(service services.BooksService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var book domain.Book
		if err := c.BodyParser(&book); err != nil {
			slog.Warn("UpdateBook request parsing failed", "error", err)
			return sendError(c, fiber.StatusBadRequest, "invalid request")
		}
		err := service.UpdateBook(c.UserContext(), book)
		if err != nil {
			slog.Error("UpdateBook failed", "error", err)
			return sendError(c, fiber.StatusInternalServerError, "internal error")
		}
		return c.SendStatus(fiber.StatusCreated)
	}
}

func sendError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(domain.ErrorResponse{
		Error: message,
	})
}
