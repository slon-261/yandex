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
	d "slon-261/yandex/internal/decompress"
	l "slon-261/yandex/internal/logger"
	"slon-261/yandex/internal/models"
	s "slon-261/yandex/internal/storage"
	"strings"
)

// Хранилище ссылок
var storage *s.StorageType

func postPage(w http.ResponseWriter, r *http.Request) {

	// Получаем ссылку из body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Body error", http.StatusBadRequest)
		return
	}
	originalURL := strings.TrimSpace(string(body))

	// Сохраняем короткую ссылку
	shortURL, errCreate := s.CreateShortURL(storage, originalURL, "")
	response := flagBaseURL + "/" + shortURL

	//Если при создании ссылки была ошибка - возвращаем код 409, но в теле указыаем ссылку
	if errCreate != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(response))
		return
	}

	// Выводим новую ссылку на экран
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

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

	// Сохраняем короткую ссылку
	shortURL, errCreate := s.CreateShortURL(storage, req.URL, "")
	var resp models.Response
	resp.Result = flagBaseURL + "/" + shortURL
	responseJSON, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		http.Error(w, "JSON error", http.StatusInternalServerError)
		return
	}

	//Если при создании ссылки была ошибка - возвращаем код 409, но в теле указыаем ссылку
	if errCreate != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write(responseJSON)
		return
	}

	// Выводим новую ссылку на экран
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}

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
		shortURL, _ := s.CreateShortURL(storage, element.URL, element.CorrelationID)
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

func getPage(w http.ResponseWriter, r *http.Request) {

	// Получаем короткую ссылку
	//shortURL := chi.URLParam(r, "url")
	shortURL := strings.Trim(string(r.RequestURI), " /")

	// Ищем ссылку в хранилище
	url, err := s.GetURL(storage, shortURL)
	if err != nil {
		http.Error(w, "Not found", http.StatusBadRequest)
	} else {
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func pingPage(w http.ResponseWriter, r *http.Request) {
	err := s.Ping(storage)

	if err != nil {
		http.Error(w, "Connect failed", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func createRouter() *chi.Mux {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := chi.NewRouter()
	r.Use(l.RequestLogger(logger)) // Логгирование
	r.Use(middleware.Compress(5))  // Сжатие ответа
	r.Use(d.Decompress)            // Распаковка сжатого запроса
	r.Post("/", postPage)
	r.Post("/api/shorten", postJSONPage)
	r.Post("/api/shorten/batch", postBatchPage)
	r.Get("/ping", pingPage)
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

	// r передаётся как http.Handler
	http.ListenAndServe(flagRunAddr, r)
}
