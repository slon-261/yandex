package models

// Запрос на создание ссылки
type Request struct {
	URL string `json:"url"`
}

// Ответ при создании ссылки
type Response struct {
	Result string `json:"result"`
}

// Запрос на массовое создание ссылок
type RequestBatch struct {
	URL           string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
}

// Ответ при массовом создании ссылок
type ResponseBatch struct {
	ShortURL      string `json:"short_url"`
	CorrelationID string `json:"correlation_id"`
}

// Ответ при запросе ссылок по пользователю
type ResponseUserUrls struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
