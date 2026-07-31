package response

import "micro-warehouse/product-service/pkg/pagination"

type ProductResponse struct {
	ID         uint             `json:"id"`
	Name       string           `json:"name"`
	Barcode    string           `json:"barcode"`
	Price      int              `json:"price"`
	About      int              `json:"about"`
	CategoryID uint             `json:"category_id"`
	Thumbnail  string           `json:"thumnail"`
	IsPopular  int              `json:"is_popular"`
	Category   CategoryResponse `json:"category"`
}

type GetAllProductResponse struct {
	Products []ProductResponse            `json:"categories"`
	Pagination pagination.PaginationResponse `json:"pagination"`
}
