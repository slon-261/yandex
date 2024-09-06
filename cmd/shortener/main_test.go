package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlPage(t *testing.T) {

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
		{method: http.MethodGet, target: "/QrPnX5IUXS", expectedCode: http.StatusTemporaryRedirect, expectedBody: ""},
		{method: http.MethodGet, target: "/111", expectedCode: http.StatusBadRequest, expectedBody: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.target, strings.NewReader(tc.body))
			w := httptest.NewRecorder()

			// вызовем хендлер как обычную функцию, без запуска самого сервера
			urlPage(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			// проверим корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				// assert.JSONEq помогает сравнить две JSON-строки
				assert.Equal(t, tc.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}
