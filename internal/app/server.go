package app

import (
	"io"
	"net/http"
	"net/url"
)

type Server struct {
	storage Storage
}

func NewServer(storage Storage) *Server {
	return &Server{storage: storage}
}

func (s *Server) ListenAndServe() error {
	http.HandleFunc("/", s.handleRequests)
	return http.ListenAndServe("localhost:8080", nil)
}

func (s *Server) handleRequests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.handlePostRequest(w, r)
	case http.MethodGet:
		s.handleGetRequest(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handlePostRequest(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u := string(b)
	if !isValidURL(u) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	key, err := s.storage.Save(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + key))
}

func (s *Server) handleGetRequest(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]

	u, err := s.storage.Load(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Location", u)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func isValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
