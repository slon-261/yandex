package storage

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
)

// ErrNotFound Ошибка "Не найдено"
var ErrNotFound = errors.New("NOT_FOUND")

// ErrShortURLExist Ошибка "Короткая ссылка существует"
var ErrShortURLExist = errors.New("SHORT_URL_EXIST")

// ErrNotSupported Ошибка "Не поддерживается"
var ErrNotSupported = errors.New("NOT_SUPPORTED")

// ErrShortURLDeleted Ошибка "Удалено"
var ErrShortURLDeleted = errors.New("DELETED")

// URL Информация о ссылке
type URL struct {
	ID            int    `json:"id"`
	CorrelationID string `json:"correlation_id"`
	UserID        string `json:"user_id"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
	DeletedFlag   bool   `json:"deleted_flag"`
}

// Storage Интерфейс для хранилищ
type Storage interface {
	Load() error
	Save(newURL URL) (int, error)
	CreateShortURL(originalURL string, correlationID string, userID string) (string, error)
	GetURL(shortURL string) (string, error)
	GetUserURLs(userID string) ([]URL, error)
	DeleteUserURLs(userID string, urls []string) error
	Ping() error
	Close() error
}

// StorageType Структура, которая содержит один из 3 типов хранилища (Mem, File, DB)
type StorageType struct {
	sType Storage
}

// NewStorage При отсутствии переменной окружения DATABASE_DSN или флага командной строки -d или при их пустых значениях вернитесь последовательно к:
// хранению сокращённых URL в файле при наличии соответствующей переменной окружения или флага командной строки;
// хранению сокращённых URL в памяти.
func NewStorage(flagDataBaseDSN string, flagFilePath string) *StorageType {
	//Храним в БД
	if flagDataBaseDSN != "" {
		return &StorageType{NewDBStorage(flagDataBaseDSN)}
		//Храним в файле
	} else if flagFilePath != "" {
		return &StorageType{NewFileStorage(flagFilePath)}
		//Храним в памяти
	} else {
		return &StorageType{NewMemStorage()}
	}
}

// Load Создаём мапу с ссылками и подгружаем туда данные
func Load(storage *StorageType) error {
	return storage.sType.Load()
}

// CreateShortURL Создаём короткую ссылку
func CreateShortURL(storage *StorageType, shortURL string, correlationID string, userID string) (string, error) {
	return storage.sType.CreateShortURL(shortURL, correlationID, userID)
}

// GetURL Ищем ссылку в хранилище
func GetURL(storage *StorageType, shortURL string) (string, error) {
	return storage.sType.GetURL(shortURL)
}

// GetUserURLs Получаем все ссылки текущего пользователя
func GetUserURLs(storage *StorageType, userID string) ([]URL, error) {
	return storage.sType.GetUserURLs(userID)
}

// DeleteUserURLs Удаление ссылок
func DeleteUserURLs(storage *StorageType, userID string, urls []string) error {
	return storage.sType.DeleteUserURLs(userID, urls)
}

// Ping Пинг БД
func Ping(storage *StorageType) error {
	return storage.sType.Ping()
}

// Close Закрытие соединения
func Close(storage *StorageType) error {
	return storage.sType.Close()
}

// Encryption Шифрование строки
func Encryption(str string) string {
	// Генерируем короткую ссылку
	h := sha256.New()
	h.Write([]byte(str))
	hashString := base64.StdEncoding.EncodeToString(h.Sum(nil))
	// Удаляем / из короткой ссылки
	hashString = strings.ReplaceAll(hashString, "/", "")
	return hashString[:10]
}
