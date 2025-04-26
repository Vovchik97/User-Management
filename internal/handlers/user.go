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
// @Security BearerAuth
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

	input.Sanitize()

	// Хешируем пароль
	hashedPassword, errPassword := utils.HashPassword(input.Password)
	if errPassword != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Ошибка при хешировании пароля"})
		return
	}

	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		RoleID:       3,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: "Не удалось создать пользователя"})
		return
	}

	// Загружаем роль для возвращаемого пользователя
	config.DB.Preload("Role").First(&user, user.ID)

	if currentUser, exists := c.Get("currentUser"); exists {
		userInfo := currentUser.(UserInfo)
		utils.LogAction(userInfo.ID, fmt.Sprintf("Создал пользователя: %s", user.Name))
	}

	c.JSON(http.StatusCreated, user)
}

// GetUsers godoc
// @Summary Получение списка пользователей
// @Description Получение списка всех пользователей или отфильтрованных по роли.
// @Tags Users
// @Security BearerAuth
// @Produce  json
// @Param role query string false "Роль пользователя"
// @Success 200 {array} models.User
// @Failure 500 {object} ResponseError "Ошибка при получении списка пользователей"
// @Router /users [get]
func GetUsers(c *gin.Context) {
	_, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, ResponseError{Message: "Необходима авторизация"})
		return
	}

	roleName := c.Query("role")
	var users []models.User

	query := config.DB
	if roleName != "" {
		// Используем внешний ключ RoleID, а не строку
		var role models.Role
		if err := config.DB.Where("name = ?", roleName).First(&role).Error; err == nil {
			query = query.Where("role_id = ?", role.ID)
		} else {
			c.JSON(http.StatusBadRequest, ResponseError{Message: "Роль не найдена"})
			return
		}
	}

	if err := query.Preload("Role").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Ошибка при получении списка пользователей"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser godoc
// @Summary Обновление пользователя
// @Description Обновление данных пользователя по его ID.
// @Tags Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "ID пользователя"
// @Param input body UpdateUserInput true "Параметры обновления пользователя"
// @Success 200 {object} models.User
// @Failure 400 {object} ResponseError "Ошибка при обновлении пользователя"
// @Failure 404 {object} ResponseError "Пользователь не найден"
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	currentUserRaw, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, ResponseError{Message: "Необходима авторизация"})
		return
	}
	userInfo := currentUserRaw.(UserInfo)
	currentUser := models.User{
		ID:   userInfo.ID,
		Role: userInfo.Role,
	}

	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Пользователь не найден"})
		return
	}

	// Проверка прав доступа:
	// Если текущий пользователь не является администратором или модератором, он не может редактировать чужие данные
	if currentUser.ID != user.ID && currentUser.Role.Name != "admin" && currentUser.Role.Name != "moderator" {
		c.JSON(http.StatusForbidden, ResponseError{Message: "Недостаточно прав для реадктирования других пользователей"})
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	input.Sanitize()

	// Обновляем данные пользователя
	if err := config.DB.Model(&user).Updates(models.User{
		Name:  input.Name,
		Email: input.Email,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Не удалось обновить данные пользователя"})
		return
	}

	utils.LogAction(currentUser.ID, fmt.Sprintf("Обновил пользователя: %s", user.Name))

	// Подгружаем роль перед возвратом
	config.DB.Preload("Role").First(&user, user.ID)

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Удаление пользователя
// @Description Удаление пользователя по его ID.
// @Tags Users
// @Security BearerAuth
// @Param id path int true "ID пользователя"
// @Success 200 {object} ResponseError "Пользователь удален"
// @Failure 404 {object} ResponseError "Пользователь не найден"
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	currentUserRaw, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, ResponseError{Message: "Необходима авторизация"})
		return
	}
	userInfo := currentUserRaw.(UserInfo)
	currentUser := models.User{
		ID:   userInfo.ID,
		Role: userInfo.Role,
	}

	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Пользователь не найден"})
		return
	}

	config.DB.Unscoped().Delete(&user)

	utils.LogAction(currentUser.ID, fmt.Sprintf("Удалил пользователя: %s", user.Name))

	c.JSON(http.StatusOK, ResponseError{Message: "Пользователь удален"})
}

// UpdateUserRole godoc
// @Summary Назначить роль пользователю
// @Description Обновление роли пользователя по ID
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param input body UpdateUserRoleInput true "Новая роль"
// @Success 200 {object} models.User
// @Failure 400 {object} ResponseError "Неверный ввод"
// @Failure 404 {object} ResponseError "Пользователь не найден"
// @Router /users/{id}/role [patch]
func UpdateUserRole(c *gin.Context) {
	currentUserRaw, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, ResponseError{Message: "Необходима авторизация"})
		return
	}
	userInfo := currentUserRaw.(UserInfo)

	var user models.User
	if err := config.DB.Preload("Role").First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, ResponseError{Message: "Пользователь не найден"})
		return
	}

	var input UpdateUserRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	// Найти роль по имени
	var role models.Role
	if err := config.DB.Where("name = ?", input.RoleName).First(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: "Указанная роль не найдена"})
		return
	}

	oldRole := user.Role.Name
	user.RoleID = role.ID
	user.Role = role

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Не удалось обновить роль пользователя"})
		return
	}

	utils.LogAction(userInfo.ID, fmt.Sprintf("Обновил роль пользователя %s: %s -> %s", user.Name, oldRole, user.Role.Name))

	// Подгружаем новую роль перед возвратом
	config.DB.Preload("Role").First(&user, user.ID)

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
// @Security BearerAuth
func BanUser(c *gin.Context) {
	currentUserRaw, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, ResponseError{Message: "Необходима авторизация"})
		return
	}
	userInfo := currentUserRaw.(UserInfo)
	currentUser := models.User{
		ID:   userInfo.ID,
		Role: userInfo.Role,
	}

	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
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

	utils.LogAction(currentUser.ID, fmt.Sprintf("Заблокировал пользователя: %s", user.Name))

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
// @Security BearerAuth
func UnbanUser(c *gin.Context) {
	currentUserRaw, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, ResponseError{Message: "Необходима авторизация"})
		return
	}
	userInfo := currentUserRaw.(UserInfo)
	currentUser := models.User{
		ID:   userInfo.ID,
		Role: userInfo.Role,
	}

	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
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
	utils.LogAction(currentUser.ID, fmt.Sprintf("Разблокировал пользователя: %s", user.Name))

	c.JSON(http.StatusOK, ResponseMessage{Message: "Пользователь разблокирован"})
}

// GetProfile godoc
// @Summary Получение профиля текущего пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} models.User "Профиль пользователя"
// @Failure 401 {object} ResponseError "Неавторизованный доступ"
// @Failure 404 {object} ResponseError "Пользователь не найден"
// @Failure 500 {object} ResponseError "Ошибка при получении профиля пользователя"
// @Router /users/me [get]
// @Security BearerAuth
func GetProfile(c *gin.Context) {
	currentUserRaw, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, ResponseError{Message: "Пользователь не найден"})
		return
	}

	userInfo := currentUserRaw.(UserInfo)

	var user models.User
	if err := config.DB.Preload("Role").First(&user, userInfo.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Ошибка при получении профиля пользователя"})
		return
	}

	c.JSON(http.StatusOK, user)
}
