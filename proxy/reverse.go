package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ReverseProxy структура для хранения целевого сервера
type ReverseProxy struct {
	host string
	port string
}

// NewReverseProxy создаёт новый инстанс прокси-сервера
func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

// ReverseProxy мидлварь для проксирования запросов
func (rp *ReverseProxy) ReverseProxy() http.Handler {
	targetURL, err := url.Parse("http://" + rp.host + ":" + rp.port)
	if err != nil {
		panic(err) // Обработка ошибки, если URL недопустим
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request for %s", r.URL.Path)

		if r.URL.Path == "/api" || r.URL.Path == "/api/" {
			log.Println("Returning 'Hello from API'")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello from API"))
			return
		}

		log.Println("Proxying request to", targetURL)
		proxy.ServeHTTP(w, r)
	})

}
