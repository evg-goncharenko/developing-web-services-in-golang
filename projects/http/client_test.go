package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func MakeServer() Server {
	return Server{
		Result{},
		make(map[string]bool),
		FileName,
		make(map[string]bool),
		false,
	}
}

func TestCantOpen(t *testing.T) {
	s := MakeServer()
	s.Init("ErrorFile.xml")
	ts := httptest.NewServer(http.HandlerFunc(s.SearchServer))
	cases := []SearchRequest{
		{
			Limit:      1,
			Offset:     1,
			Query:      "",
			OrderField: "",
			OrderBy:    1,
		},
	}

	for _, req := range cases {
		srv := SearchClient{Token, ts.URL}
		w, err := srv.FindUsers(req)
		if w != nil || err.Error() != "cant unpack error json: unexpected end of JSON input" {
			t.Errorf("File Error")
		}
	}

	ts.Close()
}

func TestErrorTimeout(t *testing.T) {
	s := MakeServer()
	s.Init("dataset.xml")
	ts := httptest.NewServer(http.HandlerFunc(s.SearchServer))

	cases := []SearchRequest{
		{
			Limit:      4,
			Offset:     1,
			Query:      "sleep",
			OrderField: "",
			OrderBy:    1,
		},
	}

	for _, req := range cases {
		srv := SearchClient{Token, ts.URL}
		w, err := srv.FindUsers(req)
		if w != nil || err.Error() != "timeout for limit=5&offset=1&order_by=1&order_field=&query=sleep" {
			t.Errorf("Timeout Error")
		}
	}
	ts.Close()
}

func TestErrorUrl(t *testing.T) {
	s := MakeServer()
	s.Init("dataset.xml")
	ts := httptest.NewServer(http.HandlerFunc(s.SearchServer))

	cases := []SearchRequest{
		{
			Limit:      4,
			Offset:     1,
			Query:      "",
			OrderField: "",
			OrderBy:    -2,
		},
	}

	for _, req := range cases {
		srv := SearchClient{Token, "something" + ts.URL}
		w, err := srv.FindUsers(req)
		if w != nil || err.Error() != "unknown error Get \"something"+ts.URL+"?limit=5&offset=1&order_by=-2&order_field=&query=\": unsupported protocol scheme \"somethinghttp\"" {
			t.Errorf("URL Error")
		}
	}
	ts.Close()
}

func TestNotFound(t *testing.T) {
	s := MakeServer()
	s.Init("dataset.xml")
	ts := httptest.NewServer(http.HandlerFunc(s.SearchServer))

	cases := []SearchRequest{
		{
			Limit:      4,
			Offset:     10000,
			Query:      "",
			OrderField: "",
			OrderBy:    1,
		},
		{
			Limit:      4,
			Offset:     10,
			Query:      "nnn",
			OrderField: "",
			OrderBy:    1,
		},
		{
			Limit:      4,
			Offset:     1,
			Query:      "",
			OrderField: "",
			OrderBy:    -2,
		},
	}

	for _, req := range cases {
		srv := SearchClient{Token, ts.URL}
		w, err := srv.FindUsers(req)
		if w != nil {
			fmt.Println(err)
		}
	}

	ts.Close()
}

func TestGetFullInfo(t *testing.T) {
	s := MakeServer()
	s.Init("dataset.xml")
	ts := httptest.NewServer(http.HandlerFunc(s.SearchServer))

	cases := []SearchRequest{
		{
			Limit:      25,
			Offset:     20,
			Query:      "",
			OrderField: "Name",
			OrderBy:    1,
		},
	}

	for _, req := range cases {
		srv := SearchClient{Token, ts.URL}
		w, err := srv.FindUsers(req)
		if w == nil {
			fmt.Println(err)
		}
	}

	ts.Close()
}

func TestGetUserWithBadServerError(t *testing.T) {
	s := MakeServer()
	ts := httptest.NewServer(http.HandlerFunc(s.SearchServer))

	cases := []SearchRequest{
		{
			Limit:      1,
			Offset:     0,
			Query:      "",
			OrderField: "",
			OrderBy:    1,
		},
	}

	for _, req := range cases {
		srv := SearchClient{Token, ts.URL}
		w, err := srv.FindUsers(req)
		if w != nil || err.Error() != "SearchServer fatal error" {
			t.Errorf("Server Error")
		}
	}

	ts.Close()
}

func TestGetUserWithBadReq(t *testing.T) {
	s := MakeServer()
	s.Init("dataset.xml")
	ts := httptest.NewServer(http.HandlerFunc(s.SearchServer))

	cases := []SearchRequest{
		{
			Limit:      1,
			Offset:     0,
			Query:      "",
			OrderField: "BadField",
			OrderBy:    1,
		},
		{
			Limit:      -1,
			Offset:     0,
			Query:      "",
			OrderField: "BadField",
			OrderBy:    1,
		},
		{
			Limit:      26,
			Offset:     0,
			Query:      "",
			OrderField: "BadField",
			OrderBy:    1,
		},
		{
			Limit:      25,
			Offset:     -1,
			Query:      "",
			OrderField: "BadField",
			OrderBy:    1,
		},
	}

	for _, req := range cases {
		srv := SearchClient{Token, ts.URL}
		w, err := srv.FindUsers(req)
		if w != nil {
			fmt.Println(err)
		}
	}

	ts.Close()
}

func TestGetUserWithBadToken(t *testing.T) {
	s := MakeServer()
	s.Init("dataset.xml")
	ts := httptest.NewServer(http.HandlerFunc(s.SearchServer))

	cases := []SearchRequest{
		{
			Limit:      4,
			Offset:     0,
			Query:      "",
			OrderField: "",
			OrderBy:    1,
		},
	}

	for _, req := range cases {
		srv := SearchClient{Token + "UNCORRECT", ts.URL}
		w, err := srv.FindUsers(req)
		if w != nil || err.Error() != "Bad AccessToken" {
			t.Errorf("Token Error")
		}
	}

	ts.Close()
}

func TestGetUser(t *testing.T) {
	s := MakeServer()
	s.Init("dataset.xml")
	ts := httptest.NewServer(http.HandlerFunc(s.SearchServer))

	cases := []SearchRequest{
		{
			Limit:      4,
			Offset:     0,
			Query:      "",
			OrderField: "",
			OrderBy:    1,
		},
		{
			Limit:      4,
			Offset:     0,
			Query:      "",
			OrderField: "Age",
			OrderBy:    1,
		},
		{
			Limit:      4,
			Offset:     0,
			Query:      "",
			OrderField: "Id",
			OrderBy:    1,
		},

		{
			Limit:      4,
			Offset:     0,
			Query:      "",
			OrderField: "Name",
			OrderBy:    -1,
		},
		{
			Limit:      4,
			Offset:     0,
			Query:      "",
			OrderField: "Age",
			OrderBy:    -1,
		},
	}

	for _, req := range cases {
		srv := SearchClient{Token, ts.URL}
		w, err := srv.FindUsers(req)
		if w == nil {
			fmt.Println(err)
		}
	}
	ts.Close()
}
