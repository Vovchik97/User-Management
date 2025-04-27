package dto

import "userManagement/internal/utils"

// GroupInput используется для создания или обновления группы
type GroupInput struct {
	Name string `json:"name" binding:"required"`
}

func (i *GroupInput) Sanitize() {
	i.Name = utils.SanitizeInput(i.Name)
}

// UserGroupInput используется для добавления пользователя в группу
type UserGroupInput struct {
	UserID uint `json:"user_id" binding:"required"`
}
