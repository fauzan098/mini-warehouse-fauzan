package request

type CreateProductRequest struct {
	Name       string `json:"name" validate:"required"`
	Barcode    string `json:"barcode" validate:"required"`
	Price      int    `json:"price" validate:"required"`
	About      string `json:"about" validate:"required"`
	CategoryID uint   `json:"category_id" validate:"required"`
	Thumbnail  string `json:"thumbnail" validate:"required"`
	IsPopular  string `json:"is_popular" validate:"required"`
}

type GetAllProductRequest struct {
	page      int `query:"page"`
	Limit     int `query:"limit"`
	Search    int `query:"search"`
	SortBy    int `query:"sort_by"`
	SortOrder int `query:"sort_order"`
}
