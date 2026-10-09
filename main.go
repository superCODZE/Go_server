package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Number int    `json:"number"`
}

var (
	Users = make(map[int]User)
	mutex sync.RWMutex
)

func main() {

	mux := http.NewServeMux()

	// Define the routes
	mux.HandleFunc("GET /{$}", handleRoot)
	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc("GET /users/{id}", getUser)
	mux.HandleFunc("DELETE /users/{id}", deleteUser)

	fmt.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mutex.RLock()
	defer mutex.RUnlock()

	if err := tmpl.Execute(w, Users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	Users[len(Users)+1] = user

	w.WriteHeader(http.StatusCreated)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	mutex.RLock()
	defer mutex.RUnlock()
	user, exists := Users[id]

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	display, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(display)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if _, exists := Users[id]; !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	delete(Users, id)

	w.WriteHeader(http.StatusNoContent)
}
