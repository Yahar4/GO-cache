package cache

import (
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := New(0, 0)

	cache.Set("key1", "value1", 0)
	value, found := cache.GetItem("key1")

	if !found || value != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}
}

func TestCache_GetAllItems(t *testing.T) {
	cache := New(0, 0)

	cache.Set("key1", "value1", 10*time.Second)
	cache.Set("key2", "value2", 10*time.Second)

	allItems := cache.GetAllItems()

	if len(allItems) != 2 {
		t.Errorf("Expected 2 items, got %d", len(allItems))
	}
}

func TestCache_Delete(t *testing.T) {
	cache := New(0, 0)

	cache.Set("key1", "value1", 0)
	err := cache.Delete("key1")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	_, found := cache.GetItem("key1")
	if found {
		t.Error("Expected key1 to be deleted")
	}
}

func TestCache_Delete_NonExistentKey(t *testing.T) {
	cache := New(0, 0)

	err := cache.Delete("nonExistentKey")
	if err == nil {
		t.Error("Expected error for non-existent key, got nil")
	}
}

func TestCache_Count(t *testing.T) {
	cache := New(0, 0)

	cache.Set("key1", "value1", 0)
	cache.Set("key2", "value2", 0)

	count := cache.Count()
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

func TestCache_RenameKey(t *testing.T) {
	cache := New(0, 0)

	cache.Set("oldKey", "value1", 0)
	err := cache.RenameKey("oldKey", "newKey")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	value, found := cache.GetItem("newKey")
	if !found || value != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}

	_, found = cache.GetItem("oldKey")
	if found {
		t.Error("Expected oldKey to be deleted")
	}
}

func TestCache_RenameKey_ExistingNewKey(t *testing.T) {
	cache := New(0, 0)

	cache.Set("oldKey", "value1", 0)
	cache.Set("newKey", "value2", 0)

	err := cache.RenameKey("oldKey", "newKey")
	if err == nil {
		t.Error("Expected error for existing new key, got nil")
	}
}

func TestCache_GetItem_Expired(t *testing.T) {
	cache := New(0, 0)

	cache.Set("key1", "value1", 1*time.Nanosecond)
	time.Sleep(2 * time.Nanosecond)

	_, found := cache.GetItem("key1")
	if found {
		t.Error("Expected key1 to be expired")
	}
}

func TestCache_Increment(t *testing.T) {
	cache := New(0, 0)

	cache.Set("key1", int64(10), 0)
	err := cache.Increment("key1", 5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	val, _ := cache.GetItem("key1")
	if val != int64(15) {
		t.Errorf("Expected value 15, got %v", val)
	}
}

func TestCache_IncrementNonExistentKey(t *testing.T) {
	cache := New(0, 0)

	err := cache.Increment("non-existent", int64(5))
	if err == nil || err.Error() != "element to increment not found" {
		t.Errorf("expected 'element to increment not found' error, got %v", err)
	}
}

func TestCache_IncrementNonIntData(t *testing.T) {
	cache := New(0, 0)

	cache.Set("key1", "value1", 0)
	err := cache.Increment("key1", "1")
	if err == nil {
		t.Errorf("expected error, got %v", err)
	}
}
