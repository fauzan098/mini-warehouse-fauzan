package request

type CreateCategoryRequest struct {
	Name    string `json:"name" validate:"required"`
	Tagline string `json:"tagline" validate:"required"`
	Photo   string `json:"photo" validate:"required"`
}

type GetAllCategoryRequest struct {
	page      int `query:"page"`
	Limit     int `query:"limit"`
	Search    int `query:"search"`
	SortBy    int `query:"sort_by"`
	SortOrder int `query:"sort_order"`
}
