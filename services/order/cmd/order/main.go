package main

import (
	"log"
	"net/http"
	"shopik/services/catalog/pkg/catalogpb"
	"shopik/services/order/internal/grpcclient"
	"shopik/services/order/internal/handler"
	"shopik/services/order/internal/repository"
	"shopik/services/order/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	catalogAddr := "localhost:50051"
	conn, err := grpc.NewClient(catalogAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("start catalog client fatal error: %v", err)
	}
	defer conn.Close()
	catalogClient := grpcclient.NewCatalogClient(catalogpb.NewCatalogServiceClient(conn))
	repo := repository.NewInMemoryRepository()
	serv := service.NewOrderService(repo, catalogClient)
	hndl := handler.NewOrderHandler(serv)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Post("/orders", hndl.CreateOrder)
	addr := ":8083"
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("ListenAndServe error: %+v", err)
	}
}
