package handlers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
)

func Openai(openaiClient *openai.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		movies, err := openaiMovieCompletion(openaiClient, c.Query("query"))
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}
		return c.Status(200).JSON(movies)
	}
}

func openaiMovieCompletion(openaiClient *openai.Client, query string) (string, error) {
	request := openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo-1106",
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: "You are a movie search engine. Users will input text describing movies they are trying to find. Your job is to return the relevent movie. Respont in JSON with only the relevent movie titles and years they were released. Use the format from the example: { \"movies\": [ {  \"title\": \"Lost in Translation\", \"year\": 2003 }, {\"title\": \"Her\",\"year\": 2013}]}"},
			{Role: "user", Content: query},
		},
	}
	resp, err := openaiClient.CreateChatCompletion(context.TODO(), request)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return resp.Choices[0].Message.Content, nil // should i keep this as a string?, or cast to a json equivilent?
}
