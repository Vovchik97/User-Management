package handlers

import (
	"net/http"
	"userManagement/internal/config"
	"userManagement/internal/models"

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

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	config.DB.Model(&user).Updates(models.User{
		Name:  input.Name,
		Email: input.Email,
	})

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

	user.Role = input.Role
	config.DB.Save(&user)

	c.JSON(http.StatusOK, user)
}
