package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"os"
	"strings"
)

func main() {
	app := fiber.New()
	app.All("/*", func(ctx *fiber.Ctx) error {
		proxyHost, _ := strings.CutSuffix(ctx.Get("h-proxy-host", os.Getenv("PROXY_DOMAIN")), "/")
		path, _ := strings.CutPrefix(ctx.OriginalURL(), "/")
		proxyUrl := fmt.Sprintf("https://%s/%s", proxyHost, path)
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
