package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase1 struct {
	ID         string
	Response   string
	StatusCode int
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("id")
	if key == "42" {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": 200, "resp": {"user": 42}}`)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status": 500, "err": "db_error"}`)
	}
}

func TestGetUser(t *testing.T) { // табличное тестирование
	cases := []TestCase1{
		TestCase1{
			ID:         "42",
			Response:   `{"status": 200, "resp": {"user": 42}}`,
			StatusCode: http.StatusOK,
		},
		TestCase1{
			ID:         "500",
			Response:   `{"status": 500, "err": "db_error"}`,
			StatusCode: http.StatusInternalServerError,
		},
	}
	for caseNum, item := range cases {
		url := "http://example.com/api/user?id=" + item.ID

		// "заглушки":
		req := httptest.NewRequest("GET", url, nil) // создаем новый реквест
		w := httptest.NewRecorder()                 // сюда пишем результат

		GetUser(w, req) // аналог http запроса

		if w.Code != item.StatusCode { // тесты
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		bodyStr := string(body)
		if bodyStr != item.Response { // сравниваем результаты
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				caseNum, bodyStr, item.Response)
		}
	}
}
