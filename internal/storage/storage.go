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
	ID          int    `json:"id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Storage interface {
	Load() error
	Save(data []byte) error
	CreateShortURL(originalURL string) string
}

// Массив URL + указатель на файл
type FileStorage struct {
	filename string
	file     *os.File
	scanner  *bufio.Scanner
	urls     map[string]URL
	mu       sync.Mutex
}

func NewFileStorage(filename string) *FileStorage {
	return &FileStorage{filename: filename}
}

func (fs *FileStorage) Load() error {
	var err error
	// Пытаемся создать директорию
	os.MkdirAll(filepath.Dir(fs.filename), 0666)
	// Создаём файл
	fs.file, err = os.OpenFile(fs.filename, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return err
	}
	// создаём новый scanner
	fs.scanner = bufio.NewScanner(fs.file)
	fs.urls = map[string]URL{}
	// перебираем все строки
	for fs.scanner.Scan() {
		// читаем данные из scanner
		data := fs.scanner.Bytes()

		url := URL{}
		err := json.Unmarshal(data, &url)

		fs.urls[url.ShortURL] = url

		if err != nil {
			return err
		}
	}
	return nil
}

func (fs *FileStorage) Save(data []byte) (int, error) {
	return fs.file.Write(data)
}

// Создаём короткую ссылку
func (fs *FileStorage) CreateShortURL(originalURL string) string {
	// Получаем хэш
	shortURL := encryption(originalURL)
	// Ищем ссылку в хранилище. Если не нашли - добавляем
	_, err := fs.GetURL(shortURL)
	if err != nil {
		newURL := URL{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
			ID:          len(fs.urls) + 1,
		}
		// Добавляем данные в мапу
		fs.urls[shortURL] = newURL

		data, _ := json.Marshal(&newURL)
		// добавляем перенос строки
		data = append(data, '\n')

		// Добавляем данные в файл
		fs.file.Write(data)
	}

	// Возвращаем короткую ссылку
	return shortURL
}

// Ищем ссылку в хранилище
func (fs *FileStorage) GetURL(shortURL string) (string, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	url, ok := fs.urls[shortURL]
	if ok {
		return url.OriginalURL, nil
	} else {
		return "", errors.New("NOT_FOUND")
	}
}

func (fs *FileStorage) Close() error {
	return fs.file.Close()
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
