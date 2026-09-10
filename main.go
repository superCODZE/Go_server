package main 

import "fmt"
import "net/http"

func main() {
    
	fmt.Println("Starting server on :8080")
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}" ,handlRoot)
	mux.HandleFunc("/hello", handlRoot)
	
	http.ListenAndServe(":8080", mux)
}

func handlRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}