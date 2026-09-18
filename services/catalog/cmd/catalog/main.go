package main

import (
	"log"
	"net/http"
	"shopik/services/catalog/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	r.Get("/products", handler.GetAllProducts)
	addr := ":8082"
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("catalog server error: %v", err)
	}
}
