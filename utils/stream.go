package utils

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sashabaranov/go-openai"
	"io"
	"os"
	"strings"
)

func ChatCompletionsStream(ctx *fiber.Ctx) error {
	openaiApikey := strings.Replace(ctx.Get("Authorization", os.Getenv("OPENAI_API_KEY")), "Bearer ", "", -1)
	var req openai.ChatCompletionRequest
	_ = json.Unmarshal(ctx.Body(), &req)
	client := openai.NewClient(openaiApikey)
	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		req,
	)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return err
	}
	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Transfer-Encoding", "chunked")
	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				fmt.Printf("\nStream error: %v\n", err)
				break
			}
			jsonStr, _ := json.Marshal(response)
			_, _ = fmt.Fprintln(w, fmt.Sprintf("data: %s\n", string(jsonStr)))
			_ = w.Flush()
		}
		_, _ = fmt.Fprintln(w, fmt.Sprintf("data: %s\n", "[DONE]"))
		_ = w.Flush()
	})
	return nil
}
