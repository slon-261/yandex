package storage

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Информация о ссылке
type URL struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// Массив URL + указатель на файл
type Storage struct {
	file    *os.File
	scanner *bufio.Scanner
	urls    map[string]URL
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
			UUID:        len(storage.urls) + 1,
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
	storage.file, err = os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)

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
		log.Print(url)
		err := json.Unmarshal(data, &url)

		log.Print(err)
		log.Print(url)
		storage.urls[url.ShortURL] = url

		if err != nil {
			return err
		}
	}
	return nil
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
