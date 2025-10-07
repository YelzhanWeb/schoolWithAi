package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Interest string `json:"interest"`
	Level    string `json:"level"`
}

type Course struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Level    string `json:"level"`
}

type RecommendResponse struct {
	User        string   `json:"user"`
	Recommended []Course `json:"recommended"`
}

func main() {
	user := User{
		ID:       1,
		Name:     "Али",
		Age:      13,
		Interest: "Программирование",
		Level:    "Начальный",
	}

	body, _ := json.Marshal(user)
	resp, err := http.Post("http://localhost:8000/recommend?limit=10", "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respData, _ := io.ReadAll(resp.Body)
	var result RecommendResponse
	json.Unmarshal(respData, &result)

	fmt.Printf("Рекомендации для %s:\n", result.User)
	for _, c := range result.Recommended {
		fmt.Printf("➡️  %s (%s, %s)\n", c.Name, c.Category, c.Level)
	}
}
