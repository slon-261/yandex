package main

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
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
	return chi.NewRouter()
}

func main() {
	// обрабатываем аргументы командной строки
	parseFlags()

	r := createRouter()
	r.Post("/", postPage)
	r.Get("/{url}", getPage)

	log.Print("Running server on ", flagRunAddr)

	// r передаётся как http.Handler
	http.ListenAndServe(flagRunAddr, r)
}
