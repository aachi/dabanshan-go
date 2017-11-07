package model

// UserStatus 用户状态
type UserStatus int

const (
	// UserStatusUnknown 未知状态
	UserStatusUnknown UserStatus = iota
	// OrderStatusCreated 新建正常
	UserStatusCreated
	// UserStatusLocked 锁定
	UserStatusLocked
)
