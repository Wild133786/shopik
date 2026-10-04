package grpcserver

import (
	"context"
	"shopik/services/catalog/internal/service"
	"shopik/services/catalog/pkg/catalogpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CatalogServer struct {
	catalogpb.UnimplementedCatalogServiceServer
	service *service.CatalogService
}

func NewCatalogServer(svc *service.CatalogService) *CatalogServer {
	return &CatalogServer{service: svc}
}

func (s *CatalogServer) GetProductsByIDs(ctx context.Context,
	req *catalogpb.GetProductsByIDsRequest) (*catalogpb.GetProductsByIDsResponse,
	error) {
	var response catalogpb.GetProductsByIDsResponse
	products, missingIDs, err := s.service.GetProductsByIDs(ctx, req.ProductIds)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get product")
	}
	response.MissingProductIds = missingIDs

	for _, product := range products {
		respProduct := catalogpb.Product{
			Id:    product.ID,
			Name:  product.Name,
			Price: product.Price,
			Stock: product.Stock,
		}
		response.Products = append(response.Products, &respProduct)
	}

	return &response, nil
}
