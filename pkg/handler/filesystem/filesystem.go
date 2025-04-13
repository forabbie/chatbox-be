package filesystem

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

var Static fiber.Handler = filesystem.New(filesystem.Config{
	Root: http.Dir("./tmp"),
})
