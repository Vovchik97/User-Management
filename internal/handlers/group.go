package handlers

import (
	"fmt"
	"net/http"
	"userManagement/internal/config"
	"userManagement/internal/models"
	"userManagement/internal/utils"

	"github.com/gin-gonic/gin"
)

// CreateGroups godoc
// @Summary Создание новой группы
// @Tags Groups
// @Accept json
// @Produce json
// @Param group body GroupInput true "Название группы"
// @Success 201 {object} ResponseMessage
// @Failure 400 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Router /groups [post]
// @Security UserID
func CreateGroups(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	if err := config.DB.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Не удалось создать группу"})
		return
	}

	if userID, exists := c.Get("userID"); exists {
		utils.LogAction(userID.(uint), fmt.Sprintf("Создана группа: %s", group.Name))
	}

	c.JSON(http.StatusCreated, group)
}

// GetGroups godoc
// @Summary Получение списка всех групп
// @Tags Groups
// @Produce json
// @Success 200 {array} models.Group
// @Failure 403 {object} ResponseError
// @Router /groups [get]
// @Security UserID
func GetGroups(c *gin.Context) {
	var groups []models.Group

	// Загружаем все группы, включая пользователей
	if err := config.DB.Preload("Users").Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Ошибка при загрузке групп"})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// UpdateGroup godoc
// @Summary Обновление названия группы
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path int true "ID группы"
// @Param group body GroupInput true "Новое название"
// @Success 200 {object} ResponseMessage
// @Failure 400 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Router /groups/{id} [put]
// @Security UserID
func UpdateGroup(c *gin.Context) {
	id := c.Param("id")
	var group models.Group
	if err := config.DB.First(&group, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Группа не найдена"})
		return
	}

	var input models.Group
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	oldName := group.Name
	group.Name = input.Name
	config.DB.Save(&group)

	if userID, exists := c.Get("userID"); exists {
		utils.LogAction(userID.(uint), fmt.Sprintf("Обновлена группа: %s -> %s", oldName, group.Name))
	}

	c.JSON(http.StatusOK, group)
}

// DeleteGroup godoc
// @Summary Удаление группы
// @Tags Groups
// @Produce json
// @Param id path int true "ID группы"
// @Success 200 {object} ResponseMessage
// @Failure 403 {object} ResponseError
// @Router /groups/{id} [delete]
// @Security UserID
func DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	var group models.Group
	if err := config.DB.First(&group, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Группа не найдена"})
		return
	}

	config.DB.Unscoped().Delete(&group)

	if userID, exists := c.Get("userID"); exists {
		utils.LogAction(userID.(uint), fmt.Sprintf("Удалил группу %s", group.Name))
	}

	c.JSON(http.StatusOK, ResponseError{Message: "Группа успешно удалена"})
}

// AddUserToGroup godoc
// @Summary Добавление пользователя в группу
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path int true "ID группы"
// @Param user body UserGroupInput true "ID пользователя для добавления"
// @Success 200 {object} ResponseMessage
// @Failure 400 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Router /groups/{id}/users [post]
// @Security UserID
func AddUserToGroup(c *gin.Context) {
	groupID := c.Param("id")
	var input struct {
		UserID uint `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	var group models.Group
	if err := config.DB.Preload("Users").First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Группа не найдена"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, input.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Пользователь не найден"})
		return
	}

	if err := config.DB.Model(&group).Association("Users").Append(&user); err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Не удалось добавить пользователя в группу"})
		return
	}

	if userID, exists := c.Get("userID"); exists {
		utils.LogAction(userID.(uint), fmt.Sprintf("Добавлен пользователь %s в группу %s", user.Name, group.Name))
	}

	c.JSON(http.StatusOK, ResponseError{Message: "Пользователь добавлен в группу"})
}

// RemoveUserFromGroup godoc
// @Summary Удаление пользователя из группы
// @Tags Groups
// @Produce json
// @Param id path int true "ID группы"
// @Param user_id path int true "ID пользователя"
// @Success 200 {object} ResponseMessage
// @Failure 400 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Router /groups/{id}/users/{user_id} [delete]
// @Security UserID
func RemoveUserFromGroup(c *gin.Context) {
	groupId := c.Param("id")
	userId := c.Param("user_id")

	var group models.Group
	if err := config.DB.Preload("Users").First(&group, groupId).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Группа не найдена"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Пользователь не найден"})
		return
	}

	// Проверка: состоит ли пользователь в группе
	found := false
	for _, u := range group.Users {
		if u.ID == user.ID {
			found = true
			break
		}
	}
	if !found {
		c.JSON(http.StatusBadRequest, ResponseError{Message: "Пользователь не состоит в группе"})
		return
	}

	if err := config.DB.Model(&group).Association("Users").Unscoped().Delete(&user); err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Не удалось удалить пользователя из группы"})
		return
	}

	if userID, exists := c.Get("userID"); exists {
		utils.LogAction(userID.(uint), fmt.Sprintf("Удален пользователь %s из группы %s", user.Name, group.Name))
	}

	c.JSON(http.StatusOK, ResponseError{Message: "Пользователь удален из группы"})
}
