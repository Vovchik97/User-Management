package handlers

import (
	"fmt"
	"net/http"
	"userManagement/internal/config"
	"userManagement/internal/dto"
	"userManagement/internal/models"
	"userManagement/internal/services"
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
// @Param input body dto.CreateUserInput true "Параметры пользователя"
// @Success 201 {object} models.User
// @Failure 400 {object} dto.ResponseError "Ошибка при создании пользователя"
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var input dto.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("Некорректный ввод при создании пользователя: %v", err)
		c.JSON(http.StatusBadRequest, dto.ResponseError{Message: err.Error()})
		return
	}

	input.Sanitize()

	// Хешируем пароль
	hashedPassword, errPassword := utils.HashPassword(input.Password)
	if errPassword != nil {
		utils.Log.Errorf("Ошибка хеширования пароля: %v", errPassword)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Message: "Ошибка при хешировании пароля"})
		return
	}

	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		RoleID:       3,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.Log.Errorf("Не удалось создать пользователя: %v", err)
		c.JSON(http.StatusBadRequest, dto.ResponseError{Message: "Не удалось создать пользователя"})
		return
	}

	// Загружаем роль для возвращаемого пользователя
	config.DB.Preload("Role").First(&user, user.ID)

	if currentUser, exists := c.Get("currentUser"); exists {
		userInfo := currentUser.(dto.UserInfo)
		err := services.LogAction(userInfo.ID, fmt.Sprintf("Создал пользователя: %s", user.Name))
		if err != nil {
			utils.Log.Errorf("Ошибка при логировании действия: %v", err)
		}
	}

	utils.Log.Errorf("Пользователь %s успешно создан", user.Name)
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
// @Failure 500 {object} dto.ResponseError "Ошибка при получении списка пользователей"
// @Router /users [get]
func GetUsers(c *gin.Context) {
	_, exists := c.Get("currentUser")
	if !exists {
		utils.Log.Warn("Попытка неавторизованного доступа к списку пользователей")
		c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Необходима авторизация"})
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
			utils.Log.Warnf("Роль %s не найдена при фильтрации", roleName)
			c.JSON(http.StatusBadRequest, dto.ResponseError{Message: "Роль не найдена"})
			return
		}
	}

	if err := query.Preload("Role").Find(&users).Error; err != nil {
		utils.Log.Errorf("Ошибка при получении пользователей: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Message: "Ошибка при получении списка пользователей"})
		return
	}

	utils.Log.Info("Получен список пользователей")
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
// @Param input body dto.UpdateUserInput true "Параметры обновления пользователя"
// @Success 200 {object} models.User
// @Failure 400 {object} dto.ResponseError "Ошибка при обновлении пользователя"
// @Failure 404 {object} dto.ResponseError "Пользователь не найден"
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	currentUserRaw, exists := c.Get("currentUser")
	if !exists {
		utils.Log.Warn("Попытка неавторизованного обновления пользователя")
		c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Необходима авторизация"})
		return
	}
	userInfo := currentUserRaw.(dto.UserInfo)
	currentUser := models.User{
		ID:   userInfo.ID,
		Role: userInfo.Role,
	}

	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		utils.Log.Warnf("Пользователь с ID %d не найден для обновления", c.Param("id"))
		c.JSON(http.StatusNotFound, dto.ResponseError{Message: "Пользователь не найден"})
		return
	}

	// Проверка прав доступа:
	// Если текущий пользователь не является администратором или модератором, он не может редактировать чужие данные
	if currentUser.ID != user.ID && currentUser.Role.Name != "admin" && currentUser.Role.Name != "moderator" {
		utils.Log.Warnf("Пользователь %d пытался обновить чужие данные", currentUser.ID)
		c.JSON(http.StatusForbidden, dto.ResponseError{Message: "Недостаточно прав для реадктирования других пользователей"})
		return
	}

	var input dto.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("Ошибка биндинга при обновлении пользователя: %v", err)
		c.JSON(http.StatusBadRequest, dto.ResponseError{Message: err.Error()})
		return
	}

	input.Sanitize()

	// Обновляем данные пользователя
	if err := config.DB.Model(&user).Updates(models.User{
		Name:  input.Name,
		Email: input.Email,
	}).Error; err != nil {
		utils.Log.Errorf("Не удалось обновить пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Message: "Не удалось обновить данные пользователя"})
		return
	}

	err := services.LogAction(currentUser.ID, fmt.Sprintf("Обновил пользователя: %s", user.Name))
	if err != nil {
		utils.Log.Errorf("Ошибка при логировании действия: %v", err)
	}

	// Подгружаем роль перед возвратом
	config.DB.Preload("Role").First(&user, user.ID)
	utils.Log.Infof("Пользователь %s обновлен", user.Email)
	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Удаление пользователя
// @Description Удаление пользователя по его ID.
// @Tags Users
// @Security BearerAuth
// @Param id path int true "ID пользователя"
// @Success 200 {object} dto.ResponseError "Пользователь удален"
// @Failure 404 {object} dto.ResponseError "Пользователь не найден"
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	currentUserRaw, exists := c.Get("currentUser")
	if !exists {
		utils.Log.Warn("Попытка неавторизованного удаления пользователя")
		c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Необходима авторизация"})
		return
	}
	userInfo := currentUserRaw.(dto.UserInfo)
	currentUser := models.User{
		ID:   userInfo.ID,
		Role: userInfo.Role,
	}

	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		utils.Log.Warnf("Пользователь с ID %d не найден для удаления", c.Param("id"))
		c.JSON(http.StatusNotFound, dto.ResponseError{Message: "Пользователь не найден"})
		return
	}

	config.DB.Unscoped().Delete(&user)

	err := services.LogAction(currentUser.ID, fmt.Sprintf("Удалил пользователя: %s", user.Name))
	if err != nil {
		utils.Log.Errorf("Ошибка при логировании действия: %v", err)
	}

	utils.Log.Infof("Пользователь %s удален", user.Email)
	c.JSON(http.StatusOK, dto.ResponseError{Message: "Пользователь удален"})
}

// UpdateUserRole godoc
// @Summary Назначить роль пользователю
// @Description Обновление роли пользователя по ID
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param input body dto.UpdateUserRoleInput true "Новая роль"
// @Success 200 {object} models.User
// @Failure 400 {object} dto.ResponseError "Неверный ввод"
// @Failure 404 {object} dto.ResponseError "Пользователь не найден"
// @Router /users/{id}/role [patch]
func UpdateUserRole(c *gin.Context) {
	currentUserRaw, exists := c.Get("currentUser")
	if !exists {
		utils.Log.Warnf("Попытка неавторизованного обновления роли пользователя")
		c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Необходима авторизация"})
		return
	}
	userInfo := currentUserRaw.(dto.UserInfo)

	var user models.User
	if err := config.DB.Preload("Role").First(&user, c.Param("id")).Error; err != nil {
		utils.Log.Warnf("Пользователь с ID %d не найден для обновления роли", userInfo.ID)
		c.JSON(http.StatusNotFound, dto.ResponseError{Message: "Пользователь не найден"})
		return
	}

	var input dto.UpdateUserRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("Неверный ввод при смене роли пользователя: %v", err)
		c.JSON(http.StatusBadRequest, dto.ResponseError{Message: err.Error()})
		return
	}

	// Найти роль по имени
	var role models.Role
	if err := config.DB.Where("name = ?", input.RoleName).First(&role).Error; err != nil {
		utils.Log.Warnf("Роль с именем %s не найдена", input.RoleName)
		c.JSON(http.StatusBadRequest, dto.ResponseError{Message: "Указанная роль не найдена"})
		return
	}

	oldRole := user.Role.Name
	user.RoleID = role.ID
	user.Role = &role

	if err := config.DB.Save(&user).Error; err != nil {
		utils.Log.Errorf("Не удалось обновить роль пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Message: "Не удалось обновить роль пользователя"})
		return
	}

	err := services.LogAction(userInfo.ID, fmt.Sprintf("Обновил роль пользователя %s: %s -> %s", user.Name, oldRole, user.Role.Name))
	if err != nil {
		utils.Log.Errorf("Ошибка при логировании действия: %v", err)
	}

	// Подгружаем новую роль перед возвратом
	config.DB.Preload("Role").First(&user, user.ID)

	utils.Log.Infof("Роль пользователя %s обновлена на %s", user.Name, user.Role.Name)
	c.JSON(http.StatusOK, user)
}

// BanUser godoc
// @Summary Временная блокировка пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} dto.ResponseMessage
// @Failure 400 {object} dto.ResponseError
// @Failure 403 {object} dto.ResponseError
// @Router /users/{id}/ban [patch]
// @Security BearerAuth
func BanUser(c *gin.Context) {
	currentUserRaw, exists := c.Get("currentUser")
	if !exists {
		utils.Log.Warnf("Попытка неавторизованного блокирования пользователя")
		c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Необходима авторизация"})
		return
	}
	userInfo := currentUserRaw.(dto.UserInfo)
	currentUser := models.User{
		ID:   userInfo.ID,
		Role: userInfo.Role,
	}

	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		utils.Log.Warnf("Пользователь с ID %s не найден для блокировки", currentUser.ID)
		c.JSON(http.StatusNotFound, dto.ResponseError{Message: "Пользователь не найден"})
		return
	}

	// Проверяяем, что пользователь не заблокирован
	if user.IsBanned {
		utils.Log.Warnf("Попытка блокирования уже заблокированного пользователя")
		c.JSON(http.StatusBadRequest, dto.ResponseError{Message: "Пользователь уже заблокирован"})
		return
	}

	user.IsBanned = true
	if err := config.DB.Save(&user).Error; err != nil {
		utils.Log.Errorf("Не удалось заблокировать пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Message: "Не удалось заблокировать пользователя"})
		return
	}

	err := services.LogAction(currentUser.ID, fmt.Sprintf("Заблокировал пользователя: %s", user.Name))
	if err != nil {
		utils.Log.Errorf("Ошибка при логировании действия: %v", err)
	}

	utils.Log.Infof("Пользователь %s заблокирован", user.Email)
	c.JSON(http.StatusOK, dto.ResponseMessage{Message: "Пользователь заблокирован"})
}

// UnbanUser godoc
// @Summary Разблокировка пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} dto.ResponseMessage
// @Failure 400 {object} dto.ResponseError
// @Failure 403 {object} dto.ResponseError
// @Router /users/{id}/unban [patch]
// @Security BearerAuth
func UnbanUser(c *gin.Context) {
	currentUserRaw, exists := c.Get("currentUser")
	if !exists {
		utils.Log.Warnf("Попытка неавторизованного разблокирования пользователя")
		c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Необходима авторизация"})
		return
	}
	userInfo := currentUserRaw.(dto.UserInfo)
	currentUser := models.User{
		ID:   userInfo.ID,
		Role: userInfo.Role,
	}

	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		utils.Log.Warnf("Пользователь с ID %s не найден для разблокировки", currentUser.ID)
		c.JSON(http.StatusNotFound, dto.ResponseError{Message: "Пользователь не найден"})
		return
	}

	// Проверяем, что пользователь не заблокирован
	if !user.IsBanned {
		utils.Log.Warnf("Попытка разблокирования не заблокированного пользователя")
		c.JSON(http.StatusBadRequest, dto.ResponseError{Message: "Пользователь не заблокирован"})
		return
	}

	user.IsBanned = false
	if err := config.DB.Save(&user).Error; err != nil {
		utils.Log.Errorf("Не удалось разблокировать пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Message: "Не удалось разблокировать пользователя"})
		return
	}

	// Логируем действие
	err := services.LogAction(currentUser.ID, fmt.Sprintf("Разблокировал пользователя: %s", user.Name))
	if err != nil {
		utils.Log.Errorf("Ошибка при логировании действия: %v", err)
	}

	utils.Log.Infof("Пользователь %s разблокирован", user.Email)
	c.JSON(http.StatusOK, dto.ResponseMessage{Message: "Пользователь разблокирован"})
}

// GetProfile godoc
// @Summary Получение профиля текущего пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} models.User "Профиль пользователя"
// @Failure 401 {object} dto.ResponseError "Неавторизованный доступ"
// @Failure 404 {object} dto.ResponseError "Пользователь не найден"
// @Failure 500 {object} dto.ResponseError "Ошибка при получении профиля пользователя"
// @Router /users/me [get]
// @Security BearerAuth
func GetProfile(c *gin.Context) {
	currentUserRaw, exists := c.Get("currentUser")
	if !exists {
		utils.Log.Warnf("Попытка получения профиля неавторизованного пользователя")
		c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Пользователь не найден"})
		return
	}

	userInfo := currentUserRaw.(dto.UserInfo)

	var user models.User
	if err := config.DB.Preload("Role").First(&user, userInfo.ID).Error; err != nil {
		utils.Log.Errorf("Не удалось получить профиль пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Message: "Ошибка при получении профиля пользователя"})
		return
	}

	utils.Log.Infof("Получен профиль пользователя %s", user.Email)
	c.JSON(http.StatusOK, user)
}
