package jane

import "testing"

func TestMemoryStoreGet(t *testing.T) {
	store := NewMemoryStore()
	store.Set("name", "jojo")
	if val := store.Get("name"); val != "jojo" {
		t.Errorf("Get() = %v, want %v", val, "jojo")
	}
}

func TestMemoryStoreInvalidGet(t *testing.T) {
	store := NewMemoryStore()
	if val := store.Get("some rando key"); val != "" {
		t.Errorf("Get() = %v, want %v", val, "")
	}
}

func TestMemoryStoreClose(t *testing.T) {
	store := NewMemoryStore()
	if !store.Close() {
		t.Errorf("Close() failed")
	}
}
