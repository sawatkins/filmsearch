package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
)

type Movie struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Keywords []string `json:"keywords"`
}

type Movies struct {
	Movies []Movie `json:"movies"`
}

func Search(c *fiber.Ctx) error {
	query := c.Query("q")
	releventMovies := make([]string, 0)
	// read data into map
	// data, err := getData()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// loop through map & select each movie
	// fmt.Println(data)
	// movies, ok := data["movies"].(map[string]interface{})
	// if !ok {
	// 	fmt.Println("Could not convert data['movies'] to []string")
	// 	return err
	// }
	// for key, _ := range movies {
	// 	fmt.Println("title: " + key)
	// }

	jsonFile, err := os.Open("scripts/topmovieswkeywords.json")
	if err != nil {
		fmt.Println("Couldn't open the file:", err)
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var movies Movies
	// Unmarshal the JSON data into our Movies struct
	json.Unmarshal(byteValue, &movies)

	// Loop through each movie and print its values
	for _, movie := range movies.Movies {
		// fmt.Printf("Movie %d:\n", i+1)
		// fmt.Println("ID:", movie.ID)
		// fmt.Println("Title:", movie.Title)
		// fmt.Println("Keywords:", movie.Keywords)
		// first attempt ==> loop through every item, if keyword match, print title
		if slices.Contains(movie.Keywords, query) {
			releventMovies = append(releventMovies, movie.Title)
		}
	}

	return c.Render("search", fiber.Map{
		"Query":   query,
		"Results": releventMovies,
	}, "layouts/main")
}

// func getData() (map[string]interface{}, error) {
// 	file, err := os.Open("scripts/topmovieswkeywords.json")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	bytes, err := io.ReadAll(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var data map[string]interface{}
// 	if err := json.Unmarshal(bytes, &data); err != nil {
// 		return nil, err
// 	}

// 	return data, nil
// }
