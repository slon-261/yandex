package models

// Request описывает запрос пользователя
type Request struct {
	URL string `json:"url"`
}

// Response описывает ответ сервера.
type Response struct {
	Result string `json:"result"`
}

// Request описывает запрос пользователя
type RequestBatch struct {
	URL           string `json:"original_url"`
	CorrelationId string `json:"correlation_id"`
}

// Response описывает ответ сервера.
type ResponseBatch struct {
	ShortURL      string `json:"short_url"`
	CorrelationId string `json:"correlation_id"`
}
