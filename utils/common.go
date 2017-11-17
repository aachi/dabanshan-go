package utils

// Pagination ...
type Pagination struct {
	Count     int
	PageIndex int
	PageSize  int
	Sortor    []string
	Data      interface{}
}

// Image struct
type Image struct {
	Filepath []byte
	Md5      string
}

// Images muti imgs.
type Images []*Image
