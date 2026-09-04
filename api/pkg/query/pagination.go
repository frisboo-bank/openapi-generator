package query

import "fmt"

type Pagination struct {
	Page int64 `json:"page"`
	Size int64 `json:"size"`
}

func NewPagination(page int64, size int64) *Pagination {
	return &Pagination{
		Page: page,
		Size: size,
	}
}

func (p *Pagination) Validate() error {
	if p.Size <= 0 {
		return fmt.Errorf("pagination size must be positive, got %d", p.Size)
	}
	if p.Page <= 0 {
		return fmt.Errorf("pagination page must be positive, got %d", p.Page)
	}

	return nil
}

func (p *Pagination) Offset() int64 {
	return (p.Page - 1) * p.Size
}

func (p *Pagination) TotalPages(totalItems int64) int64 {
	if p.Size <= 0 {
		return 0
	}
	return (totalItems + p.Size - 1) / p.Size
}
