package dto

// ResponseError - структура для ответа с ошибкой
type ResponseError struct {
	Message string `json:"message"`
}

// ResponseMessage - структура для ответа с данными
type ResponseMessage struct {
	Message string `json:"message"`
}

// AuthResponse структура для ответа при успешной аутентификации
type AuthResponse struct {
	Token string `json:"token"`
}
