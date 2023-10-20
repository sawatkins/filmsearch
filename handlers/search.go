package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
)

func Search(c *fiber.Ctx) error {
	query := c.Query("q")

	// read data into map
	data, err := getData()
	if err != nil {
		fmt.Println(err)
	}
	// loop through map & select each movie
	// fmt.Println(data)
	movies, ok := data["movies"].(map[string]interface{})
	if !ok {
		fmt.Println("Could not convert data['movies'] to []string")
		return err
	}
	for key, _ := range movies {
		fmt.Println("title: " + key)
	}

	return c.Render("search", fiber.Map{
		"Query":   query,
		"Results": "not yet",
	}, "layouts/main")
}

func getData() (map[string]interface{}, error) {
	file, err := os.Open("scripts/topmovieswkeywords.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	return data, nil
}
