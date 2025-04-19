package handlers

import (
	"fmt"
	"net/http"
	"userManagement/internal/config"
	"userManagement/internal/models"
	"userManagement/internal/utils"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary Создание нового пользователя
// @Description Создание нового пользователя с указанием имени, email, пароля и роли.
// @Tags Users
// @Security UserID
// @Accept  json
// @Produce  json
// @Param input body CreateUserInput true "Параметры пользователя"
// @Success 201 {object} models.User
// @Failure 400 {object} ResponseError "Ошибка при создании пользователя"
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: "Не удалось создать пользователя"})
		return
	}

	if userID, exists := c.Get("userID"); exists {
		utils.LogAction(userID.(uint), fmt.Sprintf("Создан пользователь: %s", user.Name))
	}

	c.JSON(http.StatusCreated, user)
}

// GetUsers godoc
// @Summary Получение списка пользователей
// @Description Получение списка всех пользователей или отфильтрованных по роли.
// @Tags Users
// @Security UserID
// @Produce  json
// @Param role query string false "Роль пользователя"
// @Success 200 {array} models.User
// @Failure 500 {object} ResponseError "Ошибка при получении списка пользователей"
// @Router /users [get]
func GetUsers(c *gin.Context) {
	role := c.Query("role")
	var users []models.User

	query := config.DB
	if role != "" {
		query = query.Where("role = ?", role)
	}

	query.Find(&users)

	c.JSON(http.StatusOK, users)
}

// UpdateUser godoc
// @Summary Обновление пользователя
// @Description Обновление данных пользователя по его ID.
// @Tags Users
// @Security UserID
// @Accept  json
// @Produce  json
// @Param id path int true "ID пользователя"
// @Param input body UpdateUserInput true "Параметры обновления пользователя"
// @Success 200 {object} models.User
// @Failure 400 {object} ResponseError "Ошибка при обновлении пользователя"
// @Failure 404 {object} ResponseError "Пользователь не найден"
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Пользователь не найден"})
		return
	}

	// Получаем текущего пользователя из контекста
	currentUserRaw, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, ResponseError{Message: "Необходима авторизация"})
		return
	}
	currentUser := currentUserRaw.(models.User)

	// Проверка прав доступа:
	// Если текущий пользователь не является администратором или модератором, он не может редактировать чужие данные
	if currentUser.ID != user.ID && currentUser.Role != "admin" && currentUser.Role != "moderator" {
		c.JSON(http.StatusForbidden, ResponseError{Message: "Недостаточно прав для реадктирования других пользователей"})
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	// Обновляем данные пользователя
	if err := config.DB.Model(&user).Updates(models.User{
		Name:  input.Name,
		Email: input.Email,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Не удалось обновить данные пользователя"})
		return
	}

	if userID, exists := c.Get("userID"); exists {
		utils.LogAction(userID.(uint), fmt.Sprintf("Обновлен пользователь: %s", user.Name))
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Удаление пользователя
// @Description Удаление пользователя по его ID.
// @Tags Users
// @Security UserID
// @Param id path int true "ID пользователя"
// @Success 200 {object} ResponseError "Пользователь удален"
// @Failure 404 {object} ResponseError "Пользователь не найден"
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Пользователь не найден"})
		return
	}

	config.DB.Unscoped().Delete(&user)

	if userId, exists := c.Get("userID"); exists {
		utils.LogAction(userId.(uint), fmt.Sprintf("Удален пользователь: %s", user.Name))
	}

	c.JSON(http.StatusOK, ResponseError{Message: "Пользователь удален"})
}

// UpdateUserRole godoc
// @Summary Назначить роль пользователю
// @Description Обновление роли пользователя по ID
// @Tags Users
// @Security UserID
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param input body UpdateUserRoleInput true "Новая роль"
// @Success 200 {object} models.User
// @Failure 400 {object} ResponseError "Неверный ввод"
// @Failure 404 {object} ResponseError "Пользователь не найден"
// @Router /users/{id}/role [patch]
func UpdateUserRole(c *gin.Context) {
	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Пользователь не найден"})
		return
	}

	var input UpdateUserRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	oldRole := user.Role
	user.Role = input.Role
	config.DB.Save(&user)

	if userID, exists := c.Get("userID"); exists {
		utils.LogAction(userID.(uint), fmt.Sprintf("Обновлена роль пользователя %s: %s -> %s", user.Name, oldRole, user.Role))
	}

	c.JSON(http.StatusOK, user)
}

// BanUser godoc
// @Summary Временная блокировка пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} ResponseMessage
// @Failure 400 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Router /users/{id}/ban [patch]
// @Security UserID
func BanUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Пользователь не найден"})
		return
	}

	// Проверяяем, что пользователь не заблокирован
	if user.IsBanned {
		c.JSON(http.StatusBadRequest, ResponseError{Message: "Пользователь уже заблокирован"})
		return
	}

	user.IsBanned = true
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Не удалось заблокировать пользователя"})
		return
	}

	if adminID, exists := c.Get("userID"); exists {
		utils.LogAction(adminID.(uint), fmt.Sprintf("Заблокировал пользователя %s", user.Name))
	}

	c.JSON(http.StatusOK, ResponseMessage{Message: "Пользователь заблокирован"})
}

// UnbanUser godoc
// @Summary Разблокировка пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} ResponseMessage
// @Failure 400 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Router /users/{id}/unban [patch]
// @Security UserID
func UnbanUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Пользователь не найден"})
		return
	}

	// Проверяем, что пользователь не заблокирован
	if !user.IsBanned {
		c.JSON(http.StatusBadRequest, ResponseError{Message: "Пользователь не заблокирован"})
		return
	}

	user.IsBanned = false
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Не удалось разблокировать пользователя"})
		return
	}

	// Логируем действие
	if adminID, exists := c.Get("userID"); exists {
		utils.LogAction(adminID.(uint), fmt.Sprintf("Разблокировал пользователя %s", user.Name))
	}

	c.JSON(http.StatusOK, ResponseMessage{Message: "Пользователь разблокирован"})
}
