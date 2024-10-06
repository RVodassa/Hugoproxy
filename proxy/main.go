package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	// Инициализация роутера chi
	r := chi.NewRouter()

	// Использование стандартных middleware от chi
	r.Use(middleware.Logger) // Логирование каждого запроса

	// Инициализация структуры прокси-сервера
	rp := NewReverseProxy("hugo", "1313")

	// Добавляем маршрут для /api, который возвращает текст "Hello from API"
	r.Route("/api", func(r chi.Router) {
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello from API"))
		})
	})

	// Прокси для всех других маршрутов
	r.NotFound(rp.ReverseProxy(nil).ServeHTTP)

	log.Println("Server starting at :8080")
	http.ListenAndServe(":8080", r)
}

const content = ``

func WorkerTest() {
	t := time.NewTicker(1 * time.Second)
	var b byte = 0
	for {
		select {
		case <-t.C:
			err := os.WriteFile("/app/static/_index.md", []byte(fmt.Sprintf(content, b)), 0644)
			if err != nil {
				log.Println(err)
			}
			b++
		}
	}
}
