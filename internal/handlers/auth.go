package handlers

import (
	"net/http"
	"time"
	"userManagement/internal/config"
	"userManagement/internal/dto"
	"userManagement/internal/models"
	"userManagement/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary Регистрация нового пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterInput true "Регистрационные данные"
// @Success 201 {object} dto.ResponseMessage "Регистрация прошла успешно"
// @Failure 400 {object} dto.ResponseError "Ошибка при валидации данных"
// @Failure 500 {object} dto.ResponseError "Ошибка при хешировании пароля или сохранении данных"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var input dto.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("Ошибка валидации при регистрации: %v", err)
		c.JSON(http.StatusBadRequest, dto.ResponseError{Message: err.Error()})
		return
	}

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
		RoleID:       3, // обычный пользователь
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.Log.Errorf("Ошибка создания пользователя: %v", err)
		c.JSON(http.StatusBadRequest, dto.ResponseError{Message: "Не удалось зарегистрироваться"})
		return
	}

	utils.Log.Info("Пользователь %s (%s) зарегистрирован", user.Name, user.Email)
	c.JSON(http.StatusCreated, dto.ResponseMessage{Message: "Регистрация прошла успешно"})
}

// Login godoc
// @Summary Вход пользователя в систему
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body dto.LoginInput true "Данные для входа"
// @Success 200 {object} dto.AuthResponse "JWT токен"
// @Failure 400 {object} dto.ResponseError "Ошибка при валидации данных"
// @Failure 401 {object} dto.ResponseError "Неверный email или пароль"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Log.Warnf("Ошибка валидации при входе: %v", err)
		c.JSON(http.StatusBadRequest, dto.ResponseError{Message: err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		utils.Log.Warnf("Попытка входа с неверным email: %s", input.Email)
		c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Неверный email или пароль"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		utils.Log.Warnf("Неверный пароль для пользователя: %s", user.Email)
		c.JSON(http.StatusUnauthorized, dto.ResponseError{Message: "Неверный email или пароль"})
		return
	}

	// Создаем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"role":   user.Role,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		utils.Log.Errorf("Ошибка при создании токена для %s: %v", user.Email, err)
		c.JSON(http.StatusInternalServerError, dto.ResponseError{Message: "Ошибка создания токена"})
		return
	}

	utils.Log.Infof("Пользователь %s успешно вошел в систему", user.Email)
	c.JSON(http.StatusOK, dto.AuthResponse{Token: tokenString})
}
