package userAPI;

import (
	"fmt"
	"encoding/json"
	"net/http"
)

type Response struct{
	Message string `json:"Message"`
	Data any `json:"Data"`
}

type User struct{
	Name string `json:"Name"`
	Age int `json:"Age"`
}


func GetUserHandler(w http.ResponseWriter,r *http.Request){
	users := []User{
		{
			Name : "A",
			Age : 23,
		},
		{
			Name : "B",
			Age : 22,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

}

func PostUserHandler (w http.ResponseWriter,r *http.Request){
	var user User;

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

	w.WriteHeader(http.StatusCreated)

    json.NewEncoder(w).Encode(user)
	
}

func reqHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method {
	case http.MethodGet:
		GetUserHandler(w,r)
	case http.MethodPost:
		PostUserHandler(w,r)	
	}
}


func main(){
	fmt.Println("USER API")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /user",GetUserHandler)
	mux.HandleFunc("POST /user",PostUserHandler)
	http.ListenAndServe(":8080",mux)
}

