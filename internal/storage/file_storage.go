package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// FileStorage Массив URL + указатель на файл
type FileStorage struct {
	filename string
	file     *os.File
	scanner  *bufio.Scanner
	urls     map[string]URL
	mu       sync.RWMutex
}

// NewFileStorage Создаём новое хранилище, открываем файл
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

// Load Создаём мапу с ссылками и подгружаем туда данные из файла
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

// Save Сохраняем данные в мапе и в файле
func (fs *FileStorage) Save(newURL URL) (int, error) {
	// Добавляем данные в мапу
	fs.urls[newURL.ShortURL] = newURL

	data, err := json.Marshal(&newURL)
	if err != nil {
		return 0, err
	}
	// добавляем перенос строки
	data = append(data, '\n')

	return fs.file.Write(data)
}

// CreateShortURL Создаём короткую ссылку
func (fs *FileStorage) CreateShortURL(originalURL string, correlationID string, userID string) (string, error) {
	// Получаем хэш
	shortURL := Encryption(originalURL)
	//Возвращаемая ошибка
	var errReturn error
	// Ищем ссылку в хранилище. Если не нашли - добавляем
	_, err := fs.GetURL(shortURL)
	if err != nil {
		newURL := URL{
			ShortURL:      shortURL,
			OriginalURL:   originalURL,
			CorrelationID: correlationID,
			UserID:        userID,
			ID:            len(fs.urls) + 1,
		}

		// Добавляем данные в файл
		fs.Save(newURL)
		errReturn = nil
	} else {
		//Если ссылка уже создана ранее - возвращаем ошибку
		errReturn = ErrShortURLExist
	}

	// Возвращаем короткую ссылку
	return shortURL, errReturn
}

// GetURL Ищем ссылку в хранилище
func (fs *FileStorage) GetURL(shortURL string) (string, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	url, ok := fs.urls[shortURL]
	if ok {
		return url.OriginalURL, nil
	} else {
		return "", ErrNotFound
	}
}

// GetUserURLs Получаем все ссылки текущего пользователя
func (fs *FileStorage) GetUserURLs(userID string) ([]URL, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	var resp []URL
	// Перебираем всю мапу, берем только нужные объекты
	for _, element := range fs.urls {
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
func (fs *FileStorage) DeleteUserURLs(userID string, urls []string) error {
	return ErrNotSupported
}

// Ping Пинг БД, не поддерживается
func (fs *FileStorage) Ping() error {
	return ErrNotSupported
}

// Close Закртыие файла
func (fs *FileStorage) Close() error {
	return fs.file.Close()
}
