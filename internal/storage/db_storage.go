package storage

import (
	"database/sql"
	"errors"
	"log"
	"sync"
)

// Массив URL + указатель на файл
type DBStorage struct {
	DSN  string
	db   *sql.DB
	urls map[string]URL
	mu   sync.Mutex
}

// Создаём новое хранилище, подключаемся к БД и создаём таблицу
func NewDBStorage(DSN string) *DBStorage {
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS urls (id integer, correlation_id varchar, short_url varchar, original_url varchar);")
	if err != nil {
		log.Print(err)
	}
	return &DBStorage{DSN: DSN, db: db}
}

// Создаём мапу с ссылками и подгружаем туда данные из БД
func (ds *DBStorage) Load() error {
	ds.urls = map[string]URL{}

	rows, err := ds.db.Query("SELECT id, correlation_id, short_url, original_url from urls")
	if err != nil {
		return err
	}
	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var url URL
		err = rows.Scan(&url.ID, &url.CorrelationID, &url.ShortURL, &url.OriginalURL)
		if err != nil {
			return err
		}
		ds.urls[url.ShortURL] = url
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

// Сохраняем данные в мапе и БД
func (ds *DBStorage) Save(newURL URL) (int, error) {
	// Добавляем данные в мапу
	ds.urls[newURL.ShortURL] = newURL

	_, err := ds.db.Exec(`
        INSERT INTO urls
        (id, correlation_id, short_url, original_url)
        VALUES
        ($1, $2, $3, $4);
		`, newURL.ID, newURL.CorrelationID, newURL.ShortURL, newURL.OriginalURL)
	if err != nil {
		log.Print(err)
	}
	return 1, err
}

// Создаём короткую ссылку
func (ds *DBStorage) CreateShortURL(originalURL string, correlationID string) (string, error) {
	// Получаем хэш
	shortURL := encryption(originalURL)
	//Возвращаемая ошибка
	var errReturn error
	// Ищем ссылку в хранилище. Если не нашли - добавляем
	_, err := ds.GetURL(shortURL)
	if err != nil {
		newURL := URL{
			ShortURL:      shortURL,
			OriginalURL:   originalURL,
			CorrelationID: correlationID,
			ID:            len(ds.urls) + 1,
		}
		// Добавляем данные в БД
		ds.Save(newURL)

		errReturn = nil
	} else {
		errReturn = errors.New("SHORT_URL_EXIST")
	}

	// Возвращаем короткую ссылку
	return shortURL, errReturn
}

// Ищем ссылку в хранилище
func (ds *DBStorage) GetURL(shortURL string) (string, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	url, ok := ds.urls[shortURL]
	if ok {
		return url.OriginalURL, nil
	} else {
		return "", errors.New("NOT_FOUND")
	}
}

// Пинг БД
func (ds *DBStorage) Ping() error {
	return ds.db.Ping()
}

func (ds *DBStorage) Close() error {
	return ds.db.Close()
}
