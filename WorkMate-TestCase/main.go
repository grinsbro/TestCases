package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"
)

// Создаю тип для записи статусов выполнения задачи
type TaskStatus string

// Записываю константы для статусов, чтобы затем можно было легко добавлять различные статусы
const (
	StatusRunning TaskStatus = "Задача выполняется..."
	StatusPending TaskStatus = "Задача обрабатывается..."
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
	mux.HandleFunc("/result/", resultTaskHandler)

	// Добавлю еще хендлер для вывода всех задач
	mux.HandleFunc("/tasks/", allTasksHandler)

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

// Добавляю функцию, которая будет выполнять задачу
func executeTask(task *Task) {
	// Меняю статус задачи
	taskMutex.Lock()
	task.TaskStatus = StatusPending
	taskMutex.Unlock()

	// Симулирую долгую I/O задачу
	time.Sleep(time.Minute * time.Duration(rand.IntN(3)+3))

	// Меняю статус задачи на успешный
	taskMutex.Lock()
	task.TaskStatus = StatusDone
	task.Result = "Задача успешно выполнена!"
	taskMutex.Unlock()
}

// Прописываю хендлер для вывода результата задачи
func resultTaskHandler(resp http.ResponseWriter, req *http.Request) {
	// Получаю id и query параметров запроса
	id := req.URL.Query().Get("id")
	if id == "" {
		http.Error(resp, "Не передан id задачи", http.StatusBadRequest)
		return
	}

	// Получаю значение из мапы с задачами
	taskMutex.Lock()
	task, ok := tasks[id]
	taskMutex.Unlock()
	if !ok {
		http.Error(resp, "Не найдено такой задачи :(", http.StatusNotFound)
		return
	}

	// Если все проверки пройдены, то вывожу json с данными задачи
	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(task)
}

// Прописываю хендлер для вывода всех существующих задач
func allTasksHandler(resp http.ResponseWriter, req *http.Request) {
	taskMutex.Lock()
	task := tasks
	taskMutex.Unlock()
	if len(task) == 0 {
		http.Error(resp, "Нет ни одной задачи", http.StatusNotFound)
	}

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(task)
}
