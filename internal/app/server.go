package app

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
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

func (s *Server) HandlePostRequest(c *gin.Context) {
	b, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	u := string(b)
	if !isValidURL(u) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	key, err := s.storage.Save(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusCreated, "http://localhost:8080/"+key)
}

func (s *Server) HandleAPIShorten(c *gin.Context) {
	var request struct {
		URL string `json:"url"`
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	defer c.Request.Body.Close()

	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if !isValidURL(request.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	key, err := s.storage.Save(request.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{"result": "http://localhost:8080/" + key}
	respBytes, err := json.Marshal(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Write(respBytes)
}

func (s *Server) HandleGetRequest(c *gin.Context) {
	id := c.Param("id")

	u, err := s.storage.Load(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, u)
}

func isValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
