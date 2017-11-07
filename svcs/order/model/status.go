package model

// OrderStatus 订单状态
type OrderStatus int

const (
	// OrderStatusUnknown 未知状态
	OrderStatusUnknown OrderStatus = iota
	// OrderStatusCreated 新建，待付款
	OrderStatusCreated
	// OrderStatusPaymented 已付款
	OrderStatusPaymented
	// OrderStatusDispatched 已发货
	OrderStatusDispatched
	// OrderStatusFinished 完成
	OrderStatusFinished
	// OrderStatusCanceled 关闭
	OrderStatusCanceled
)
