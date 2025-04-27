package dto

import "userManagement/internal/models"

// UserInfo используется для передачи информации о пользователе
type UserInfo struct {
	ID     uint         `json:"id"`
	RoleID uint         `json:"role_id"`
	Role   *models.Role `json:"role"`
}
