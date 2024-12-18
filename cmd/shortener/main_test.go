package main

import (
	"net/http"
	"net/http/httptest"
	s "slon-261/yandex/internal/storage"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostPage(t *testing.T) {

	// Парсим флаги (в том числе, чтобы задать flagBaseURL)
	parseFlags()

	storage = s.NewStorage(flagDataBaseDSN, flagFilePath)
	// Загружаем из файла\БД все ранее сгенерированные ссылки
	s.Load(storage)
	defer s.Close(storage)

	// описываем набор данных: метод запроса, ожидаемый код ответа, ожидаемое тело
	testCases := []struct {
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodPost, body: "https://practicum.yandex.ru/", expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/QrPnX5IUXS"},
		{method: http.MethodPost, body: "https://practicum.yandex.ru/test", expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/50K3Dd+Erq"},
		{method: http.MethodPost, body: "https://practicum.yandex.ru/test", expectedCode: http.StatusConflict, expectedBody: "http://localhost:8080/50K3Dd+Erq"},
		{method: http.MethodPost, body: "https://e1.ru/", expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/QpZyjSjq5e"},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/", strings.NewReader(tc.body))
			w := httptest.NewRecorder()

			// вызовем хендлер как обычную функцию, без запуска самого сервера
			postPage(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			// проверим корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				// assert.JSONEq помогает сравнить две JSON-строки
				assert.Equal(t, tc.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}

func TestPostJsonPage(t *testing.T) {

	// описываем набор данных: метод запроса, ожидаемый код ответа, ожидаемое тело
	testCases := []struct {
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodPost, body: "{\"url\":\"https://practicum.yandex.ru/JSON\"}", expectedCode: http.StatusCreated, expectedBody: "{\n   \"result\": \"http://localhost:8080/3ABEBnUYiI\"\n}"},
		{method: http.MethodPost, body: "{\"url\":\"https://practicum.yandex.ru/JSON\"}", expectedCode: http.StatusConflict, expectedBody: "{\n   \"result\": \"http://localhost:8080/3ABEBnUYiI\"\n}"},
		{method: http.MethodPost, body: "{\"url\":\"https://practicum.yandex.ru/testJSON\"}", expectedCode: http.StatusCreated, expectedBody: "{\n   \"result\": \"http://localhost:8080/qNk2xxBG+2\"\n}"},
		{method: http.MethodPost, body: "{\"url\":\"https://e1.ru/\"}", expectedCode: http.StatusConflict, expectedBody: "{\n   \"result\": \"http://localhost:8080/QpZyjSjq5e\"\n}"},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/api/shorten", strings.NewReader(tc.body))
			w := httptest.NewRecorder()

			// вызовем хендлер как обычную функцию, без запуска самого сервера
			postJSONPage(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			// проверим корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				// assert.JSONEq помогает сравнить две JSON-строки
				assert.Equal(t, tc.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}

func TestPostBatchPage(t *testing.T) {

	// описываем набор данных: метод запроса, ожидаемый код ответа, ожидаемое тело
	testCases := []struct {
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodPost, body: "[{\"correlation_id\": \"qqq\",\"original_url\": \"http://du2mkj9ffffffffffff\"},{\"correlation_id\": \"www\",\"original_url\": \"http://e1.ru\"}]", expectedCode: http.StatusCreated, expectedBody: "[\n   {\n      \"short_url\": \"http://localhost:8080/D3pOAbtqFc\",\n      \"correlation_id\": \"qqq\"\n   },\n   {\n      \"short_url\": \"http://localhost:8080/FYyo4hlW2g\",\n      \"correlation_id\": \"www\"\n   }\n]"},
		{method: http.MethodPost, body: "[{\"correlggggggggation_id\": \"qqq\",\"original_url\": \"http://du2mkj9ffffffffffff\"}}]", expectedCode: http.StatusBadRequest, expectedBody: "JSON error\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/api/shorten/batch", strings.NewReader(tc.body))
			w := httptest.NewRecorder()

			// вызовем хендлер как обычную функцию, без запуска самого сервера
			postBatchPage(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			// проверим корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				// assert.JSONEq помогает сравнить две JSON-строки
				assert.Equal(t, tc.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}

func TestGetPage(t *testing.T) {

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
			r := httptest.NewRequest(tc.method, tc.target, nil)
			w := httptest.NewRecorder()

			// вызовем хендлер как обычную функцию, без запуска самого сервера
			getPage(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			// проверим корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				// assert.JSONEq помогает сравнить две JSON-строки
				assert.Equal(t, tc.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}
