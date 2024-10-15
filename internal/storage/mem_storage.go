package storage

import (
	"errors"
	"sync"
)

// Массив URL
type MemStorage struct {
	urls map[string]URL
	mu   sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{}
}

func (ms *MemStorage) Load() error {
	ms.urls = map[string]URL{}
	return nil
}

func (ms *MemStorage) Save(shortURL string, newURL URL) (int, error) {
	// Добавляем данные в мапу
	ms.urls[shortURL] = newURL
	return 0, nil
}

// Создаём короткую ссылку
func (ms *MemStorage) CreateShortURL(originalURL string) string {
	// Получаем хэш
	shortURL := encryption(originalURL)
	// Ищем ссылку в хранилище. Если не нашли - добавляем
	_, err := ms.GetURL(shortURL)
	if err != nil {
		newURL := URL{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
			ID:          len(ms.urls) + 1,
		}

		// Добавляем данные в мапу
		ms.Save(shortURL, newURL)
	}

	// Возвращаем короткую ссылку
	return shortURL
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

func (ms *MemStorage) Ping() error {
	return errors.New("NOT_SUPPORTED")
}

func (ms *MemStorage) Close() error {
	return nil
}
