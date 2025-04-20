package main

import (
	"log"
	"net/http"
)

// Создаю тип для записи статусов выполнения задачи
type TaskStatus string

// Записываю константы для статусов, чтобы затем можно было легко добавлять различные статусы
const (
	StatusRunning TaskStatus = "Задача выполняется..."
	StatusPending TaskStatus = "Задача в обработке..."
	StatusDone    TaskStatus = "Задача выполнена..."
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
	fmt.fPrintf("Привет от сервера на главной странице!")
}
