package main

import (
	"log"
	"net/http"
	"shopik/services/user/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("user service hello"))
	})
	r.Post("/register", handler.Register)
	addr := ":8081"
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("http server error: %v", err)
	}
}
