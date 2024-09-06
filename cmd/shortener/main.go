package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

// Таблица со ссылками
var tableURL = make(map[string]string)

func postPage(w http.ResponseWriter, r *http.Request) {

	// Получаем ссылку из body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url := strings.TrimSpace(string(body))

	// Генерируем короткую ссылку
	h := sha256.New()
	h.Write([]byte(url))
	hashString := base64.StdEncoding.EncodeToString(h.Sum(nil))
	shortURL := hashString[:10]

	// Сохраняем короткую ссылку
	tableURL[shortURL] = url
	response := flagBaseUrl + "/" + shortURL

	// Выводим новую ссылку на экран
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func getPage(w http.ResponseWriter, r *http.Request) {

	// Получаем короткую ссылку
	shortURL := chi.URLParam(r, "url")

	// Ищем ссылку в таблице
	url, ok := tableURL[shortURL]
	if ok {
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	// обрабатываем аргументы командной строки
	parseFlags()

	r := chi.NewRouter()
	r.Post("/", postPage)
	r.Get("/{url}", getPage)

	fmt.Println("Running server on", flagRunAddr)

	// r передаётся как http.Handler
	http.ListenAndServe(flagRunAddr, r)
}
