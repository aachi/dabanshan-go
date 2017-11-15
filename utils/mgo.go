package utils

// Pagination ...
type Pagination struct {
	Count     int
	PageIndex int
	PageSize  int
	Sortor    []string
	Data      interface{}
}
