package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
)

var tableUrl = make(map[string]string)

func postPage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url := strings.TrimSpace(string(body))

	h := md5.New()
	h.Write([]byte(url))
	shortUrl := hex.EncodeToString(h.Sum(nil))

	tableUrl[shortUrl] = url
	response := "http://localhost:8080/" + shortUrl

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func getPage(w http.ResponseWriter, r *http.Request) {
	shortUrl := strings.Trim(string(r.RequestURI), " /")

	url, ok := tableUrl[shortUrl]
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
