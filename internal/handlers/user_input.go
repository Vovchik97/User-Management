package handlers

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"omitempty,email"`
}

type UpdateUserRoleInput struct {
	Role string `json:"role" binding:"required,oneof=admin moderator user"`
}

type UserInfo struct {
	ID   uint
	Role string
}
