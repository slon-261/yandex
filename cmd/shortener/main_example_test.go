package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	s "slon-261/yandex/internal/storage"
	"strings"
)

func Example_postPage() {

	storage = s.NewStorage(flagDataBaseDSN, flagFilePath)
	// Загружаем из файла\БД все ранее сгенерированные ссылки
	s.Load(storage)
	defer s.Close(storage)

	body := "https://practicum.yandex.ru/"
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	w := httptest.NewRecorder()
	// вызовем хендлер как обычную функцию, без запуска самого сервера
	postPage(w, r)

	fmt.Println(w.Body.String())

	// Output:
	// http://localhost:8080/QrPnX5IUXS
}

func Example_postJsonPage() {

	body := "{\"url\":\"https://practicum.yandex.ru/JSON\"}"
	r := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(body))
	w := httptest.NewRecorder()
	// вызовем хендлер как обычную функцию, без запуска самого сервера
	postJSONPage(w, r)

	fmt.Println(w.Body.String())

	// Output:
	// {
	//    "result": "http://localhost:8080/3ABEBnUYiI"
	// }
}

func Example_postBatchPage() {

	body := "[{\"correlation_id\": \"qqq\",\"original_url\": \"http://du2mkj9ffffffffffff\"},{\"correlation_id\": \"www\",\"original_url\": \"http://e1.ru\"}]"
	r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(body))
	w := httptest.NewRecorder()
	// вызовем хендлер как обычную функцию, без запуска самого сервера
	postBatchPage(w, r)

	fmt.Println(w.Body.String())

	// Output:
	// [
	//    {
	//       "short_url": "http://localhost:8080/D3pOAbtqFc",
	//       "correlation_id": "qqq"
	//    },
	//    {
	//       "short_url": "http://localhost:8080/FYyo4hlW2g",
	//       "correlation_id": "www"
	//    }
	// ]
}

func Example_getPage() {

	r := httptest.NewRequest(http.MethodPost, "/QrPnX5IUXS", nil)
	w := httptest.NewRecorder()
	// вызовем хендлер как обычную функцию, без запуска самого сервера
	getPage(w, r)

	fmt.Println(w.Code)

	// Output:
	// 307
}
