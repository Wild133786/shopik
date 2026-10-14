package grpcclient

import (
	"context"
	"shopik/services/catalog/pkg/catalogpb"
	"shopik/services/order/internal/service"
)

type CatalogClient struct {
	client catalogpb.CatalogServiceClient
}

func NewCatalogClient(client catalogpb.CatalogServiceClient) *CatalogClient {
	return &CatalogClient{client: client}
}

func (c *CatalogClient) GetProductsByIDs(ctx context.Context, ids []int64) ([]service.CatalogProduct, []int64, error) {
	resp, err := c.client.GetProductsByIDs(ctx, &catalogpb.GetProductsByIDsRequest{ProductIds: ids})
	if err != nil {
		return nil, nil, err
	}

	products := make([]service.CatalogProduct, len(resp.Products))
	for i, p := range resp.Products {
		products[i] = service.CatalogProduct{ID: p.Id, Name: p.Name, Price: p.Price, Stock: p.Stock}
	}

	return products, resp.MissingProductIds, nil
}
