package handlers

import (
	"net/http"
	"time"
	"userManagement/internal/config"
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
// @Param user body handlers.RegisterInput true "Регистрационные данные"
// @Success 201 {object} ResponseMessage "Регистрация прошла успешно"
// @Failure 400 {object} ResponseError "Ошибка при валидации данных"
// @Failure 500 {object} ResponseError "Ошибка при хешировании пароля или сохранении данных"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Ошибка при хешировании пароля"})
		return
	}

	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		RoleID:       3,
	}

	if err1 := config.DB.Create(&user).Error; err1 != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: "Не удалось зарегистрироваться"})
		return
	}

	c.JSON(http.StatusCreated, ResponseMessage{Message: "Регистрация прошла успешно"})
}

// Login godoc
// @Summary Вход пользователя в систему
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body handlers.LoginInput true "Данные для входа"
// @Success 200 {object} AuthResponse "JWT токен"
// @Failure 400 {object} ResponseError "Ошибка при валидации данных"
// @Failure 401 {object} ResponseError "Неверный email или пароль"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, ResponseError{Message: "Неверный email или пароль"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, ResponseError{Message: "Неверный email или пароль"})
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
		c.JSON(http.StatusInternalServerError, ResponseError{Message: "Ошибка создания токена"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{Token: tokenString})
}
