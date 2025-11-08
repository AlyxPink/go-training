package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestUserStore(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - add proper TODOs")

	store := NewUserStore()

	// Test Create
	user := store.Create("Alice")
	if user.Name != "Alice" {
		t.Errorf("Name = %q, want \"Alice\"", user.Name)
	}
	
	// Test Get
	got, ok := store.Get(user.ID)
	if !ok {
		t.Fatal("Get() failed")
	}
	if got.Name != "Alice" {
		t.Errorf("Name = %q, want \"Alice\"", got.Name)
	}
	
	// Test GetAll
	users := store.GetAll()
	if len(users) != 1 {
		t.Errorf("GetAll() length = %d, want 1", len(users))
	}
	
	// Test Delete
	if !store.Delete(user.ID) {
		t.Error("Delete() failed")
	}
	
	if len(store.GetAll()) != 0 {
		t.Error("User not deleted")
	}
}

func TestServer_HandleUsers_GET(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - add proper TODOs")

	server := NewServer()
	server.store.Create("Alice")
	server.store.Create("Bob")
	
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	
	server.HandleUsers(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}
	
	var users []User
	json.NewDecoder(w.Body).Decode(&users)
	
	if len(users) != 2 {
		t.Errorf("Got %d users, want 2", len(users))
	}
}

func TestServer_HandleUsers_POST(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - add proper TODOs")

	server := NewServer()
	
	body := bytes.NewBufferString(`{"name":"Charlie"}`)
	req := httptest.NewRequest("POST", "/users", body)
	w := httptest.NewRecorder()
	
	server.HandleUsers(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusCreated)
	}
	
	var user User
	json.NewDecoder(w.Body).Decode(&user)
	
	if user.Name != "Charlie" {
		t.Errorf("Name = %q, want \"Charlie\"", user.Name)
	}
}

func TestServer_HandleUser_GET(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - add proper TODOs")

	server := NewServer()
	user := server.store.Create("Alice")
	
	req := httptest.NewRequest("GET", "/users/"+strconv.Itoa(user.ID), nil)
	w := httptest.NewRecorder()
	
	server.HandleUser(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestServer_HandleUser_DELETE(t *testing.T) {
	t.Fatal("TODO: This exercise is fully implemented but should be starter code - add proper TODOs")

	server := NewServer()
	user := server.store.Create("Alice")
	
	req := httptest.NewRequest("DELETE", "/users/"+strconv.Itoa(user.ID), nil)
	w := httptest.NewRecorder()
	
	server.HandleUser(w, req)
	
	if w.Code != http.StatusNoContent {
		t.Errorf("Status = %d, want %d", w.Code, http.StatusNoContent)
	}
	
	if len(server.store.GetAll()) != 0 {
		t.Error("User not deleted")
	}
}
