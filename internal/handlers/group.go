package handlers

import (
	"net/http"
	"userManagement/internal/config"
	"userManagement/internal/models"

	"github.com/gin-gonic/gin"
)

// CreateGroups godoc
// @Summary Создание новой группы
// @Tags groups
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
	c.JSON(http.StatusCreated, group)
}

// GetGroups godoc
// @Summary Получение списка всех групп
// @Tags groups
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
// @Tags groups
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

	group.Name = input.Name
	config.DB.Save(&group)
	c.JSON(http.StatusOK, group)
}

// DeleteGroup godoc
// @Summary Удаление группы
// @Tags groups
// @Produce json
// @Param id path int true "ID группы"
// @Success 200 {object} ResponseMessage
// @Failure 403 {object} ResponseError
// @Router /groups/{id} [delete]
// @Security UserID
func DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Unscoped().Delete(&models.Group{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Группа не найдена"})
		return
	}
	c.JSON(http.StatusOK, ResponseError{Message: "Группа успешно удалена"})
}

// AddUserToGroup godoc
// @Summary Добавление пользователя в группу
// @Tags groups
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

	c.JSON(http.StatusOK, ResponseError{Message: "Пользователь добавлен в группу"})
}

// RemoveUserFromGroup godoc
// @Summary Удаление пользователя из группы
// @Tags groups
// @Produce json
// @Param id path int true "ID группы"
// @Param user_id path int true "ID пользователя"
// @Success 200 {object} ResponseMessage
// @Failure 400 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Router /groups/{id}/users/{user_id} [delete]
// @Security UserID
func RemoveUserFromGroup(c *gin.Context) {
	groupID := c.Param("id")
	userID := c.Param("user_id")

	var group models.Group
	if err := config.DB.Preload("Users").First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Группа не найдена"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Пользователь не найден"})
		return
	}

	// Проверяем, состоит ли пользователь в группе
	var groupUser models.Group
	if err := config.DB.Where("group_id = ? AND user_id = ?", group.ID, user.ID).First(&groupUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: "Пользователь не состоит в группе"})
		return
	}

	if err := config.DB.Model(&group).Association("Users").Unscoped().Delete(&user); err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Не удалось удалить пользователя из группы"})
		return
	}

	c.JSON(http.StatusOK, ResponseError{Message: "Пользователь удален из группы"})
}
