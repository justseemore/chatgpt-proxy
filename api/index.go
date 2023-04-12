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
		proxyHost, _ := strings.CutSuffix(ctx.Get("h-proxy-host", os.Getenv("PROXY_DOMAIN")), "/")
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
