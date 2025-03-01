package users

// UsersController Controllers 用户控制器聚合结构体
type UsersController struct{}

// 提供控制器实例（单例模式）
var usersController = &UsersController{}
