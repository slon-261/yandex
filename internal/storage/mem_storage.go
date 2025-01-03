package storage

import (
	"sync"
)

// MemStorage Массив URL
type MemStorage struct {
	urls map[string]URL
	mu   sync.RWMutex
}

// NewMemStorage Создаём новое хранилище
func NewMemStorage() *MemStorage {
	return &MemStorage{}
}

// Load Создаём мапу с ссылками
func (ms *MemStorage) Load() error {
	ms.urls = map[string]URL{}
	return nil
}

// Save Сохраняем данные в мапе
func (ms *MemStorage) Save(newURL URL) (int, error) {
	// Добавляем данные в мапу
	ms.urls[newURL.ShortURL] = newURL
	return 0, nil
}

// CreateShortURL Создаём короткую ссылку
func (ms *MemStorage) CreateShortURL(originalURL string, correlationID string, userID string) (string, error) {
	// Получаем хэш
	shortURL := Encryption(originalURL)
	//Возвращаемая ошибка
	var errReturn error
	// Ищем ссылку в хранилище. Если не нашли - добавляем
	_, err := ms.GetURL(shortURL)
	if err != nil {
		newURL := URL{
			ShortURL:      shortURL,
			OriginalURL:   originalURL,
			CorrelationID: correlationID,
			UserID:        userID,
			ID:            len(ms.urls) + 1,
		}

		// Добавляем данные в мапу
		ms.Save(newURL)
	} else {
		//Если ссылка уже создана ранее - возвращаем ошибку
		errReturn = ErrShortURLExist
	}

	// Возвращаем короткую ссылку
	return shortURL, errReturn
}

// GetURL Ищем ссылку в хранилище
func (ms *MemStorage) GetURL(shortURL string) (string, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	url, ok := ms.urls[shortURL]
	if ok {
		return url.OriginalURL, nil
	} else {
		return "", ErrNotFound
	}
}

// GetUserURLs Получаем все ссылки текущего пользователя
func (ms *MemStorage) GetUserURLs(userID string) ([]URL, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var resp []URL
	// Перебираем всю мапу, берем только нужные объекты
	for _, element := range ms.urls {
		if element.UserID == userID {
			resp = append(resp, element)
		}
	}
	if len(resp) > 0 {
		return resp, nil
	} else {
		return nil, ErrNotFound
	}
}

// DeleteUserURLs Удаление ссылок, не поддерживается
func (ms *MemStorage) DeleteUserURLs(userID string, urls []string) error {
	return ErrNotSupported
}

// Ping Пинг БД, не поддерживается
func (ms *MemStorage) Ping() error {
	return ErrNotSupported
}

// Close Закрытие соединения, не поддерживается
func (ms *MemStorage) Close() error {
	return nil
}
