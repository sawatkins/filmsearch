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

// const NEW_MODEL string = "gpt-4o"
// const NEW_PROMPT string = `A user will provide a description or aspect of a movie or movies they are trying to find. 
// Provide a list of movies, each movie's release year, and a justification (with very brief movie summary) of specific aspects of each movie that match the user's description.

// Movie resuts should come from a very deep and extensive knowldge of movies. Justifications (and very brief summary) should be around 2 sentences and not reveal any major spoilers.

// Responses MUST be valid JSON following this example format:
// { "movies": [ {  "title": "Movie Title 1", "year": 2000, "justification": "Justification (and very brief summary) sentences for movie 1 goes here"}, 
// 	{"title": "Movie Ritle 2","year": 2000, "justification": "Justification (and very brief summary) sentences for movie 2 goes here"}, ect...]}

// If there are no relevant movies or the user input doesn't make sense, respond with "{}" only.
// If a user respond in a language other than English, write the justification (and very brief summary) sentences in that language. 
// Respond with around 3-4 most relevant movies if the user input is more open ended. 
// Respond with only one movie if the user input is seeking one specific movie. 

// Be careful to make sure not to only respond with popular movies and that the movie actaully exists. 

// --
// For context: This is part of a personal project of mine which is a search/answer engine for discovering movies using natural language. 
// The idea is a user can find movies they don't know about by inputing a specific theme or movie aspect/element/piece of dialog/ect, users can find specific movie they forgot the name for, ect...
// Searching the interet for these questions ususally reveals just popular movies, so having deep extensive knowldege and being able to uncover uncommon movies really counts. 
// --
// `

// openaiMovieCompletion makes the call to the OpenAI API to get a list of matching movies
func openaiMovieCompletion(openaiClient *openai.Client, query string) (string, error) {
	request := openai.ChatCompletionRequest{
		Model: MODEL,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
		// Messages: []openai.ChatCompletionMessage{
		// 	{Role: "user", Content: PROMPT + "\n\nUser Query: " + query}, //o1 doesn't support system prompting yet
		// },
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
