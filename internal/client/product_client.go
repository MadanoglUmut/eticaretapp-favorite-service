package client

import (
	"context"
	"encoding/json"
	"favorite_service/internal/models"
	"fmt"
	"net/http"
	"time"
)

type ProductClient struct {
	productServiceURL string
	maxRetries        int
	retryInterval     time.Duration
}

func NewProductClient(productServiceURL string, maxRetries int, retryInterval time.Duration) *ProductClient {
	return &ProductClient{
		productServiceURL: productServiceURL,
		maxRetries:        maxRetries,
		retryInterval:     retryInterval,
	}
}

func (c *ProductClient) VerifyProduct(ctx context.Context, productId int) (*models.Product, error) {
	var lastErr error

	for retry := 0; retry < c.maxRetries; retry++ {

		product, err := c.GetProduct(ctx, productId)
		if err == nil {
			return product, nil
		}

		lastErr = err

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(c.retryInterval):

		}

	}

	return nil, fmt.Errorf("istek sayisi (%d) , last error: %v", c.maxRetries, lastErr)
}

func (c *ProductClient) GetProduct(ctx context.Context, productId int) (*models.Product, error) {

	productServiceURL := fmt.Sprintf("%s/products/%d", c.productServiceURL, productId)

	reqProductService, err := http.NewRequestWithContext(ctx, "GET", productServiceURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	respProductService, err := client.Do(reqProductService)

	if err != nil {

		return nil, err

	}
	defer respProductService.Body.Close()

	if respProductService.StatusCode != http.StatusOK {

		return nil, models.ErrRecordNotFound

	}

	var productResponse models.ProductResponse

	if err := json.NewDecoder(respProductService.Body).Decode(&productResponse); err != nil {

		return nil, err

	}

	return &productResponse.SuccesData, nil

}
