package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sync"
	"time"
)

// Создаю тип для записи статусов выполнения задачи
type TaskStatus string

// Записываю константы для статусов, чтобы затем можно было легко добавлять различные статусы
const (
	StatusRunning TaskStatus = "Задача выполняется..."
	StatusPending TaskStatus = "Задача в обработке..."
	StatusDone    TaskStatus = "Задача выполнена..."
)

// Создаю структуру Task, где прописаны все поля задачи
type Task struct {
	ID         string     `json:"id"`
	TaskStatus TaskStatus `json:"task_status"`
	Result     string     `json:"result"`
	CreatedAt  time.Time  `json:"created_at"`
}

var (
	// Создаю мапу, где будут храниться все задачи
	tasks = make(map[string]*Task)
	// Также создаю мьютекс, чтобы было удобнее блокировать код в горутинах
	taskMutex sync.Mutex
)

func main() {
	// Создаю собственную переменную-маршрутизатор
	mux := http.NewServeMux()

	// Создаю маршруты и присваиваю им нужные хендлеры
	mux.HandleFunc("/", mainPageHandler)
	mux.HandleFunc("/create", createTaskHandler)
	mux.HandleFunc("/result", resultTaskHandler)

	// Прописываю запуск сервера и сразу обрабатываю ошибку
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Не удалось запустить сервер!")
	}
}

// Прописываю хендлер для главной странциы
func mainPageHandler(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "Привет от сервера на главной странице!")
}

// Прописываю хендлер для создания задачи
func createTaskHandler(resp http.ResponseWriter, req *http.Request) {
	// Создаю переменную id, которая будет передана в структуру
	id := uuid.New().String()

	task := &Task{
		ID:         id,
		TaskStatus: StatusRunning,
		CreatedAt:  time.Now(),
	}

	// Блокирую доступ к мапе, чтобы добавить в нее запись
	taskMutex.Lock()
	tasks[id] = task
	taskMutex.Unlock()

	// Вызываю горутину, которая начнет выполнение задачи
	go executeTask(task)

	// Записываю в заголовке тип данных, который будет получен клиентом
	resp.Header().Set("Content-Type", "application/json")
	// Вывожу эти данные
	json.NewEncoder(resp).Encode(task)
}
