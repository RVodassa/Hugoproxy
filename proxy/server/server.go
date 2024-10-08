package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

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

func HandleTest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // 200
}

func (s *Server) Start() error {
	// Инициализация роутера chi
	r := chi.NewRouter()

	r.Get("/test", HandleTest)
	// Использование стандартных middleware от chi
	r.Use(middleware.Logger) // Логирование каждого запроса

	// Эндпоинт для вашего swagger.yaml
	r.Get("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Запрос на swagger.yaml")
		if _, err := os.Stat("./swagger/swagger.yaml"); os.IsNotExist(err) {
			log.Println("swagger.yaml не найден!")
			http.Error(w, "swagger.yaml не найден!", http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, "./swagger/swagger.yaml")
	})

	// Раздача статических файлов из папки dist (Swagger UI)
	fs := http.FileServer(http.Dir("./swagger/static/swagger/dist"))
	r.Handle("/swagger/*", http.StripPrefix("/swagger/", fs))

	// Инициализация структуры прокси-сервера
	rp := NewReverseProxy("hugo", "1313")

	// Добавляем маршрут для /api, который возвращает текст "Hello from API"
	r.Route("/api", func(r chi.Router) {
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello from API"))
		})
	})

	r.Post("/api/address/geocode", HandleGeocode)
	r.Post("/api/address/search", HandleSearch)

	// Прокси для всех других маршрутов
	r.NotFound(rp.ReverseProxy(nil).ServeHTTP)

	port := ":8080"
	s.httpSrv = &http.Server{
		Addr:    port,
		Handler: r,
	}

	// Запуск сервера в горутине
	go func() {
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)

		}
		log.Printf("Сервер запущен на порту %s", port)
	}()

	<-s.stop // Ожидание сигнала остановки
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

// const content = ``

// func WorkerTest() {
// 	t := time.NewTicker(1 * time.Second)
// 	defer t.Stop() // Обязательно останавливаем таймер при выходе из функции
// 	var b byte = 0

// 	for range t.C { // Используем for range для обработки событий
// 		err := os.WriteFile("/app/static/_index.md", []byte(fmt.Sprint(content, b)), 0644)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		b++
// 	}
// }
