package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", handleTasks)
	mux.HandleFunc("/tasks/", handleTaskByID)

	fmt.Println("Server start")
	http.ListenAndServe(":8080", mux)

}

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var tasks []Task
var nextID = 1

func handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Отправить все задачи как JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		// Прочитать JSON из тела запроса
		var input struct {
			Title string `json:"title"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, `{"error": "неверный JSON"}`, http.StatusBadRequest)
			return
		}
		if input.Title == "" {
			http.Error(w, `{"error": "заголовок не может быть пустым"}`, http.StatusBadRequest)
			return
		}

		// Создать новую задачу
		task := Task{ID: nextID, Title: input.Title, Completed: false}
		tasks = append(tasks, task)
		nextID++

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)

	default:
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
	}

	// Логируем
	fmt.Printf("[%s] %s → %d\n", r.Method, r.URL.Path, 200) // упрощённо
}

func handleTaskByID(w http.ResponseWriter, r *http.Request) {
	// URL вида: /tasks/5
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "некорректный ID"}`, http.StatusBadRequest)
		return
	}

	// Найти задачу
	var task *Task
	for i := range tasks {
		if tasks[i].ID == id {
			task = &tasks[i]
			break
		}
	}
	if task == nil {
		http.Error(w, `{"error": "задача не найдена"}`, http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPut {
		var input struct {
			Completed bool `json:"completed"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, `{"error": "неверный JSON"}`, http.StatusBadRequest)
			return
		}

		task.Completed = input.Completed

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	} else {
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	fmt.Printf("[%s] %s → %d\n", r.Method, r.URL.Path, 200)
}
