package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
var storage = s.Storage{}

func postPage(w http.ResponseWriter, r *http.Request) {

	// Получаем ссылку из body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Body error", http.StatusBadRequest)
		return
	}
	originalURL := strings.TrimSpace(string(body))

	// Сохраняем короткую ссылку
	shortURL := storage.CreateShortURL(originalURL)
	response := flagBaseURL + "/" + shortURL

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
	shortURL := storage.CreateShortURL(req.URL)
	var resp models.Response
	resp.Result = flagBaseURL + "/" + shortURL

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
	url, err := storage.GetURL(shortURL)
	if err != nil {
		http.Error(w, "Not found", http.StatusBadRequest)
	} else {
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
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
	// Загружаем из файла все ранее сгенерированные ссылки
	storage.Load(flagFilePath)
	defer storage.Close()

	r := createRouter()

	log.Print("Running server on ", flagRunAddr)
	log.Print("File storage is ", flagFilePath)

	// r передаётся как http.Handler
	http.ListenAndServe(flagRunAddr, r)
}
