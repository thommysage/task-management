package main

import (
	"strings"
	"testing"
)

func TestStep1_Repository(t *testing.T) {
	repo := NewRepository()

	// проврека сохранения в бд
	task, err := repo.SaveToPostgres("купить хлеб")
	if err != nil {
		t.Fatalf("ошибка нижнего уровня БД: %v", err)
	}
	if task.Title != "купить хлеб" || task.ID != 1 {
		t.Errorf("данные в бд записались некорректно")
	}
	//проверяем нижнй слой и убеждемся что метод реально добавляет задачу и присваетивает айдишку

}

func TestStep2_Service_Integration(t *testing.T) {
	repo := NewRepository()
	service := NewService(repo) // Интегрируем сервис с репозиторием

	// проверка бизнес-логики
	_, err := service.CreateTask("")
	if err == nil || !strings.Contains(err.Error(), "validation_error") {
		t.Errorf("сервис должен был отклонить пустую задачу на уровне логики")
	}

	// проверка записи через сервис в бдшку
	task, err := service.CreateTask("изучить go")
	if err != nil {
		t.Fatalf("ошибка при интеграции сервиса и бд: %v", err)
	}
	if task.ID != 1 {
		t.Errorf("связка сервис - бд отработала неверно")
	}
	//проверяем средний слой и убеждаемся что логика отсекает пустые задачи
	// и тк у нас имитация бдшки нам заглушки не нужны

}

func TestStep3_Full_Integration(t *testing.T) {
	repo := NewRepository()
	service := NewService(repo)
	controller := NewController(service) // Собрали всё приложение снизу вверх

	// Имитируем отправку HTTP запроса на создание
	response := controller.HandlePostRequest("сдать лабу")

	if !strings.Contains(response, "HTTP 201") {
		t.Errorf("сквозной тест провален. ответ сервера: %s", response)
	}
	// проверяем верхний слой и смотрим собрался ли весь слой

}
