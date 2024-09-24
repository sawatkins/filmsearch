package handlers

import (
	"context"
	"log"

	openai "github.com/sashabaranov/go-openai"
)

const MODEL = "gpt-4o" 
const PROMPT = `You are search engine with very deep and extensive knowledge of movies. A user will provide text describing a movie they are trying to find. 
	Your job is to find the user's movie based on their description. Respont in valid JSON with only the relevant movie titles, years they were released, 
	and short 1-2 sentece justification of specific aspects of the movie that match the user's description. 
	Make sure to only respond with movies. If there are no relevant movies or the prompt is unclear, respond with "{}" only. 
	Use the response format from this example: 
	{ "movies": [ {  "title": "movie title 1", "year": 0000, "justification": "justification sentences for movie 1 goes here"}, 
	{"title": "movie 2 year","year": 0000, "justification": "justification sentences for movie 2 goes here"}, ect...]}`

// openaiMovieCompletion makes the call to the OpenAI API to get a list of matching movies
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
		log.Println(err)
		return "", err 
	}
	return resp.Choices[0].Message.Content, nil 
}
