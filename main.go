package main

import (
	"chatgpt-proxy/utils"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/sashabaranov/go-openai"
	"os"
	"strings"
)

func main() {
	app := fiber.New()
	app.All("/*", func(ctx *fiber.Ctx) error {
		OsAuthApiKey := os.Getenv("AUTH_API_KEY")
		if OsAuthApiKey != "" && ctx.Get("auth-api-key", "") == OsAuthApiKey {
			ctx.Request().Header.Del("Authorization")
			ctx.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))
			ctx.Request().Header.Del("auth-api-key")
		}
		if ctx.Path() == "/v1/chat/completions" {
			var req openai.ChatCompletionRequest
			_ = json.Unmarshal(ctx.Body(), &req)
			if req.Stream {
				return utils.ChatCompletionsStream(ctx)
			}
		}
		proxyHost, _ := strings.CutSuffix(ctx.Get("h-proxy-host", os.Getenv("PROXY_DOMAIN")), "/")
		proxyUrl := fmt.Sprintf("https://%s%s", proxyHost, ctx.Path())
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
