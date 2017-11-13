package model

// ProductStatus 商品状态
type ProductStatus int

const (
	// ProductStatusUnknown 未知状态
	ProductStatusUnknown ProductStatus = iota
	// ProductStatusNormal 上架
	ProductStatusNormal
	// ProductStatusLocked 下架
	ProductStatusLocked
	// ProductStatusViolate 违反商品，暂不用
	ProductStatusViolate
)
