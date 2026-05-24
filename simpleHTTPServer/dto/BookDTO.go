package dto

type BookFilter struct {
	Page     string `json:"page"`
	Limit    string `json:"limit"`
	AuthorID string `json:"author_id"`
	MinPrice string `json:"min_price"`
	MaxPrice string `json:"max_price"`
	InStock  string `json:"in_stock"`
}

type CreateBookRequest struct {
	Title    string  `json:"title"`
	AuthorID int     `json:"author_id"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
}
