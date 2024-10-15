package storage

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

// Информация о ссылке
type URL struct {
	ID            int    `json:"id"`
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
}

// Интерфейс для хранилищ
type Storage interface {
	Load() error
	Save(newURL URL) (int, error)
	CreateShortURL(originalURL string, correlationID string) (string, error)
	GetURL(shortURL string) (string, error)
	Ping() error
	Close() error
}

// Структура, которая содержит один из 3 типов хранилища (Mem, File, DB)
type StorageType struct {
	sType Storage
}

// При отсутствии переменной окружения DATABASE_DSN или флага командной строки -d или при их пустых значениях вернитесь последовательно к:
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

// Реализация методов для интерфейса
func Load(storage *StorageType) error {
	return storage.sType.Load()
}
func CreateShortURL(storage *StorageType, shortURL string, correlationID string) (string, error) {
	return storage.sType.CreateShortURL(shortURL, correlationID)
}
func GetURL(storage *StorageType, shortURL string) (string, error) {
	return storage.sType.GetURL(shortURL)
}
func Ping(storage *StorageType) error {
	return storage.sType.Ping()
}
func Close(storage *StorageType) error {
	return storage.sType.Close()
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
