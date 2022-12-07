package pagination

type Pagination struct {
	Limit int  `json:"limit"`
	Page  int  `json:"page"`
	Asc   bool `json:"asc"`
}
