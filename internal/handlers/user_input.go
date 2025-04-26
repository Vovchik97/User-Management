package handlers

import (
	"userManagement/internal/models"
	"userManagement/internal/utils"
)

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

type UpdateUserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"omitempty,email"`
}

func (i *UpdateUserInput) Sanitize() {
	i.Name = utils.SanitizeInput(i.Name)
	i.Email = utils.SanitizeInput(i.Email)
}

type UpdateUserRoleInput struct {
	RoleName string `json:"role_name" binding:"required"`
}

type UserInfo struct {
	ID     uint        `json:"id"`
	RoleID uint        `json:"role_id"`
	Role   models.Role `json:"role"`
}
