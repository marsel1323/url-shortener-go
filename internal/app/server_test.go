package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlePostRequest(t *testing.T) {
	storage := NewInMemoryStorage()
	server := NewServer(storage)

	request, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	responseRecorder := httptest.NewRecorder()

	server.handleRequests(responseRecorder, request)

	if responseRecorder.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, responseRecorder.Code)
	}
}

func TestHandlePostRequestInvalidURL(t *testing.T) {
	storage := NewInMemoryStorage()
	server := NewServer(storage)

	invalidURL := "invalid-url"
	request, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(invalidURL))
	responseRecorder := httptest.NewRecorder()

	server.handleRequests(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for invalid URL, got %d", http.StatusBadRequest, responseRecorder.Code)
	}
}

func TestHandleGetRequest(t *testing.T) {
	storage := NewInMemoryStorage()
	server := NewServer(storage)

	testURL := "https://example.com"
	key, _ := storage.Save(testURL)
	request, _ := http.NewRequest(http.MethodGet, "/"+key, nil)
	responseRecorder := httptest.NewRecorder()
	server.handleRequests(responseRecorder, request)

	if responseRecorder.Code != http.StatusTemporaryRedirect || responseRecorder.Header().Get("Location") != testURL {
		t.Errorf("GET request for existing URL failed, expected %d and redirect to %s, got %d and %s", http.StatusTemporaryRedirect, testURL, responseRecorder.Code, responseRecorder.Header().Get("Location"))
	}

	request, _ = http.NewRequest(http.MethodGet, "/nonexistent", nil)
	responseRecorder = httptest.NewRecorder()
	server.handleRequests(responseRecorder, request)

	if responseRecorder.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d for nonexistent URL, got %d", http.StatusNotFound, responseRecorder.Code)
	}
}

func TestHandlePostRequestEmptyBody(t *testing.T) {
	storage := NewInMemoryStorage()
	server := NewServer(storage)

	request, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	responseRecorder := httptest.NewRecorder()

	server.handleRequests(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for empty body, got %d", http.StatusBadRequest, responseRecorder.Code)
	}
}

func TestHandlePostRequestLongURL(t *testing.T) {
	storage := NewInMemoryStorage()
	server := NewServer(storage)

	longURL := strings.Repeat("a", 10000) // очень длинный URL
	request, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(longURL))
	responseRecorder := httptest.NewRecorder()

	server.handleRequests(responseRecorder, request)

	// Здесь ожидаемый результат зависит от логики вашего приложения.
	// Например, можно ожидать http.StatusBadRequest если есть ограничение на длину URL.
}

func TestHandleGetRequestRootPath(t *testing.T) {
	storage := NewInMemoryStorage()
	server := NewServer(storage)

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	responseRecorder := httptest.NewRecorder()

	server.handleRequests(responseRecorder, request)

	if responseRecorder.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d for root path, got %d", http.StatusNotFound, responseRecorder.Code)
	}
}
