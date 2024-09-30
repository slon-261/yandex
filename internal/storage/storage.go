package storage

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Информация о ссылке
type URL struct {
	Id          int    `json:"id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// Массив URL + указатель на файл
type Storage struct {
	file    *os.File
	scanner *bufio.Scanner
	urls    map[string]URL
	mu      sync.Mutex
}

// Создаём короткую ссылку
func (storage *Storage) CreateShortURL(originalURL string) string {
	// Получаем хэш
	shortURL := encryption(originalURL)
	// Ищем ссылку в хранилище. Если не нашли - добавляем
	_, err := storage.GetURL(shortURL)
	if err != nil {
		newURL := URL{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
			Id:          len(storage.urls) + 1,
		}
		// Добавляем данные в мапу
		storage.urls[shortURL] = newURL

		data, _ := json.Marshal(&newURL)
		// добавляем перенос строки
		data = append(data, '\n')

		// Добавляем данные в файл
		storage.file.Write(data)
	}

	// Возвращаем короткую ссылку
	return shortURL
}

// Ищем ссылку в хранилище
func (storage *Storage) GetURL(shortURL string) (string, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	url, ok := storage.urls[shortURL]
	if ok {
		return url.OriginalURL, nil
	} else {
		return "", errors.New("NOT_FOUND")
	}
}

// Загружаем из файла все ранее сгенерированные ссылки
func (storage *Storage) Load(filename string) error {
	var err error
	// Пытаемся создать директорию
	os.MkdirAll(filepath.Dir(filename), 0666)
	// Создаём файл
	storage.file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return err
	}
	// создаём новый scanner
	storage.scanner = bufio.NewScanner(storage.file)
	storage.urls = map[string]URL{}
	// перебираем все строки
	for storage.scanner.Scan() {
		// читаем данные из scanner
		data := storage.scanner.Bytes()

		url := URL{}
		err := json.Unmarshal(data, &url)

		storage.urls[url.ShortURL] = url

		if err != nil {
			return err
		}
	}
	return nil
}

func (storage *Storage) Close() error {
	return storage.file.Close()
}

func encryption(str string) string {
	// Генерируем короткую ссылку
	h := sha256.New()
	h.Write([]byte(str))
	hashString := base64.StdEncoding.EncodeToString(h.Sum(nil))
	// Удаляем / из короткой ссылки
	hashString = strings.ReplaceAll(hashString, "/", "")
	return hashString[:10]
}
