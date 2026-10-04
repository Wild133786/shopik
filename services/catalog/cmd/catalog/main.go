package main

import (
	"log"
	"net"
	"net/http"
	"shopik/services/catalog/internal/grpcserver"
	"shopik/services/catalog/internal/handler"
	"shopik/services/catalog/internal/repository"
	"shopik/services/catalog/internal/service"
	"shopik/services/catalog/pkg/catalogpb"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
)

func main() {
	repo := repository.NewInMemoryRepository()
	svc := service.NewCatalogService(repo)
	handl := handler.NewCatalogHandler(svc)
	catalogServer := grpcserver.NewCatalogServer(svc)
	grpcServer := grpc.NewServer()
	catalogpb.RegisterCatalogServiceServer(grpcServer, catalogServer)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	r.Get("/products", handl.GetAllProducts)
	addrHTTP := ":8082"
	addrgRPC := ":50051"

	go func() {
		lis, err := net.Listen("tcp", addrgRPC)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	err := http.ListenAndServe(addrHTTP, r)
	if err != nil {
		log.Fatalf("catalog server error: %v", err)
	}
}
