package route

import (
	"github.com/gofiber/fiber/v2"

	cuser "chatbox/app/controller/user"

	hjwt "chatbox/pkg/handler/jwt"
)

func Route(router fiber.Router) {
	router.Post("/user", cuser.Register)

	router.Post("/user/login", cuser.Login)

	router.Post("/user/auth/refresh", hjwt.ValidateRefreshToken, cuser.Refresh)

	router.Post("/user/logout", cuser.Logout)
}
