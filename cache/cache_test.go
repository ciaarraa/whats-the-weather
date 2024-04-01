package cache

import (
	"os"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := tempDir + "/db"
	cacheFolder := tempDir + "/.cache"
	NewCache(dbPath, cacheFolder)
	if _, err := os.Stat(cacheFolder); os.IsNotExist(err) {
		t.Fatalf("Cache folder not created in the expected location: %v", err)
	}
}

func TestAdd(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := tempDir + "/db"
	cacheFolder := tempDir + "/.cache"
	cache := NewCache(dbPath, cacheFolder)
	cache.Add([]byte("test"), "key")

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatalf("DB folder not in expected location %v", err)
	}

	files, _ := os.ReadDir(cacheFolder)
	if len(files) != 1 {
		t.Fatalf("file not added to cache folder")
	}
}

func TestAddWithTTL(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := tempDir + "/db"
	cacheFolder := tempDir + "/.cache"
	cache := NewCache(dbPath, cacheFolder)
	cache.AddWithTTL([]byte("test"), "key", time.Hour)

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatalf("DB folder not in expected location %v", err)
	}

	files, _ := os.ReadDir(cacheFolder)
	if len(files) != 1 {
		t.Fatalf("file not added to cache folder")
	}
}

func TestGet(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := tempDir + "/db"
	cacheFolder := tempDir + "/.cache"
	cache := NewCache(dbPath, cacheFolder)
	cacheKey := "key"
	cache.Add([]byte("test"), cacheKey)

	result, _ := cache.Get(cacheKey)
	 if result != "test" {
	 	t.Fatalf("cached value not retrieved")
	 }
}
