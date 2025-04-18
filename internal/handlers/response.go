package handlers

// ResponseError - структура для ответа с ошибкой
type ResponseError struct {
	Message string `json:"message"`
}

// ResponseMessage - структура для ответа с данными
type ResponseMessage struct {
	Message string `json:"message"`
}
