package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct{
	Message string `json:"Message"`
	Name string `json:"Name"`
	Age string `json:"Age"`
	Sex string `json:"Sex"`
}


func main(){
	port := ":8080"

	homeHandler := func(res http.ResponseWriter,req *http.Request){
		fmt.Fprintf(res,"Hello World")
	}
	healthHandler := func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w,"Server is Healthy.")
	}
	queryHandler := func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintln(w,"This is for the Query Test")
		fmt.Fprintln(w,"Method : ",r.Method)
		fmt.Fprintln(w,"URL : ",r.URL)
		fmt.Fprintln(w,"Path : ",r.URL.Path)
		fmt.Fprintln(w,"RawQuery : ",r.URL.RawQuery)
		fmt.Fprintln(w,"Name : ",r.URL.Query().Get("name"))
		fmt.Fprintln(w,"Hobby : ",r.URL.Query()["hobby"])
		fmt.Fprintf(w,"\n\n")


		fmt.Fprintln(w,"Header Details : ")
		headerDetails := r.Header
		for key,val := range headerDetails{
			fmt.Fprintf(w,"%s : ",key)
			fmt.Fprint(w,val)
			fmt.Fprintln(w)
		}

		authValue,exists := headerDetails["Auth"]
		if exists {
			fmt.Fprintln(w,authValue)
		}

		body,err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w,"Unable to ReadBody",http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w,string(body))

	}

	resBodyHandler := func (w http.ResponseWriter,r *http.Request){
		w.Header().Set("Auth-Type","OAuth")
		w.Header().Set("AppName","Go-Backend")

		w.Write([]byte("Response Body"))

		w.WriteHeader(http.StatusOK)
	}

	userDetails := func(w http.ResponseWriter, r *http.Request){
		reqQuery := r.URL.Query()
		Name := reqQuery.Get("Name")
		Age := reqQuery.Get("Age")
		Sex := reqQuery.Get("Sex")

		w.WriteHeader(http.StatusAccepted)
		w.Header().Set("Application","Go-Backend")
		
		response := Response{
			Message : "User Fetched Successfully!!",
			Name: Name,
			Age: Age,
			Sex: Sex,
		}

		json.NewEncoder(w).Encode(response)
	}

	productCategoryHandler := func(w http.ResponseWriter,r *http.Request){
		category := r.URL.Query().Get("category")
		api := r.Header.Get("Api-key")

		if api == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(struct{
				Message string `json:"message"`
			}{
				Message: "API Key Not Found.",
			})
			return
		}

		if category == ""{
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(struct{
				Message string `json:"message"`
			}{
				Message: "Category not Found.",
			})
			return
		}

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(struct{
				Message string `json:"message"`
			}{
				Message: "Category Found.",
			})
	}

	http.HandleFunc("/",homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/requestBody",queryHandler)
	http.HandleFunc("/responseBody",resBodyHandler)
	http.HandleFunc("/user",userDetails)
	http.HandleFunc("/product",productCategoryHandler)

	fmt.Println("Server is running on PORT : 8080")
	http.ListenAndServe(port,nil)
}