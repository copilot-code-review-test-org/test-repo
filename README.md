# Todo API

A basic Golang REST API for managing an in-memory to-do list.

## Features

- Add new todo items
- Fetch all todo items
- Mark items as done
- Delete items
- Thread-safe in-memory storage

## Running the API

```bash
go run main.go
```

The server will start on `http://localhost:8080`.

## API Endpoints

### 1. Add a Todo Item

**POST** `/todos`

Request body:
```json
{
  "title": "Buy groceries"
}
```

Response (201 Created):
```json
{
  "id": 1,
  "title": "Buy groceries",
  "done": false,
  "created_at": "2025-10-20T23:54:16.986Z"
}
```

### 2. Get All Todo Items

**GET** `/todos`

Response (200 OK):
```json
[
  {
    "id": 1,
    "title": "Buy groceries",
    "done": false,
    "created_at": "2025-10-20T23:54:16.986Z"
  },
  {
    "id": 2,
    "title": "Walk the dog",
    "done": true,
    "created_at": "2025-10-20T23:55:30.123Z"
  }
]
```

### 3. Mark Todo as Done

**PUT** `/todos/{id}/done`

Response (200 OK):
```json
{
  "id": 1,
  "title": "Buy groceries",
  "done": true,
  "created_at": "2025-10-20T23:54:16.986Z"
}
```

### 4. Delete a Todo Item

**DELETE** `/todos/{id}`

Response (204 No Content)

## Testing

Run the tests with:

```bash
go test -v
```

## Building

Build the binary:

```bash
go build -o todo-api
```

Run the binary:

```bash
./todo-api
```

## Example Usage with curl

```bash
# Add a todo item
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Buy groceries"}'

# Get all todos
curl http://localhost:8080/todos

# Mark todo as done (replace 1 with actual ID)
curl -X PUT http://localhost:8080/todos/1/done

# Delete a todo (replace 1 with actual ID)
curl -X DELETE http://localhost:8080/todos/1
```
