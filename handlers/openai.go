package handlers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
)

func Openai(openaiClient *openai.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := openai.ChatCompletionRequest{
			Model: "gpt-3.5-turbo-1106",
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
			Messages: []openai.ChatCompletionMessage{
				{Role: "system", Content: "You are a movie search engine. Users will try to input text relating to a movie they are trying to find. Your job is to return the relevent movie. Respont in JSON."},
				{Role: "user", Content: c.Query("query")},
			},
		}
		resp, err := openaiClient.CreateChatCompletion(context.TODO(), request)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}
		return c.Status(200).JSON(fiber.Map{
			"response": resp.Choices.
		})
	}
}
