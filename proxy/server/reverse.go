package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ReverseProxy структура для хранения целевого сервера (Hugo)
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

// ReverseProxy middleware для проксирования запросов
// ReverseProxy мидлварь для проксирования запросов на сервер hugo
func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проксируем запрос на сервер Hugo
		link := fmt.Sprintf("http://%s:%s", rp.host, rp.port)
		uri, err := url.Parse(link)
		if err != nil {
			http.Error(w, "Invalid proxy URL", http.StatusInternalServerError)
			return
		}

		r.Header.Set("Reverse-Proxy", "true")

		proxy := httputil.ReverseProxy{Director: func(req *http.Request) {
			req.URL.Scheme = uri.Scheme
			req.URL.Host = uri.Host
			req.URL.Path = r.URL.Path
			req.Host = uri.Host
		}}

		proxy.ServeHTTP(w, r)
	})
}
