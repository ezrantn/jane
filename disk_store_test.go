package jane

import (
	"os"
	"testing"
)

func TestDiskStoreGet(t *testing.T) {
	store, err := NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("error creating new disk to store %v", err)
	}
	defer os.Remove("test.db")
	err = store.Set("name", "ezra")
	if err != nil {
		t.Fatalf("error set key value to store %v", err)
	}
	if value, _ := store.Get("name"); value != "ezra" {
		t.Errorf("Get() = %v, want %v", value, "ezra")
	}
}

func TestDiskStoreInvalid(t *testing.T) {
	store, err := NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("error creating new disk to store %v", err)
	}
	defer os.Remove("test.db")
	if value, _ := store.Get("some key"); value != "" {
		t.Errorf("Get() = %v, want %v", value, "")
	}
}

func TestDiskStoreSetWithPersistence(t *testing.T) {
	store, err := NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("error creating new disk to store %v", err)
	}
	defer os.Remove("test.db")

	tests := map[string]string{
		"username": "john_doe",
		"email":    "john.doe@example.com",
		"address":  "1234 Elm Street",
		"phone":    "555-1234",
		"country":  "USA",
		"language": "English",
		"currency": "USD",
		"company":  "Acme Corp",
		"role":     "Developer",
		"status":   "Active",
	}

	for key, value := range tests {
		err = store.Set(key, value)
		if err != nil {
			t.Fatalf("error set key value to store %v", err)
		}

		got, errGetKey := store.Get(key)
		if errGetKey != nil {
			t.Fatalf("error getting value from store for key %s: %v", key, err)
		}

		if got != value {
			t.Errorf("Get() for key %s = %v, want %v", key, got, value)
		}
	}

	store.Close()
	store, err = NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}
	for key, value := range tests {
		got, errGetKey := store.Get(key)
		if errGetKey != nil {
			t.Fatalf("error getting value from store for key %s: %v", key, err)
		}

		if got != value {
			t.Errorf("Get() for key %s = %v, want %v", key, got, value)
		}
	}
	store.Close()
}

func TestDiskStoreDelete(t *testing.T) {
	store, err := NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("error creating new disk to store %v", err)
	}
	defer os.Remove("test.db")

	tests := map[string]string{
		"username": "john_doe",
		"email":    "john.doe@example.com",
		"address":  "1234 Elm Street",
		"phone":    "555-1234",
		"country":  "USA",
		"language": "English",
		"currency": "USD",
		"company":  "Acme Corp",
		"role":     "Developer",
		"status":   "Active",
	}

	for key, value := range tests {
		err = store.Set(key, value)
		if err != nil {
			t.Fatalf("error set key value to store %v", err)
		}
	}

	for key, _ := range tests {
		err = store.Set(key, "")
		if err != nil {
			t.Fatalf("error set key value to store %v", err)
		}
	}
	err = store.Set("end", "yes")
	store.Close()

	store, err = NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}

	for key := range tests {
		value, errGetKey := store.Get(key)
		if errGetKey != nil {
			t.Fatalf("error getting value from store for key %s: %v", key, err)
		}
		if value != "" {
			t.Errorf("expected value to be deleted for key %s, found: %s", key, value)
		}
	}

	value, err := store.Get("end")
	if err != nil {
		t.Fatalf("error getting value from store for key end: %v", err)
	}
	if value != "yes" {
		t.Errorf("expected value 'yes' for key 'end', found: %s", value)
	}

	store.Close()
}
