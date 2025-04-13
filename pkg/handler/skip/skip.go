package skip

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/skip"
)

var ProxyTrusted fiber.Handler = skip.New(func(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusForbidden)
}, func(c *fiber.Ctx) bool {
	return c.IsProxyTrusted()
})
