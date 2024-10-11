package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func main() {
	s := NewServer()
	s.Start()
}

// Структура данных сервера
type Server struct {
	stop    chan struct{}
	handler *chi.Mux
	httpSrv *http.Server
}

func NewServer() *Server {
	return &Server{
		stop:    make(chan struct{}),
		handler: chi.NewRouter(),
	}
}

// Start запускает сервер
func (s *Server) Start() error {
	// Инициализация роутера chi
	r := chi.NewRouter()

	// Настройка CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:1313"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}))

	r.Use(middleware.Logger) // Логирование каждого запроса

	// Создание инстанса реверс-прокси для запросов с префиксом /api/
	rp := NewReverseProxy("localhost", "8080")
	r.Mount("/api/", rp.ReverseProxy())

	// Создание реверс-прокси для всех остальных запросов
	frontendProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "hugo:1313",
	})

	// Обработчик для всех остальных запросов
	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		frontendProxy.ServeHTTP(w, r)
	}))

	// Конфигурация HTTP-сервера
	port := ":8080"
	s.httpSrv = &http.Server{
		Addr:    port,
		Handler: r,
	}

	// Запуск сервера в горутине
	go func() {
		log.Printf("Сервер запущен на порту %s", port)
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ожидание сигнала остановки
	<-s.stop
	return nil
}

func (s *Server) Stop() error {
	// Закрываем канал, чтобы сигнализировать об остановке
	close(s.stop)

	// Ожидаем завершения активных соединений (graceful shutdown)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpSrv.Shutdown(ctx); err != nil {
		return fmt.Errorf("ошибка остановки сервера %v", err)
	}
	log.Println("Сервер успешно остановлен")
	return nil
}
