package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ResponseTemplate struct {
	Message string `json:"Message"`
	Data []any `json:"Data"`
}

func LoggingMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Time : ",time.Now(),"Request : ",r.Method,r.URL)
		next.ServeHTTP(w,r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		if r.Header.Get("Authentication") == ""{
			w.WriteHeader(http.StatusNetworkAuthenticationRequired)
			json.NewEncoder(w).Encode(
				ResponseTemplate{
					Message: "Authentication Header Not Found.",
					Data: []any{},
				},
			)
			return
		}

		if r.Header.Get("Authentication") != "Password" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(
				ResponseTemplate{
					Message: "Invalid Authentication",
					Data: []any{},
				},
			)
			return
		}

		next.ServeHTTP(w,r)
	})
}

func HomeRouteHandler(w http.ResponseWriter, r *http.Request){
	res := ResponseTemplate{
		Message: "Home Route",
		Data: []any{},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func UserRouteHandler(w http.ResponseWriter, r *http.Request){
	res := ResponseTemplate{
		Message: "User Route",
		Data: []any{},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}



func main(){
	fmt.Println("Project API")
	port := ":8080"
	mux := http.NewServeMux()

	mux.HandleFunc("GET /",HomeRouteHandler)
	mux.HandleFunc("GET /user",UserRouteHandler)
	handler := LoggingMiddleware(AuthMiddleware(mux))
	

	http.ListenAndServe(port,handler)
}