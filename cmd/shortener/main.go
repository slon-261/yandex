package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	l "slon-261/yandex/internal/logger"
	"slon-261/yandex/internal/models"
	"strings"
)

// Хэш-таблица со ссылками
var tableURL = make(map[string]string)

func postPage(w http.ResponseWriter, r *http.Request) {

	// Получаем ссылку из body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url := strings.TrimSpace(string(body))

	shortURL := encryption(url)

	// Сохраняем короткую ссылку
	tableURL[shortURL] = url
	response := flagBaseURL + "/" + shortURL

	// Выводим новую ссылку на экран
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func postJsonPage(w http.ResponseWriter, r *http.Request) {

	// Получаем ссылку из body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.Request
	if err = json.Unmarshal([]byte(body), &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := encryption(req.Url)

	// Сохраняем короткую ссылку
	tableURL[shortURL] = req.Url
	var resp models.Response
	resp.Result = flagBaseURL + "/" + shortURL

	responseJSON, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	// Ищем ссылку в таблице
	url, ok := tableURL[shortURL]
	if ok {
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
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

func createRouter() *chi.Mux {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := chi.NewRouter()
	r.Use(l.RequestLogger(logger))
	r.Post("/", postPage)
	r.Post("/api/shorten", postJsonPage)
	r.Get("/{url}", getPage)
	return r
}

func main() {
	// обрабатываем аргументы командной строки
	parseFlags()

	r := createRouter()

	log.Print("Running server on ", flagRunAddr)

	// r передаётся как http.Handler
	http.ListenAndServe(flagRunAddr, r)
}
