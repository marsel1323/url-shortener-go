package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
)

var hashTable map[string]string

func main() {
	hashTable = make(map[string]string)

	http.HandleFunc("/", handleRequests)

	server := &http.Server{
		Addr: "localhost:8080",
	}

	log.Fatal(server.ListenAndServe())
}

func handleRequests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePostRequest(w, r)
	case http.MethodGet:
		handleGetRequest(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := string(b)
	if !isValidURL(url) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id := generateID(url)
	hashTable[id] = url

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + id))
}

func isValidURL(url string) bool {
	return url != ""
}

func generateID(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil)[:8])
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]

	url, ok := hashTable[id]
	if !ok {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
