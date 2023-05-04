package handler

import (
	"fmt"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"net/http"
	"os"
	"strings"
)

// Handler is the main entry point of the application. Think of it like the main() method
func Handler(w http.ResponseWriter, r *http.Request) {
	// This is needed to set the proper request path in `*fiber.Ctx`
	r.RequestURI = r.URL.String()

	handler().ServeHTTP(w, r)
}

// building the fiber application
func handler() http.HandlerFunc {
	app := fiber.New()

	app.All("/*", func(ctx *fiber.Ctx) error {
		OsAuthApiKey := os.Getenv("AUTH_API_KEY")
		proxyHost := os.Getenv("PROXY_DOMAIN")
		if OsAuthApiKey != "" && ctx.Get("auth-api-key", "") == OsAuthApiKey {
			ctx.Request().Header.Del("Authorization")
			ctx.Request().Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))
			ctx.Request().Header.Del("auth-api-key")
		} else {
			proxyHost, _ = strings.CutSuffix(ctx.Get("h-proxy-host", proxyHost), "/")
		}
		path, _ := strings.CutPrefix(ctx.OriginalURL(), "/")
		proxyUrl := fmt.Sprintf("https://%s/%s", proxyHost, path)
		if err := proxy.Do(ctx, proxyUrl); err != nil {
			return err
		}
		ctx.Response().Header.Del(fiber.HeaderServer)
		return nil
	})
	return adaptor.FiberApp(app)
}
