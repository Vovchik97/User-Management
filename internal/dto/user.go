package dto

import (
	"userManagement/internal/utils"
)

// CreateUserInput используется для создания нового пользователя
type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (i *CreateUserInput) Sanitize() {
	i.Name = utils.SanitizeInput(i.Name)
	i.Email = utils.SanitizeInput(i.Email)
	i.Password = utils.SanitizeInput(i.Password)
}

// UpdateUserInput используется для обновления информации о пользователе
type UpdateUserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"omitempty,email"`
}

func (i *UpdateUserInput) Sanitize() {
	i.Name = utils.SanitizeInput(i.Name)
	i.Email = utils.SanitizeInput(i.Email)
}

// UpdateUserRoleInput используется для обновления роли пользователя
type UpdateUserRoleInput struct {
	RoleName string `json:"role_name" binding:"required"`
}
