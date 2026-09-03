package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"micro-warehouse/product-service/configs"
	"micro-warehouse/product-service/pkg/jwt"

	"github.com/gofiber/fiber/v2/log"
)

type MerchantClient struct {
	UrlApiGateway string
	httpClient    *http.Client
	config        configs.Config
}

type MerchantProductResponse struct {
	ProductID  uint `json:"product_id"`
	TotalStock int  `json:"total_stock"`
}

type MerchantProductServiceResponse struct {
	Message string                  `json:"message"`
	Data    MerchantProductResponse `json:"data"`
	Error   uint                    `json:"error,omitempty"`
}

func (mc *MerchantClient) generateInternalToken() (string, error) {
	return jwt.GenerateInternalToken(mc.config)
}

func NewMerchantClient(cfg configs.Config) *MerchantClient {
	return &MerchantClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		UrlApiGateway: cfg.App.UrlApiGateway,
		config:        cfg,
	}
}

func (mc *MerchantClient) GetProductStockAcrossMerchant(ctx context.Context, productID uint) (int, error) {
	url := fmt.Sprintf("%s/api/v1/merchant-products/%d/total-stock", mc.UrlApiGateway, productID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[MerchantClient] GetProductStockAcrossMerchant - 1: %v", err)
		return 0, err
	}

	token, err := mc.generateInternalToken()
	if err != nil {
		log.Errorf("[MerchantClient] GetProductStockAcrossMerchants - 2: %v", err)
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Internal-Request", "true")
	req.Header.Set("X-Gateway", "warehouse-api-gateway")

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		log.Errorf("[MerchantClient] GetProductStockAcrossMerchant - 2: %v", err)
		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[MerchantClient] GetProductStockAcrossMerchant - 3: %v", err)
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[MerchantClient] GetProductStockAcrossMerchant - 4: %s", string(body))
		return 0, errors.New("failed to get product stock across merchants")
	}

	var stockResp MerchantProductServiceResponse
	err = json.Unmarshal(body, &stockResp)
	if err != nil {
		log.Errorf("[MerchantClient] GetProductStockAcrossMerchant - 5: %v", err)
		return 0, errors.New("failed to parse response from merchant service")
	}

	return stockResp.Data.TotalStock, nil
}

func (mc *MerchantClient) DeleteAllProductMerchantProducts(ctx context.Context, productID uint) error {
	url := fmt.Sprintf("%s/api/v1/merchant-products/%d/total-stock", mc.UrlApiGateway, productID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		log.Errorf("[MerchantClient] DeleteAllProductMerchantProducts - 1: %v", err)
		return err
	}

	token, err := mc.generateInternalToken()
	if err != nil {
		log.Errorf("[MerchantClient] DeleteAllProductMerchantProducts - 2: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Internal-Request", "true")
	req.Header.Set("X-Gateway", "warehouse-api-gateway")

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		log.Errorf("[MerchantClient] DeleteAllProductMerchantProducts - 2: %v", err)
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[MerchantClient] DeleteAllProductMerchantProducts - 3: %v", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[MerchantClient] DeleteAllProductMerchantProducts - 4: %s", string(body))
		return errors.New("failed to get product stock across merchants")
	}

	return nil
}
