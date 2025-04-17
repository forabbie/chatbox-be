package route

import (
	"github.com/gofiber/fiber/v2"

	cuser "chatbox/app/controller/user"

	hjwt "chatbox/pkg/handler/jwt"
)

func Route(router fiber.Router) {

	router.Post("/auth/register", cuser.Register)

	router.Post("/auth/login", cuser.Login)

	router.Post("/auth/refresh", hjwt.ValidateRefreshToken, cuser.Refresh)

	router.Post("/auth/logout", hjwt.ValidateAccessToken, cuser.Logout)

	router.Get("/users", hjwt.ValidateAccessToken, cuser.GetUsers)

	router.Get("/user/profile", hjwt.ValidateAccessToken, cuser.GetCurrentUser)

	router.Get("/user/:id", hjwt.ValidateAccessToken, cuser.GetUserDetails)
}
