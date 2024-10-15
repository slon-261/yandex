package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// Массив URL + указатель на файл
type FileStorage struct {
	filename string
	file     *os.File
	scanner  *bufio.Scanner
	urls     map[string]URL
	mu       sync.Mutex
}

func NewFileStorage(filename string) *FileStorage {
	// Пытаемся создать директорию
	os.MkdirAll(filepath.Dir(filename), 0666)
	// Создаём файл
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	return &FileStorage{filename: filename, file: file}
}

func (fs *FileStorage) Load() error {
	fs.urls = map[string]URL{}
	// создаём новый scanner
	fs.scanner = bufio.NewScanner(fs.file)
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

func (fs *FileStorage) Save(shortURL string, newURL URL) (int, error) {
	// Добавляем данные в мапу
	fs.urls[shortURL] = newURL

	data, _ := json.Marshal(&newURL)
	// добавляем перенос строки
	data = append(data, '\n')

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

		// Добавляем данные в файл
		fs.Save(shortURL, newURL)
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

func (fs *FileStorage) Ping() error {
	return errors.New("NOT_SUPPORTED")
}

func (fs *FileStorage) Close() error {
	return fs.file.Close()
}
