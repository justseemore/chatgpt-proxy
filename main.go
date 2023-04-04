package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"log"
	"os"
)

func main() {
	app := fiber.New()
	app.Static("/", "./api/views")
	app.All("/v1/*", func(ctx *fiber.Ctx) error {
		proxyUrl := fmt.Sprintf("https://%s/v1/%s", os.Getenv("PROXY_DOMAIN"), ctx.Params("*", ""))
		log.Println(proxyUrl)
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
