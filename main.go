package main

import (
	"errors"
	"fmt"
)

// наши модкли
type Task struct {
	ID    int
	Title string
}

// репозитрой типа бдшки нашей
type TaskRepository struct {
	db     map[int]string
	nextID int
}

func NewRepository() *TaskRepository {
	return &TaskRepository{db: make(map[int]string), nextID: 1}
}

func (r *TaskRepository) SaveToPostgres(title string) (Task, error) {
	if title == "вызвать_ошибку_бд" {
		return Task{}, errors.New("postgres_error: connection timeout")
	}
	id := r.nextID
	r.db[id] = title
	r.nextID++
	return Task{ID: id, Title: title}, nil
}

// сервсс бищзес логика и валидация
type TaskService struct {
	repo *TaskRepository
}

func NewService(repo *TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(title string) (Task, error) {
	if title == "" {
		return Task{}, errors.New("validation_error: title cannot be empty")
	}
	return s.repo.SaveToPostgres(title)
}

// контролер интерфейса/хттп
type TaskController struct {
	service *TaskService
}

func NewController(service *TaskService) *TaskController {
	return &TaskController{service: service}
}

func (c *TaskController) HandlePostRequest(title string) string {
	task, err := c.service.CreateTask(title)
	if err != nil {
		return fmt.Sprintf("HTTP 400 Bad Request: %v", err)
	}
	return fmt.Sprintf("HTTP 201: Task ID %d saved successfully", task.ID)
}

func main() {
	fmt.Println("приложение готово к тестированию.")
}
