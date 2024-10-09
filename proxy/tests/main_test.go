package tests

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"test/server"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	err := errors.New("Не получилось собрать структуру сервера")
	s := server.NewServer()
	if s == nil {
		t.Errorf("%v", err)
	}
}
func TestHandleRouteTest(t *testing.T) {
	req := httptest.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()

	server.HandleTest(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
}
func TestStart(t *testing.T) {
	// Создаем новый сервер
	s := server.NewServer()

	errChan := make(chan error, 1)

	// Запускаем сервер
	go func() {
		if err := s.Start(); err != nil {
			errChan <- fmt.Errorf("Ошибка при запуске сервера: %v", err)
			return
		}
		errChan <- nil
	}()

	// Таймаут для ожидания ответа от горутины
	select {
	case err := <-errChan:
		if err != nil {
			t.Fatalf("Ошибка при запуске сервера: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Превышено время ожидания запуска сервера")
	default:

		// Делаем запрос к серверу для проверки его работы
		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			t.Fatalf("Запрос к серверу не удался: %v", err)
		}
		// Проверяем статус ответа
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Ожидался статус 200, но получен %d", resp.StatusCode)
		}

		// Останавливаем сервер
		err = s.Stop()
		if err != nil {
			t.Errorf("%v", err)
		}
	}
}
