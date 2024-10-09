package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"test/server"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleSearch(t *testing.T) {
	// Установка переменных окружения для теста

	// Создание тестового запроса
	requestBody, _ := json.Marshal(server.SearchRequest{Query: "test query"})
	req, err := http.NewRequest(http.MethodPost, "/api/address/search", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	// Запуск тестового сервера
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleSearch)

	// Вызов обработчика
	handler.ServeHTTP(rr, req)

	// Проверка ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	var response []server.Address
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

}

func TestHandleGeocode(t *testing.T) {
	// Создание тестового запроса
	requestBody, _ := json.Marshal(server.GeocodeRequest{Lat: "55.878", Lng: "37.653"})
	req, err := http.NewRequest(http.MethodPost, "/api/address/geocode", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	// Запуск тестового сервера
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleGeocode)

	// Вызов обработчика
	handler.ServeHTTP(rr, req)

	// Проверка ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Логирование тела ответа для диагностики
	body := rr.Body.String()
	t.Logf("Response body: %s", body)

	// Проверка содержимого ответа
	var response server.Address
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

}
