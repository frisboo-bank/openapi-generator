package query

type ListQuery struct {
	Pagination *Pagination    `json:"pagination"`
	Filters    *QueryFilters  `json:"filters"`
	Orders     *QueryOrders   `json:"orders"`
	Searches   *QuerySearches `json:"searches"`
}
