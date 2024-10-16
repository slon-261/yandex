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
	mu       sync.RWMutex
}

// Создаём новое хранилище, открываем файл
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

// Создаём мапу с ссылками и подгружаем туда данные из файла
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

// Сохраняем данные в мапе и в файле
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

// Создаём короткую ссылку
func (fs *FileStorage) CreateShortURL(originalURL string, correlationID string) (string, error) {
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
			ID:            len(fs.urls) + 1,
		}

		// Добавляем данные в файл
		fs.Save(newURL)
		errReturn = nil
	} else {
		//Если ссылка уже создана ранее - возвращаем ошибку
		errReturn = errors.New("SHORT_URL_EXIST")
	}

	// Возвращаем короткую ссылку
	return shortURL, errReturn
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

// Пинг БД, не поддерживается
func (fs *FileStorage) Ping() error {
	return errors.New("NOT_SUPPORTED")
}

func (fs *FileStorage) Close() error {
	return fs.file.Close()
}
