package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Row struct {
	Name          xml.Name `xml:"row"`
	Id            int      `xml:"id"`
	Guid          string   `xml:"guid"`
	IsActive      string   `xml:"isActive"`
	Balance       string   `xml:"balance"`
	Picture       string   `xml:"picture"`
	Age           int      `xml:"age"`
	EyeColor      string   `xml:"eyeColor"`
	FirstName     string   `xml:"first_name"`
	LastName      string   `xml:"last_name"`
	Gender        string   `xml:"gender"`
	Company       string   `xml:"company"`
	Email         string   `xml:"email"`
	Phone         string   `xml:"phone"`
	Address       string   `xml:"address"`
	About         string   `xml:"about"`
	Registered    string   `xml:"registered"`
	FavoriteFruit string   `xml:"favoriteFruit"`
}

type Result struct {
	Name xml.Name `xml:"root"`
	Rows []Row    `xml:"row"`
}

type Server struct {
	Data       Result
	Tokens     map[string]bool
	FileName   string
	SortFields map[string]bool
	WasInit    bool
}

type CmpArr struct {
	A     []Row
	Query string
	Order string
}

var (
	FileName = "dataset.xml"
	Token = "123"
)

func serializer(r Row) User {
	var res User
	res.Id = r.Id
	res.Name = r.FirstName + " " + r.LastName
	res.Age = r.Age
	res.About = r.About
	res.Gender = r.Gender
	return res
}

func (s *Server) Init(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Can't read file - TEST PASS")
	}
	ok := xml.Unmarshal([]byte(data), &(s.Data))
	if ok != nil {
		fmt.Println("Can't open xml file - TEST PASS")
	}
	s.SortFields["Name"] = true
	s.SortFields["Age"] = true
	s.SortFields["Id"] = true

	s.Tokens[Token] = true
	s.WasInit = true
}

func GetField(r Row, query string) interface{} {
	switch query {
	case "Age":
		{
			return r.Age
		}
	case "Id":
		{
			return r.Id
		}
	default:
		{
			return r.FirstName + r.LastName
		}
	}
}

func (c CmpArr) Len() int {
	return len(c.A)
}

func (c CmpArr) Less(i, j int) bool {
	first := GetField(c.A[i], c.Query)
	second := GetField(c.A[j], c.Query)
	strVal1, ok := first.(string)
	strVal2, ok := second.(string)

	if ok {
		if c.Order == "-1" {
			return strVal1 > strVal2
		}
		return strVal1 < strVal2
	}
	intVal1, _ := first.(int)
	intVal2, _ := second.(int)

	if c.Order == "-1" {
		return intVal1 > intVal2
	}
	return intVal1 < intVal2
}

func (c CmpArr) Swap(i, j int) { c.A[i], c.A[j] = c.A[j], c.A[i] }

func (s *Server) SearchServer(w http.ResponseWriter, r *http.Request) {
	if !s.WasInit {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"Error": "Server is not init"}`))
		return
	}
	_, val := s.Tokens[r.Header.Get("AccessToken")]
	if !val {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "unknown user"}`))
		return
	}

	limit, _ := strconv.Atoi(r.FormValue("limit"))
	query := strings.ToLower(r.FormValue("query"))
	orderField := r.FormValue("order_field")
	order := r.FormValue("order_by")
	offset, _ := strconv.Atoi(r.FormValue("offset"))

	if query == "sleep" {
		time.Sleep(3 * time.Second)
	}

	if orderField == "" {
		orderField = "Name"
	}
	_, ok := s.SortFields[orderField]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Error": "ErrorBadOrderField"}`))
		return
	}

	if order != "-1" && order != "0" && order != "1" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Error": "ErrorOrderBy"}`))
		return
	}

	if offset > len(s.Data.Rows) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c := CmpArr{
		make([]Row, 0),
		orderField,
		order,
	}

	for _, row := range s.Data.Rows {
		if strings.Contains(strings.ToLower(row.FirstName), query) ||
			strings.Contains(strings.ToLower(row.LastName), query) ||
			strings.Contains(strings.ToLower(row.About), query) {
			c.A = append(c.A, row)
		}
	}

	if c.Order != "0" {
		sort.Sort(c)
	}

	users := make([]User, 0)
	cnt := 0
	if offset >= len(c.A) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	for i := offset; i < len(c.A); i++ {
		users = append(users, serializer(c.A[i]))
		cnt++
		if cnt == limit {
			break
		}
	}
	responseData, _ := json.Marshal(users)
	w.Write(responseData)
}

func main() {
}
