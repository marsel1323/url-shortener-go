package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupRouter() (*gin.Engine, Storage) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	storage := NewInMemoryStorage()
	server := NewServer(storage)

	r.POST("/", server.HandlePostRequest)
	r.GET("/:id", server.HandleGetRequest)

	return r, storage
}

func TestHandlePostRequest(t *testing.T) {
	router, _ := setupRouter()

	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, responseRecorder.Code)
	}
}

func TestHandlePostRequestInvalidURL(t *testing.T) {
	router, _ := setupRouter()

	responseRecorder := httptest.NewRecorder()
	invalidURL := "invalid-url"
	request, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(invalidURL))

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for invalid URL, got %d", http.StatusBadRequest, responseRecorder.Code)
	}
}

func TestHandleGetRequest(t *testing.T) {
	router, storage := setupRouter()

	responseRecorder := httptest.NewRecorder()
	testURL := "https://example.com"
	key, _ := storage.Save(testURL)
	request, _ := http.NewRequest(http.MethodGet, "/"+key, nil)

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusTemporaryRedirect || responseRecorder.Header().Get("Location") != testURL {
		t.Errorf("GET request for existing URL failed, expected %d and redirect to %s, got %d and %s", http.StatusTemporaryRedirect, testURL, responseRecorder.Code, responseRecorder.Header().Get("Location"))
	}

	responseRecorder = httptest.NewRecorder()
	request, _ = http.NewRequest(http.MethodGet, "/nonexistent", nil)

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d for nonexistent URL, got %d", http.StatusNotFound, responseRecorder.Code)
	}
}

func TestHandlePostRequestEmptyBody(t *testing.T) {
	router, _ := setupRouter()

	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for empty body, got %d", http.StatusBadRequest, responseRecorder.Code)
	}
}

func TestHandlePostRequestLongURL(t *testing.T) {
	router, _ := setupRouter()

	responseRecorder := httptest.NewRecorder()
	longURL := strings.Repeat("a", 10000) // очень длинный URL
	request, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(longURL))

	router.ServeHTTP(responseRecorder, request)

	// Здесь ожидаемый результат зависит от логики вашего приложения.
	// Например, можно ожидать http.StatusBadRequest если есть ограничение на длину URL.
}

func TestHandleGetRequestRootPath(t *testing.T) {
	router, _ := setupRouter()

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d for root path, got %d", http.StatusNotFound, responseRecorder.Code)
	}
}
