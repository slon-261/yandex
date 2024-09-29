package storage

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"sync"
)

type autoInc struct {
	sync.Mutex
	id int
}

func (a *autoInc) ID() (id int) {
	a.Lock()
	defer a.Unlock()

	id = a.id
	a.id++
	return
}

var ai autoInc // Глобальная переменная

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

func (storage *Storage) CreateShortURL(originalURL string) string {

	shortURL := encryption(originalURL)

	_, err := storage.GetURL(shortURL)
	if err != nil {
		newURL := URL{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
			UUID:        len(storage.urls) + 1,
		}
		storage.urls[shortURL] = newURL

		data, _ := json.Marshal(&newURL)
		// добавляем перенос строки
		data = append(data, '\n')

		storage.file.Write(data)
	}

	return shortURL
}

func (storage *Storage) GetURL(shortURL string) (string, error) {
	url, ok := storage.urls[shortURL]
	if ok {
		return url.OriginalURL, nil
	} else {
		return "", errors.New("Not found")
	}
}

func (storage *Storage) Load(filename string) error {
	var err error
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
