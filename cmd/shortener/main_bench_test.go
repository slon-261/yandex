// Create handler Save load
//
// go test -bench . -v ./... -benchmem
// go tool pprof -http=":9090" bench.test base.pprof

package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkPostPage(b *testing.B) {

	for i := 0; i < b.N; i++ {
		body := "http://benchmark" + strconv.Itoa(i)
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		w := httptest.NewRecorder()
		// вызовем хендлер как обычную функцию, без запуска самого сервера
		postPage(w, r)
	}
}

func BenchmarkPostJSONPage(b *testing.B) {

	for i := 0; i < b.N; i++ {
		body := "{\"url\":\"http://benchmark" + strconv.Itoa(i) + "\"}"
		r := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(body))
		w := httptest.NewRecorder()
		// вызовем хендлер как обычную функцию, без запуска самого сервера
		postJSONPage(w, r)
	}
}

func BenchmarkPostBatchPage(b *testing.B) {

	for i := 0; i < b.N; i++ {
		body := "[{\"correlation_id\": \"" + strconv.Itoa(i) + "\",\"original_url\": \"http://benchmark" + strconv.Itoa(i) + "\"}]"
		r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(body))
		w := httptest.NewRecorder()
		// вызовем хендлер как обычную функцию, без запуска самого сервера
		postBatchPage(w, r)
	}
}

func BenchmarkGetPage(b *testing.B) {

	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest(http.MethodPost, "/"+strconv.Itoa(i), nil)
		w := httptest.NewRecorder()
		// вызовем хендлер как обычную функцию, без запуска самого сервера
		getPage(w, r)
	}
}
