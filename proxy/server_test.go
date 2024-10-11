package main

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// TestServer_Start проверяет, что сервер отвечает на запросы к /api.
func TestServer_Start(t *testing.T) {
	server := &Server{
		stop: make(chan struct{}),
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// Запускаем сервер в горутине
	go func() {
		defer wg.Done()    // Уменьшаем счетчик горутины при завершении
		_ = server.Start() // Игнорируем ошибки для простоты
	}()

	// Ждем, пока сервер запустится
	time.Sleep(100 * time.Millisecond)

	// Создаем тестовый запрос к /api
	req := httptest.NewRequest(http.MethodGet, "/api/", nil)
	w := httptest.NewRecorder()

	// Обработка запроса через ваш сервер
	server.httpSrv.Handler.ServeHTTP(w, req)

	// Проверка статуса ответа
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200, но получен %d", resp.StatusCode)
	}

	// Проверка тела ответа
	body := w.Body.String()
	if body != "Hello from API" {
		t.Fatalf("Ожидалось 'Hello from API', но получено '%s'", body)
	}

	// Остановка сервера
	if err := server.Stop(); err != nil {
		t.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	wg.Wait() // Ждем завершения горутины сервера
}
