package unit

import (
	"testing"

	"github.com/decoesp/tamborete/internal/database"
)

func TestSetGet(t *testing.T) {
	db := database.New()
	db.Set("test", "value")

	val, exists := db.Get("test")
	if !exists {
		t.Error("Key should exist")
	}
	if val != "value" {
		t.Errorf("Expected 'value', got '%s'", val)
	}
}
