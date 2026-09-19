package main

import (
	"log"
	"net/http"
	"shopik/services/order/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Post("/orders", handler.CreateOrder)
	addr := ":8083"
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("ListenAndServe error: %+v", err)
	}

}
