package app

import (
	"github.com/gin-gonic/gin"
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
