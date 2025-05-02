package main

import (
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	// raccount "chatbox/app/route/account"
	// rfile "chatbox/app/route/file"
	// rstatic "chatbox/app/route/static"
	rchannel "chatbox/app/route/channel"
	rdm "chatbox/app/route/dm"
	rmessage "chatbox/app/route/message"
	ruser "chatbox/app/route/user"

	"chatbox/pkg/channel"
	hjwt "chatbox/pkg/handler/jwt"
	hskip "chatbox/pkg/handler/skip"
	"chatbox/pkg/settings"
)

func New() *fiber.App {
	app := fiber.New(settings.FiberConfig)

	app.Use(
		recover.New(settings.RecoverConfig),
		logger.New(settings.LoggerConfig),
		requestid.New(settings.RequestIDConfig),
		cors.New(settings.CORSConfig),
		etag.New(settings.ETagConfig),
		cache.New(settings.CacheConfig),
	)

	// Skip if proxy not trusted
	app.Use(hskip.ProxyTrusted)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	ruser.Route(v1)
	rmessage.Route(v1)
	rchannel.Route(v1)
	rdm.Route(v1)
	// rfile.Route(v1)
	// rstatic.Route(v1)
	// raccount.Route(v1)

	// WS
	ws := app.Group("/ws")

	ws.Use(func(c *fiber.Ctx) error {
		if c.Query("v") == "1.0" {
			return c.Next()
		}

		return c.SendStatus(fiber.StatusNotFound)
	})

	ws.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	ws.Get("/chat/:id", hjwt.ValidateAccessToken, websocket.New(func(c *websocket.Conn) {
		channel.ChatHub.Register(c)
		defer channel.ChatHub.Unregister(c)
		for {
			_, p, err := c.ReadMessage()
			if err != nil {
				break
			}
			channel.ChatHub.Broadcast(c.Params("id"), p)
		}
	}))

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	app.Hooks().OnListen(func(listenData fiber.ListenData) error {
		fmt.Println("Server is listening...")
		return nil
	})

	app.Hooks().OnShutdown(func() error {
		fmt.Println("shutting down...")

		return nil
	})

	return app
}
