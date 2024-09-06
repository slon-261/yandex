package main

import (
	"crypto/sha256"
	"encoding/base64"
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
	response := "http://localhost:8080/" + shortURL

	// Выводим новую ссылку на экран
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func getPage(w http.ResponseWriter, r *http.Request) {

	// Получаем короткую ссылку
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

func urlPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		postPage(w, r)
	} else if r.Method == http.MethodGet {
		getPage(w, r)
	} else {
		return
	}
}

func main() {

	http.HandleFunc(`/`, urlPage)

	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}
