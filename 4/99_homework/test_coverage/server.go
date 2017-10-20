package main

import (
	"net/http"
	"strconv"
	"encoding/xml"
	"os"
	"io/ioutil"
	"encoding/json"
	"strings"
	"errors"
	"sort"
)

const (
	CorrectAccessToken = "d41d8cd98f00b204e9800998ecf8427e"
)

var (
	ErrInvalidOrderField = errors.New(ErrorBadOrderField)
)

type UserModel struct {
	Id            int    `xml:"id"`
	Guid          string `xml:"guid"`
	IsActive      bool   `xml:"isActive"`
	Balance       string `xml:"balance"`
	Picture       string `xml:"picture"`
	Age           int    `xml:"age"`
	EyeColor      string `xml:"eyeColor"`
	FirstName     string `xml:"first_name"`
	LastName      string `xml:"last_name"`
	Gender        string `xml:"gender"`
	Company       string `xml:"company"`
	Email         string `xml:"email"`
	Phone         string `xml:"phone"`
	Address       string `xml:"address"`
	About         string `xml:"about"`
	Registered    string `xml:"registered"`
	FavoriteFruit string `xml:"favoriteFruit"`
}

type Users struct {
	List []UserModel `xml:"row"`
}

type QueryParams struct {
	Limit      int
	Offset     int
	Query      string
	OrderField string
	OrderBy    int
}

type ById []User

func (a ById) Len() int {
	return len(a)
}

func (a ById) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ById) Less(i, j int) bool {
	return a[i].Id < a[j].Id
}

type ByAge []User

func (a ByAge) Len() int {
	return len(a)
}

func (a ByAge) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByAge) Less(i, j int) bool {
	return a[i].Age < a[j].Age
}

type ByName []User

func (a ByName) Len() int {
	return len(a)
}

func (a ByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

func checkOrderField(orderField string) bool {
	allowed := []string{
		"Id", "Age", "Name",
	}
	for _, param := range allowed {
		if orderField == param {
			return true
		}
	}
	return false
}

func (qp *QueryParams) getQueryParams(r *http.Request) error {
	var err error
	if limit, ok := r.URL.Query()["limit"]; ok {
		if qp.Limit, err = strconv.Atoi(limit[0]); err != nil {
			return err
		}
	}
	if offset, ok := r.URL.Query()["offset"]; ok {
		if qp.Offset, err = strconv.Atoi(offset[0]); err != nil {
			return err
		}
	}
	if query, ok := r.URL.Query()["query"]; ok {
		qp.Query = query[0]
	}
	if orderField, ok := r.URL.Query()["order_field"]; ok {
		qp.OrderField = orderField[0]
		if qp.OrderField == "" {
			qp.OrderField = "Name"
		} else if !checkOrderField(qp.OrderField) {
			return ErrInvalidOrderField

		}
	}
	if orderBy, ok := r.URL.Query()["order_by"]; ok {
		if qp.OrderBy, err = strconv.Atoi(orderBy[0]); err != nil {
			return err
		}
	}
	return nil
}

func (u *Users) parseUsersXml(fileName string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	xmlFile, err := os.Open(pwd + "/" + fileName)
	defer xmlFile.Close()
	if err != nil {
		return err
	}
	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return err
	}
	if err = xml.Unmarshal(byteValue, u); err != nil {
		return err
	}
	return nil
}

func (u *Users) getSubList(params *QueryParams) {
	if params.Offset < 0 {
		params.Offset = 0
	}
	if params.Limit < 0 {
		params.Limit = 0
	}
	if params.Offset > len(u.List) {
		params.Offset = len(u.List)
	}
	if params.Limit > len(u.List) {
		params.Limit = len(u.List)
	}
	limit := params.Offset + params.Limit
	if limit > len(u.List) {
		limit = len(u.List)
	}
	u.List = u.List[params.Offset:limit]
}

func sortUsers(params *QueryParams, users []User) {
	switch params.OrderBy {
	case OrderByDesc:
		switch params.OrderField {
		case "Id":
			sort.Sort(sort.Reverse(ById(users)))
		case "Age":
			sort.Sort(sort.Reverse(ByAge(users)))
		default:
			sort.Sort(sort.Reverse(ByName(users)))
		}
	case OrderByAsIs:
	case OrderByAsc:
		switch params.OrderField {
		case "Id":
			sort.Sort(ById(users))
		case "Age":
			sort.Sort(ByAge(users))
		default:
			sort.Sort(ByName(users))
		}
	}
}

func (u *Users) getUsers(params *QueryParams) []User {
	users := make([]User, 0)
	for i := 0; i < len(u.List); i++ {
		user := User{
			Id:     u.List[i].Id,
			Name:   u.List[i].FirstName + " " + u.List[i].LastName,
			Age:    u.List[i].Age,
			About:  u.List[i].About,
			Gender: u.List[i].Gender,
		}
		if strings.Contains(user.Name, params.Query) ||
			strings.Contains(user.About, params.Query) ||
			params.Query == "" {
			users = append(users, user)
		}
	}
	sortUsers(params, users)
	return users
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("AccessToken") != CorrectAccessToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	params := new(QueryParams)
	if err := params.getQueryParams(r); err != nil {
		errResp := SearchErrorResponse{
			Error: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		buf, err := json.Marshal(&errResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(buf)
		return
	}
	users := new(Users)
	if err := users.parseUsersXml("dataset.xml"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	users.getSubList(params)
	buf, err := json.Marshal(users.getUsers(params))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(buf)
}
