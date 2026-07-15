package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type TodoItem struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

type TodoList struct {
	mu     sync.RWMutex
	items  map[int]*TodoItem
	nextID int
}

func NewTodoList() *TodoList {
	return &TodoList{
		items:  make(map[int]*TodoItem),
		nextID: 1,
	}
}

func (tl *TodoList) Add(title string) *TodoItem {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	item := &TodoItem{
		ID:        tl.nextID,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
	tl.items[tl.nextID] = item
	tl.nextID++
	return item
}

func (tl *TodoList) GetAll() []*TodoItem {
	tl.mu.RLock()
	defer tl.mu.RUnlock()

	items := make([]*TodoItem, 0, len(tl.items))
	for _, item := range tl.items {
		items = append(items, item)
	}
	return items
}

func (tl *TodoList) MarkDone(id int) (*TodoItem, bool) {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	item, exists := tl.items[id]
	if !exists {
		return nil, false
	}
	item.Done = true
	return item, true
}

func (tl *TodoList) Count() int {
	return len(tl.items)
}

func (tl *TodoList) Delete(id int) bool {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	_, exists := tl.items[id]
	if !exists {
		return false
	}
	delete(tl.items, id)
	return true
}

var todoList = NewTodoList()

func handleAddTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	item := todoList.Add(req.Title)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func handleGetTodos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	items := todoList.GetAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func handleMarkDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/todos/"):]
	if len(idStr) > 5 && idStr[len(idStr)-5:] == "/done" {
		idStr = idStr[:len(idStr)-5]
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	item, ok := todoList.MarkDone(id)
	if !ok {
		http.Error(w, "Todo item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ok := todoList.Delete(id)
	if !ok {
		http.Error(w, "Todo item not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/todos" {
		if r.Method == http.MethodPost {
			handleAddTodo(w, r)
		} else if r.Method == http.MethodGet {
			handleGetTodos(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	if len(path) > 7 && path[:7] == "/todos/" {
		if len(path) > 12 && path[len(path)-5:] == "/done" {
			handleMarkDone(w, r)
		} else {
			handleDeleteTodo(w, r)
		}
		return
	}

	http.NotFound(w, r)
}

func main() {
	http.HandleFunc("/todos", todoHandler)
	http.HandleFunc("/todos/", todoHandler)

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
