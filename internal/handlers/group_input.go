package handlers

// GroupInput используется для создания или обновления группы
type GroupInput struct {
	Name string `json:"name" binding:"required"`
}

// UserGroupInput используется для добавления пользователя в группу
type UserGroupInput struct {
	UserID uint `json:"user_id" binding:"required"`
}
