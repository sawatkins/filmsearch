package handlers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
)

const MODEL = "gpt-4o" 
const PROMPT = "You are search engine with very deep and extensive knowledge of movies. A user will provide text describing a movie they are trying to find. Your job is to find the user's movie based on their description. Respont in valid JSON with only the relevant movie titles, years they were released, and the short 1-2 sentece description justifying why that movie is relevant user's description. Use the format from this example: { \"movies\": [ {  \"title\": \"movie title 1\", \"year\": 2003, \"justification\": \"justification sentences for movie 1 goes here\"}, {\"title\": \"movie 2 year\",\"year\": 2013, \"justification\": \"justification sentences for movie 2 goes here\"}, ect.]}"

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
		Model: MODEL,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: PROMPT},
			{Role: "user", Content: query},
		},
	}
	resp, err := openaiClient.CreateChatCompletion(context.TODO(), request)
	if err != nil {
		fmt.Println(err)
		return "", err // or nil?
	}
	return resp.Choices[0].Message.Content, nil // should i keep this as a string?, or cast to a json equivilent?
}
