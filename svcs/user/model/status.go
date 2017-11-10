package model

// UserStatus 用户状态
type UserStatus int

const (
	// UserStatusUnknown 未知状态
	UserStatusUnknown UserStatus = iota
	// UserStatusCreated 新建正常
	UserStatusCreated
	// UserStatusLocked 锁定
	UserStatusLocked
)

// UserAuthority 用户角色
type UserAuthority int

const (
	// UserAuthorityUnknown 未知
	UserAuthorityUnknown UserAuthority = iota
	// UserAuthorityCust 用户
	UserAuthorityCust
	// UserAuthorityTenant 租户/供应商
	UserAuthorityTenant
	// UserAuthorityAdmin 管理员
	UserAuthorityAdmin
)
