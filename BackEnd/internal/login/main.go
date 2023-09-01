package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/ltovarm/DomoticX/BackEnd/internal/query"
)

func main() {

	// Create a new HTTP server
	// handler := cors.Default().Handler(http.DefaultServeMux)
	r := mux.NewRouter()
	r.HandleFunc("/api/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/register", registerHandler).Methods("POST")

	http.Handle("/", r)
	log.Println("Server started on :3000")
	err := http.ListenAndServe(":3000", handlers.CompressHandler(r))
	if err != nil {
		log.Fatal(err)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var user query.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Insert the user into the database
	// Set-up connection
	my_db := query.NewDb()
	if err := my_db.ConnectToDatabaseFromEnvVar(); err != nil {
		log.Fatalf(" > Error connecting to db: %s\n", err)
	}
	defer my_db.CloseDatabase()

	println("Recibido user y pass: ", user.Username, user.Password)

	err = my_db.InsertUser(user)
	if err != nil {
		http.Error(w, "User registration error", http.StatusInternalServerError)
		return
	}
	println("Add to bd successful")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("break point 1")
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("break point 2")
	var user query.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Println("break point 3")

	// Insert the user into the database
	// Set-up connection
	my_db := query.NewDb()
	if err := my_db.ConnectToDatabaseFromEnvVar(); err != nil {
		log.Fatalf(" > Error connecting to db: %s\n", err)
	}
	defer my_db.CloseDatabase()
	log.Println("break point 4")
	println("Recibido user y pass: ", user.Username, user.Password)

	bLogin, err := my_db.CheckUser(user)
	if err != nil {
		http.Error(w, "User registration error", http.StatusInternalServerError)
		return
	}
	if bLogin {
		println("Log-in to bd successful")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Log-in successfully", http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Log-in reject", http.StatusUnauthorized)
	}
}
