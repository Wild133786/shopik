package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	mustLoad()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	proxyUsers := NewProxy(os.Getenv("USER_ADDRESS"))
	proxyCatalog := NewProxy(os.Getenv("CATALOG_ADDRESS"))
	proxyOrders := NewProxy(os.Getenv("ORDER_ADDRESS"))
	r.Handle("/api/users/*", http.StripPrefix("/api/users", proxyUsers))
	r.Handle("/api/catalog/*", http.StripPrefix("/api/catalog", proxyCatalog))
	r.Handle("api/orders/*", http.StripPrefix("api/orders", proxyOrders))
	addr := ":8080"
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("http server error: %v", err)
	}
}

func NewProxy(host string) *httputil.ReverseProxy {
	target := &url.URL{
		Scheme: "http",
		Host:   host,
	}
	return httputil.NewSingleHostReverseProxy(target)
}

func mustLoad() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("err .env loading, using enviroment OS")
	}
	envList := []string{
		"USER_ADDRESS",
		"CATALOG_ADDRESS",
		"ORDER_ADDRESS",
	}
	var envNotExist []string
	for _, envName := range envList {
		if os.Getenv(envName) == "" {
			envNotExist = append(envNotExist, envName)
		}
	}
	if len(envNotExist) > 0 {
		log.Fatalf("env not exist: %v", envNotExist)
	}
}
