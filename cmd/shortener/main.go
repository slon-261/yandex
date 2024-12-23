package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"slon-261/yandex/internal/auth"
	d "slon-261/yandex/internal/decompress"
	l "slon-261/yandex/internal/logger"
	"slon-261/yandex/internal/models"
	s "slon-261/yandex/internal/storage"
	"strings"
)

// Хранилище ссылок
var storage *s.StorageType

// postPage хэндлер сокращения ссылок
func postPage(w http.ResponseWriter, r *http.Request) {

	// Получаем ссылку из body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Body error", http.StatusBadRequest)
		return
	}
	originalURL := strings.TrimSpace(string(body))

	// Сохраняем короткую ссылку
	shortURL, errCreate := s.CreateShortURL(storage, originalURL, "", auth.GetCurrentUserID())
	//Если при создании ссылки была ошибка - возвращаем код 409, но в теле указыаем ссылку
	var status int
	if errCreate != nil {
		status = http.StatusConflict
	} else {
		status = http.StatusCreated
	}

	response := flagBaseURL + "/" + shortURL

	// Выводим новую ссылку на экран
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// postJSONPage хэндлер сокращения ссылок в формате JSON
func postJSONPage(w http.ResponseWriter, r *http.Request) {

	// Получаем ссылку из body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Body error", http.StatusBadRequest)
		return
	}

	var req models.Request
	if err = json.Unmarshal([]byte(body), &req); err != nil {
		http.Error(w, "JSON error", http.StatusBadRequest)
		return
	}

	auth.GetUserID(r)

	// Сохраняем короткую ссылку
	shortURL, errCreate := s.CreateShortURL(storage, req.URL, "", auth.GetCurrentUserID())
	//Если при создании ссылки была ошибка - возвращаем код 409, но в теле указыаем ссылку
	var status int
	if errCreate != nil {
		status = http.StatusConflict
	} else {
		status = http.StatusCreated
	}

	var resp models.Response
	resp.Result = flagBaseURL + "/" + shortURL
	responseJSON, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		http.Error(w, "JSON error", http.StatusInternalServerError)
		return
	}

	// Выводим новую ссылку на экран
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	w.Write(responseJSON)
}

// postJSONPage хэндлер пакетного сокращения ссылок
func postBatchPage(w http.ResponseWriter, r *http.Request) {

	// Получаем ссылку из body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Body error", http.StatusBadRequest)
		return
	}

	var req []models.RequestBatch
	if err = json.Unmarshal([]byte(body), &req); err != nil {
		log.Print(err)
		http.Error(w, "JSON error", http.StatusBadRequest)
		return
	}

	var resp []models.ResponseBatch
	var respCurr models.ResponseBatch
	for _, element := range req {
		// Сохраняем короткую ссылку
		shortURL, _ := s.CreateShortURL(storage, element.URL, element.CorrelationID, auth.GetCurrentUserID())
		respCurr.ShortURL = flagBaseURL + "/" + shortURL
		respCurr.CorrelationID = element.CorrelationID
		resp = append(resp, respCurr)
	}

	responseJSON, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		http.Error(w, "JSON error", http.StatusInternalServerError)
		return
	}

	// Выводим новую ссылку на экран
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}

// getPage хэндлер получения полной ссылки
func getPage(w http.ResponseWriter, r *http.Request) {

	// Получаем короткую ссылку
	//shortURL := chi.URLParam(r, "url")
	shortURL := strings.Trim(string(r.RequestURI), " /")

	// Ищем ссылку в хранилище
	url, err := s.GetURL(storage, shortURL)
	if err != nil {
		if err == s.ErrShortURLDeleted {
			http.Error(w, "Deleted", http.StatusGone)
			return
		}
		http.Error(w, "Not found", http.StatusBadRequest)
	} else {
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

// userURLsPage хэндлер получения всех ссылок для пользователя
func userURLsPage(w http.ResponseWriter, r *http.Request) {
	//Если не смогли получить из куков - ошибка
	if auth.GetUserID(r) == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем все ссылки по указанному пользователю
	urls, err := s.GetUserURLs(storage, auth.GetCurrentUserID())
	if err != nil {
		http.Error(w, "No content", http.StatusNoContent)
		return
	} else {
		// Преобразуем в нужный вид, в том числе добавляем flagBaseURL
		var resp []models.ResponseUserUrls
		var respCurr models.ResponseUserUrls
		for _, element := range urls {
			// Сохраняем короткую ссылку
			respCurr.ShortURL = flagBaseURL + "/" + element.ShortURL
			respCurr.OriginalURL = element.OriginalURL
			resp = append(resp, respCurr)
		}

		responseJSON, err := json.MarshalIndent(resp, "", "   ")
		if err != nil {
			http.Error(w, "JSON error", http.StatusInternalServerError)
			return
		}

		// Выводим новую ссылку на экран
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}
}

// deleteUserURLsPage хэндлер удаления ссылок
func deleteUserURLsPage(w http.ResponseWriter, r *http.Request) {
	// Если не смогли получить из куков - ошибка
	if auth.GetUserID(r) == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем массив ссылок из body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Body error", http.StatusBadRequest)
		return
	}

	var req []string
	if err = json.Unmarshal([]byte(body), &req); err != nil {
		log.Print(err)
		http.Error(w, "JSON error", http.StatusBadRequest)
		return
	}

	// Удаляем переданные ссылки, при условии что они принадлежат указанному пользователю
	err = s.DeleteUserURLs(storage, auth.GetCurrentUserID(), req)
	if err != nil {
		log.Print(err)
		http.Error(w, "Delete error", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

// pingPage хэндлер проверки соеднения с БД
func pingPage(w http.ResponseWriter, r *http.Request) {
	err := s.Ping(storage)

	if err != nil {
		http.Error(w, "Connect failed", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// createRouter создание роутера
func createRouter() *chi.Mux {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := chi.NewRouter()
	r.Use(l.RequestLogger(logger)) // Логгирование
	r.Use(middleware.Compress(5))  // Сжатие ответа
	r.Use(d.Decompress)            // Распаковка сжатого запроса
	r.Use(auth.Authenticator())
	r.Post("/", postPage)
	r.Post("/api/shorten", postJSONPage)
	r.Post("/api/shorten/batch", postBatchPage)
	r.Get("/ping", pingPage)
	r.Get("/api/user/urls", userURLsPage)
	r.Delete("/api/user/urls", deleteUserURLsPage)
	r.Get("/{url}", getPage)
	return r
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Print("Panic occurred: ", r)
		}
	}()

	// обрабатываем аргументы командной строки
	parseFlags()

	storage = s.NewStorage(flagDataBaseDSN, flagFilePath)

	// Загружаем из файла\БД все ранее сгенерированные ссылки
	s.Load(storage)
	defer s.Close(storage)

	r := createRouter()

	log.Print("Running server on ", flagRunAddr)
	log.Print("File storage is ", flagFilePath)
	log.Print("DB connected at ", flagDataBaseDSN)

	// создаём файл журнала профилирования памяти
	fmemPath := "./profiles/base.pprof"
	os.MkdirAll(filepath.Dir(fmemPath), 0666)
	fmem, err := os.Create(fmemPath)

	if err != nil {
		panic(err)
	}
	defer fmem.Close()
	runtime.GC() // получаем статистику по использованию памяти
	if err := pprof.WriteHeapProfile(fmem); err != nil {
		panic(err)
	}

	// r передаётся как http.Handler
	http.ListenAndServe(flagRunAddr, r)

}
