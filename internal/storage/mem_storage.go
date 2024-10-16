package storage

import (
	"errors"
	"sync"
)

// Массив URL
type MemStorage struct {
	urls map[string]URL
	mu   sync.RWMutex
}

// Создаём новое хранилище
func NewMemStorage() *MemStorage {
	return &MemStorage{}
}

// Создаём мапу с ссылками
func (ms *MemStorage) Load() error {
	ms.urls = map[string]URL{}
	return nil
}

// Сохраняем данные в мапе
func (ms *MemStorage) Save(newURL URL) (int, error) {
	// Добавляем данные в мапу
	ms.urls[newURL.ShortURL] = newURL
	return 0, nil
}

// Создаём короткую ссылку
func (ms *MemStorage) CreateShortURL(originalURL string, correlationID string) (string, error) {
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
			ID:            len(ms.urls) + 1,
		}

		// Добавляем данные в мапу
		ms.Save(newURL)
	} else {
		//Если ссылка уже создана ранее - возвращаем ошибку
		errReturn = errors.New("SHORT_URL_EXIST")
	}

	// Возвращаем короткую ссылку
	return shortURL, errReturn
}

// Ищем ссылку в хранилище
func (ms *MemStorage) GetURL(shortURL string) (string, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	url, ok := ms.urls[shortURL]
	if ok {
		return url.OriginalURL, nil
	} else {
		return "", errors.New("NOT_FOUND")
	}
}

// Пинг БД, не поддерживается
func (ms *MemStorage) Ping() error {
	return errors.New("NOT_SUPPORTED")
}

func (ms *MemStorage) Close() error {
	return nil
}
