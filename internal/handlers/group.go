package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"userManagement/internal/config"
	"userManagement/internal/dto"
	"userManagement/internal/models"
	"userManagement/internal/services"
	"userManagement/internal/utils"
)

// CreateGroups godoc
// @Summary Создание новой группы
// @Tags Groups
// @Accept json
// @Produce json
// @Param group body dto.GroupInput true "Название группы"
// @Success 201 {object} dto.ResponseMessage
// @Failure 400 {object} dto.ResponseError
// @Failure 403 {object} dto.ResponseError
// @Router /groups [post]
// @Security BearerAuth
func CreateGroups(c *gin.Context) {
	var input dto.GroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warn("Некорректный ввод при создании группы:", err)
		c.JSON(http.StatusBadRequest, dto.ResponseError{Message: err.Error()})
		return
	}

	input.Sanitize()

	group := models.Group{
		Name: input.Name,
	}

	if err := config.DB.Create(&group).Error; err != nil {
		utils.Log.Error("Ошибка при создании группы:", err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Message: "Не удалось создать группу"})
		return
	}

	utils.Log.Infof("Создана группа: %s", group.Name)

	if userID, exists := c.Get("userID"); exists {
		err := services.LogAction(userID.(uint), fmt.Sprintf("Создана группа: %s", group.Name))
		if err != nil {
			utils.Log.Warn("Ошибка при логировании действия:", err)
		}
	}

	c.JSON(http.StatusCreated, group)
}

// GetGroups godoc
// @Summary Получение списка всех групп
// @Tags Groups
// @Produce json
// @Success 200 {array} models.Group
// @Failure 403 {object} dto.ResponseError
// @Router /groups [get]
// @Security BearerAuth
func GetGroups(c *gin.Context) {
	var groups []models.Group

	// Загружаем все группы, включая пользователей
	if err := config.DB.Preload("Users.Role").Find(&groups).Error; err != nil {
		utils.Log.Error("Ошибка при получении списка групп:", err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Message: "Ошибка при загрузке групп"})
		return
	}

	utils.Log.Info("Получен список групп")
	c.JSON(http.StatusOK, groups)
}

// UpdateGroup godoc
// @Summary Обновление названия группы
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path int true "ID группы"
// @Param group body dto.GroupInput true "Новое название"
// @Success 200 {object} dto.ResponseMessage
// @Failure 400 {object} dto.ResponseError
// @Failure 403 {object} dto.ResponseError
// @Router /groups/{id} [put]
// @Security BearerAuth
func UpdateGroup(c *gin.Context) {
	id := c.Param("id")
	var group models.Group
	if err := config.DB.First(&group, id).Error; err != nil {
		utils.Log.Warnf("Группа с ID %s не найдена", id)
		c.JSON(http.StatusNotFound, dto.ResponseError{Message: "Группа не найдена"})
		return
	}

	var input dto.GroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warn("Некорректный ввод при обновлении группы:", err)
		c.JSON(http.StatusBadRequest, dto.ResponseError{Message: err.Error()})
		return
	}

	input.Sanitize()

	oldName := group.Name
	group.Name = input.Name
	config.DB.Save(&group)

	utils.Log.Infof("Обновлена группа %s -> %s", oldName, group.Name)

	if userID, exists := c.Get("userID"); exists {
		err := services.LogAction(userID.(uint), fmt.Sprintf("Обновлена группа: %s -> %s", oldName, group.Name))
		if err != nil {
			utils.Log.Warn("Ошибка при логировании действия:", err)
		}
	}

	c.JSON(http.StatusOK, group)
}

// DeleteGroup godoc
// @Summary Удаление группы
// @Tags Groups
// @Produce json
// @Param id path int true "ID группы"
// @Success 200 {object} dto.ResponseMessage
// @Failure 403 {object} dto.ResponseError
// @Router /groups/{id} [delete]
// @Security BearerAuth
func DeleteGroup(c *gin.Context) {
	id := c.Param("id")

	var group models.Group
	if err := config.DB.First(&group, id).Error; err != nil {
		utils.Log.Warnf("Попытка удалить несуществующую группу с ID %s", id)
		c.JSON(http.StatusNotFound, dto.ResponseError{Message: "Группа не найдена"})
		return
	}

	// Удалить связи с пользователями
	if err := config.DB.Model(&group).Association("Users").Clear(); err != nil {
		utils.Log.Error("Ошибка при удалении связей с пользователями:", err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Message: "Не удалось удалить связи с пользователями"})
		return
	}

	// Удалить группу
	if err := config.DB.Unscoped().Delete(&group).Error; err != nil {
		utils.Log.Error("Ошибка при удалении группы:", err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Message: "Не удалось удалить группу"})
		return
	}

	utils.Log.Infof("Удалена группа %s", group.Name)

	if userID, exists := c.Get("userID"); exists {
		err := services.LogAction(userID.(uint), fmt.Sprintf("Удалил группу %s", group.Name))
		if err != nil {
			utils.Log.Warn("Ошибка при логировании действия:", err)
		}
	}

	c.JSON(http.StatusOK, dto.ResponseError{Message: "Группа успешно удалена"})
}

// AddUserToGroup godoc
// @Summary Добавление пользователя в группу
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path int true "ID группы"
// @Param user body dto.UserGroupInput true "ID пользователя для добавления"
// @Success 200 {object} dto.ResponseMessage
// @Failure 400 {object} dto.ResponseError
// @Failure 403 {object} dto.ResponseError
// @Router /groups/{id}/users [post]
// @Security BearerAuth
func AddUserToGroup(c *gin.Context) {
	groupID := c.Param("id")
	var input struct {
		UserID uint `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warn("Некорректный ввод при добавлении пользователя в группу:", err)
		c.JSON(http.StatusBadRequest, dto.ResponseError{Message: err.Error()})
		return
	}

	var group models.Group
	if err := config.DB.Preload("Users").First(&group, groupID).Error; err != nil {
		utils.Log.Warnf("Группа с ID %s не найдена", groupID)
		c.JSON(http.StatusNotFound, dto.ResponseError{Message: "Группа не найдена"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, input.UserID).Error; err != nil {
		utils.Log.Warnf("Пользователь с ID %d не найден", input.UserID)
		c.JSON(http.StatusNotFound, dto.ResponseError{Message: "Пользователь не найден"})
		return
	}

	if err := config.DB.Model(&group).Association("Users").Append(&user); err != nil {
		utils.Log.Error("Ошибка при добавлении пользователя в группу:", err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Message: "Не удалось добавить пользователя в группу"})
		return
	}

	utils.Log.Infof("Добавлен пользователь %s в группу %s", user.Name, group.Name)

	if userID, exists := c.Get("userID"); exists {
		err := services.LogAction(userID.(uint), fmt.Sprintf("Добавлен пользователь %s в группу %s", user.Name, group.Name))
		if err != nil {
			utils.Log.Warn("Ошибка при логировании действия:", err)
		}
	}

	c.JSON(http.StatusOK, dto.ResponseError{Message: "Пользователь добавлен в группу"})
}

// RemoveUserFromGroup godoc
// @Summary Удаление пользователя из группы
// @Tags Groups
// @Produce json
// @Param id path int true "ID группы"
// @Param user_id path int true "ID пользователя"
// @Success 200 {object} dto.ResponseMessage
// @Failure 400 {object} dto.ResponseError
// @Failure 403 {object} dto.ResponseError
// @Router /groups/{id}/users/{user_id} [delete]
// @Security BearerAuth
func RemoveUserFromGroup(c *gin.Context) {
	groupId := c.Param("id")
	userId := c.Param("user_id")

	var group models.Group
	if err := config.DB.Preload("Users").First(&group, groupId).Error; err != nil {
		utils.Log.Warnf("Группа с ID %s не найдена", groupId)
		c.JSON(http.StatusNotFound, dto.ResponseError{Message: "Группа не найдена"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userId).Error; err != nil {
		utils.Log.Warnf("Пользователь с ID %s не найден", userId)
		c.JSON(http.StatusNotFound, dto.ResponseError{Message: "Пользователь не найден"})
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
		utils.Log.Warn("Попытка удалить пользователя, которого нет в группе")
		c.JSON(http.StatusBadRequest, dto.ResponseError{Message: "Пользователь не состоит в группе"})
		return
	}

	if err := config.DB.Model(&group).Association("Users").Unscoped().Delete(&user); err != nil {
		utils.Log.Error("Ошибка при удалении пользователя из группы:", err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Message: "Не удалось удалить пользователя из группы"})
		return
	}

	utils.Log.Infof("Удален пользователь %s из группы %s", user.Name, group.Name)

	if userID, exists := c.Get("userID"); exists {
		err := services.LogAction(userID.(uint), fmt.Sprintf("Удален пользователь %s из группы %s", user.Name, group.Name))
		if err != nil {
			utils.Log.Warn("Ошибка при логировании действия:", err)
		}
	}

	c.JSON(http.StatusOK, dto.ResponseError{Message: "Пользователь удален из группы"})
}
