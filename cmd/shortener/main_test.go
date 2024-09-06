package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestPostPage(t *testing.T) {

	// тип http.HandlerFunc реализует интерфейс http.Handler
	// это поможет передать хендлер тестовому серверу
	handler := http.HandlerFunc(postPage)
	// запускаем тестовый сервер, будет выбран первый свободный порт
	srv := httptest.NewServer(handler)
	// останавливаем сервер после завершения теста
	defer srv.Close()

	// Парсим флаги (в том числе, чтобы задать flagBaseUrl
	parseFlags()

	// описываем набор данных: метод запроса, ожидаемый код ответа, ожидаемое тело
	testCases := []struct {
		method       string
		target       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodPost, target: "/", body: "https://practicum.yandex.ru/", expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/QrPnX5IUXS"},
		{method: http.MethodPost, target: "/", body: "https://practicum.yandex.ru/test", expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/50K3Dd+Erq"},
		{method: http.MethodPost, target: "/", body: "https://e1.ru/", expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/QpZyjSjq5e"},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			// делаем запрос с помощью библиотеки resty к адресу запущенного сервера,
			// который хранится в поле URL соответствующей структуры
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + tc.target
			req.Body = tc.body

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, string(resp.Body()))
			}
		})
	}
}

func TestGetPage(t *testing.T) {

	// тип http.HandlerFunc реализует интерфейс http.Handler
	// это поможет передать хендлер тестовому серверу
	handler := http.HandlerFunc(getPage)
	// запускаем тестовый сервер, будет выбран первый свободный порт
	srv := httptest.NewServer(handler)
	// останавливаем сервер после завершения теста
	defer srv.Close()

	// описываем набор данных: метод запроса, ожидаемый код ответа, ожидаемое тело
	testCases := []struct {
		method       string
		target       string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodGet, target: "/QrPnX5IUXS", expectedCode: http.StatusTemporaryRedirect, expectedBody: ""},
		{method: http.MethodGet, target: "/111", expectedCode: http.StatusBadRequest, expectedBody: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			// делаем запрос с помощью библиотеки resty к адресу запущенного сервера,
			// который хранится в поле URL соответствующей структуры
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL + tc.target

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, string(resp.Body()))
			}
		})
	}
}
