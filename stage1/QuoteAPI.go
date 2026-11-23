package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var Quotes = []QuoteTemplate{
	{
		Id:    1,
		Quote: "Silence speaks what words never could.",
	},
	{
		Id:    2,
		Quote: "Some journeys are meant to be walked alone.",
	},
	{
		Id:    3,
		Quote: "Lost in thoughts, found in silence.",
	},
	{
		Id:    4,
		Quote: "Not everything beautiful is meant to stay.",
	},
	{
		Id:    5,
		Quote: "A little broken, but still becoming.",
	},
}

type QuoteTemplate struct {
	Id    int    `json:"id"`
	Quote string `json:"quote"`
}

type ResponseTemplate struct {
	Success bool            `json:"success"`
	Data    []QuoteTemplate `json:"data"`
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Time : ", time.Now(), " Request : ", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})

}

func HomeRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	res := ResponseTemplate{
		Success: true,
		Data:    nil,
	}
	json.NewEncoder(w).Encode(res)
}

func GetAllQuotes(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	res := ResponseTemplate{
		Success: true,
		Data:    Quotes,
	}
	json.NewEncoder(w).Encode(res)
}

func GetQuoteById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid quote ID", http.StatusBadRequest)
		return
	}

	for _, quote := range Quotes {
		if quote.Id == id {
			w.Header().Set("Content-Type", "application/json")
			res := ResponseTemplate{
				Success: true,
				Data:    []QuoteTemplate{quote},
			}
			json.NewEncoder(w).Encode(res)
			return
		}
	}

	http.Error(w, "quote not found", http.StatusNotFound)
}

func main() {
	fmt.Println("Quote API")

	port := ":8080"
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", HomeRoute)
	mux.HandleFunc("GET /all", GetAllQuotes)
	mux.HandleFunc("GET /quote/{id}", GetQuoteById)

	handler := LoggingMiddleware(mux)
	http.ListenAndServe(port, handler)
}
