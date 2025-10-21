package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddTodo(t *testing.T) {
	todoList = NewTodoList()

	reqBody := `{"title": "Test todo item"}`
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	handleAddTodo(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var item TodoItem
	if err := json.NewDecoder(w.Body).Decode(&item); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if item.Title != "Test todo item" {
		t.Errorf("Expected title 'Test todo item', got '%s'", item.Title)
	}

	if item.Done {
		t.Error("Expected item to not be done")
	}

	if item.ID != 1 {
		t.Errorf("Expected ID 1, got %d", item.ID)
	}
}

func TestAddTodoEmptyTitle(t *testing.T) {
	todoList = NewTodoList()

	reqBody := `{"title": ""}`
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	handleAddTodo(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetTodos(t *testing.T) {
	todoList = NewTodoList()
	todoList.Add("First item")
	todoList.Add("Second item")

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	handleGetTodos(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var items []*TodoItem
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}
}

func TestMarkDone(t *testing.T) {
	todoList = NewTodoList()
	item := todoList.Add("Test item")

	req := httptest.NewRequest(http.MethodPut, "/todos/1/done", nil)
	w := httptest.NewRecorder()

	handleMarkDone(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var updatedItem TodoItem
	if err := json.NewDecoder(w.Body).Decode(&updatedItem); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !updatedItem.Done {
		t.Error("Expected item to be marked as done")
	}

	if updatedItem.ID != item.ID {
		t.Errorf("Expected ID %d, got %d", item.ID, updatedItem.ID)
	}
}

func TestMarkDoneNotFound(t *testing.T) {
	todoList = NewTodoList()

	req := httptest.NewRequest(http.MethodPut, "/todos/999/done", nil)
	w := httptest.NewRecorder()

	handleMarkDone(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteTodo(t *testing.T) {
	todoList = NewTodoList()
	todoList.Add("Test item")

	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	w := httptest.NewRecorder()

	handleDeleteTodo(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
	}

	items := todoList.GetAll()
	if len(items) != 0 {
		t.Errorf("Expected 0 items after deletion, got %d", len(items))
	}
}

func TestDeleteTodoNotFound(t *testing.T) {
	todoList = NewTodoList()

	req := httptest.NewRequest(http.MethodDelete, "/todos/999", nil)
	w := httptest.NewRecorder()

	handleDeleteTodo(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestTodoListConcurrency(t *testing.T) {
	todoList = NewTodoList()

	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(id int) {
			todoList.Add("Item")
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	items := todoList.GetAll()
	if len(items) != 10 {
		t.Errorf("Expected 10 items, got %d", len(items))
	}
}
