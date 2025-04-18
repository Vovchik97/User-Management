package handlers

// ResponseError - структура для ответа с ошибкой
type ResponseError struct {
	Message string `json:"message"`
}
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
	Role string `json:"role" binding:"required,oneof=admin user"`
}
