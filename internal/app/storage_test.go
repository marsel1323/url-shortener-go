package app

import (
	"testing"
)

func TestInMemoryStorage_SaveAndLoad(t *testing.T) {
	storage := NewInMemoryStorage()
	url := "https://example.com"

	key, err := storage.Save(url)
	if err != nil {
		t.Fatalf("Save returned an error: %v", err)
	}

	loadedURL, err := storage.Load(key)
	if err != nil {
		t.Fatalf("Load returned an error: %v", err)
	}

	if loadedURL != url {
		t.Errorf("Expected URL %s, got %s", url, loadedURL)
	}
}

func TestInMemoryStorage_LoadNotFound(t *testing.T) {
	storage := NewInMemoryStorage()
	_, err := storage.Load("nonexistent")

	if err == nil {
		t.Errorf("Expected an error for a nonexistent key, but got none")
	}
}
