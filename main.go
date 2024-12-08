package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/verify", handleRequest)
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	log.Println("Starting server on port:", httpPort)
	if err := http.ListenAndServe(":"+httpPort, logRequest(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email address required", http.StatusBadRequest)
		return
	}

	result := VerifyResult{Email: email}
	result.Verify()
	json.NewEncoder(w).Encode(result)
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
