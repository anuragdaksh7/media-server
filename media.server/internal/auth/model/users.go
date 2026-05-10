package model

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	Base
	Email    string `gorm:"uniqueIndex"`
	Password string
	Role     UserRole `gorm:"default:user"`
}
