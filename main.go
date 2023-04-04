package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"os"
)

func main() {
	app := fiber.New()
	app.All("/*", func(ctx *fiber.Ctx) error {
		proxyUrl := fmt.Sprintf("https://%s/%s", os.Getenv("PROXY_DOMAIN"), ctx.OriginalURL())
		if err := proxy.Do(ctx, proxyUrl); err != nil {
			return err
		}
		ctx.Response().Header.Del(fiber.HeaderServer)
		return nil
	})
	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
