package storage

import (
	"database/sql"
	"log"
	"sync"
)

// Массив URL + указатель на файл
type DBStorage struct {
	DSN string
	db  *sql.DB
	mu  sync.RWMutex
}

// Создаём новое хранилище, подключаемся к БД и создаём таблицу
func NewDBStorage(DSN string) *DBStorage {
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS urls (id serial PRIMARY KEY, correlation_id varchar, user_id varchar, short_url varchar UNIQUE, original_url varchar);")
	if err != nil {
		panic(err)
	}
	return &DBStorage{DSN: DSN, db: db}
}

// Для БД неактуально, не загружаем данные в мапу
func (ds *DBStorage) Load() error {
	return nil
}

// Сохраняем данные в мапе и БД
func (ds *DBStorage) Save(newURL URL) (int, error) {
	_, err := ds.db.Exec(`
        INSERT INTO urls
        (correlation_id, user_id, short_url, original_url)
        VALUES
        ($1, $2, $3, $4);
		`, newURL.CorrelationID, newURL.UserID, newURL.ShortURL, newURL.OriginalURL)
	if err != nil {
		log.Print(err)
	}
	return 1, err
}

// Создаём короткую ссылку
func (ds *DBStorage) CreateShortURL(originalURL string, correlationID string, userID string) (string, error) {
	// Получаем хэш
	shortURL := Encryption(originalURL)
	//Возвращаемая ошибка
	var errReturn error
	// Ищем ссылку в хранилище. Если не нашли - добавляем
	_, err := ds.GetURL(shortURL)
	if err != nil {
		newURL := URL{
			ShortURL:      shortURL,
			OriginalURL:   originalURL,
			UserID:        userID,
			CorrelationID: correlationID,
		}
		// Добавляем данные в БД
		ds.Save(newURL)

		errReturn = nil
	} else {
		errReturn = ErrShortURLExist
	}

	// Возвращаем короткую ссылку
	return shortURL, errReturn
}

// Ищем ссылку в хранилище
func (ds *DBStorage) GetURL(shortURL string) (string, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	row := ds.db.QueryRow("SELECT id, correlation_id, short_url, original_url from urls WHERE short_url = $1 LIMIT 1", shortURL)
	var url URL
	err := row.Scan(&url.ID, &url.CorrelationID, &url.ShortURL, &url.OriginalURL)

	if err != nil {
		return "", err
	} else {
		return url.OriginalURL, nil
	}
}

// Получаем все ссылки текущего пользователя
func (ds *DBStorage) GetUserURLs(userID string) ([]URL, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	var resp []URL

	rows, err := ds.db.Query("SELECT id, correlation_id, short_url, original_url from urls WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var url URL
		err = rows.Scan(&url.ID, &url.CorrelationID, &url.ShortURL, &url.OriginalURL)
		if err != nil {
			return nil, err
		}
		resp = append(resp, url)
	}
	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(resp) > 0 {
		return resp, nil
	} else {
		return nil, ErrNotFound
	}
}

// Пинг БД
func (ds *DBStorage) Ping() error {
	return ds.db.Ping()
}

func (ds *DBStorage) Close() error {
	return ds.db.Close()
}
