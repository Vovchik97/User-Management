package handlers

import "userManagement/internal/models"

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
	RoleName string `json:"role_name" binding:"required"`
}

type UserInfo struct {
	ID     uint        `json:"id"`
	RoleID uint        `json:"role_id"`
	Role   models.Role `json:"role"`
}
