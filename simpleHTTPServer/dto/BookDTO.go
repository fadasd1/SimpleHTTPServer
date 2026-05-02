package dto

type CreateBookRequest struct {
	Title    string
	AuthorID int
	Price    float64
	Stock    int
}

type BookFilter struct {
	Page     string
	Limit    string
	AuthorID string
	MinPrice string
	MaxPrice string
	InStock  string
}
