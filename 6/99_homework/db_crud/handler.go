package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type MyHander struct {
	DB     *sql.DB
	Router *MyRouter
}

func (mhr *MyHander) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, baseResp := CreateContext(w, r), NewResponse()
	splittedPath := strings.Split(r.URL.Path[1:], "/")
	if !isValidEntity(splittedPath[0]) && len(splittedPath[0]) != 0 {
		baseResp.ServeError(w, http.StatusNotFound, "unknown table")
		return
	}
	for _, route := range mhr.Router.Routes {
		if route.Method == r.Method && ctx.MatchAndFill(route.Path, splittedPath) {
			route.Func(ctx, baseResp)
			return
		}
	}
	baseResp.ServeError(w, http.StatusMethodNotAllowed, "unknown method")
}

func (mhr *MyHander) GetTables(ctx *MyContext, baseResp *MyResponse) {
	rows, err := mhr.DB.Query("SHOW TABLES")
	if err != nil {
		baseResp.ServeError(ctx.Writer, http.StatusNotFound, err.Error())
		return
	}
	tables := make([]string, 0)
	for rows.Next() {
		var table string
		rows.Scan(&table)
		tables = append(tables, table)
	}
	response := make(map[string]interface{})
	response["tables"] = tables
	baseResp.Body["response"] = response
	baseResp.ServeSuccess(ctx.Writer, http.StatusOK)
}

func (mhr *MyHander) GetTableEntries(ctx *MyContext, baseResp *MyResponse) {
	table := ctx.GetPathVar("table")
	sqlQuery := "SELECT * FROM " + table
	var args []interface{}
	if limit, ok := ctx.Request.URL.Query()["limit"]; ok {
		sqlQuery += " LIMIT ?"
		num, _ := strconv.Atoi(limit[0])
		args = append(args, num)
	}
	if offset, ok := ctx.Request.URL.Query()["offset"]; ok {
		sqlQuery += " OFFSET ?"
		num, _ := strconv.Atoi(offset[0])
		args = append(args, num)
	}
	rows, err := mhr.DB.Query(sqlQuery, args...)
	if err != nil {
		baseResp.ServeError(ctx.Writer, http.StatusNotFound, err.Error())
		return
	}
	var records []interface{}
	switch table {
	case "items":
		records = getItems(rows)
	case "users":
		records = getUsers(rows)
	}
	response := make(map[string]interface{})
	response["records"] = records
	baseResp.Body["response"] = response
	baseResp.ServeSuccess(ctx.Writer, http.StatusOK)
}

func getItems(rows *sql.Rows) []interface{} {
	records := make([]interface{}, 0)
	for rows.Next() {
		var item Item
		rows.Scan(&item.ID, &item.Title, &item.Description, &item.Updated)
		records = append(records, item)
	}
	return records
}

func getUsers(rows *sql.Rows) []interface{} {
	records := make([]interface{}, 0)
	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Login, &user.Password, &user.Email, &user.Info, &user.Updated)
		records = append(records, user)
	}
	return records
}

func (mhr *MyHander) GetTableEntry(ctx *MyContext, baseResp *MyResponse) {
	table := ctx.GetPathVar("table")
	id, _ := strconv.Atoi(ctx.GetPathVar("id"))
	sqlQuery := "SELECT * FROM %s WHERE id = ?"
	row := mhr.DB.QueryRow(fmt.Sprintf(sqlQuery, table), id)
	var record interface{}
	var err error
	switch table {
	case "items":
		record, err = getItem(row)
	case "users":
		record, err = getUser(row)
	}
	if err != nil {
		baseResp.ServeError(ctx.Writer, http.StatusNotFound, "record not found")
		return
	}
	response := make(map[string]interface{})
	response["record"] = record
	baseResp.Body["response"] = response
	baseResp.ServeSuccess(ctx.Writer, http.StatusOK)
}

func getItem(row *sql.Row) (interface{}, error) {
	var item Item
	err := row.Scan(&item.ID, &item.Title, &item.Description, &item.Updated)
	return item, err
}

func getUser(row *sql.Row) (interface{}, error) {
	var user User
	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Email, &user.Info, &user.Updated)
	return user, err
}

func (mhr *MyHander) CreateTableEntry(ctx *MyContext, baseResp *MyResponse) {
	table := ctx.GetPathVar("table")
	var id int64
	var err error
	switch table {
	case "items":
		id, err = createItem(mhr.DB, ctx)
	case "users":
		id, err = createUser(mhr.DB, ctx)
	}
	if err != nil {
		baseResp.ServeError(ctx.Writer, http.StatusNotFound, "error creating entity")
		return
	}
	response := make(map[string]interface{})
	response["id"] = id
	baseResp.Body["response"] = response
	baseResp.ServeSuccess(ctx.Writer, http.StatusOK)
}

func createItem(db *sql.DB, ctx *MyContext) (int64, error) {
	sqlQuery := "INSERT INTO items (title, description, updated) VALUES (?, ?, ?)"
	item := new(Item)
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	json.Unmarshal(body, item)
	ctx.Request.Body.Close()
	res, err := db.Exec(sqlQuery, item.Title, item.Description, item.Updated)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, err
}

func createUser(db *sql.DB, ctx *MyContext) (int64, error) {
	sqlQuery := "INSERT INTO users (login, password, email, info, updated) VALUES (?, ?, ?, ?, ?)"
	user := new(User)
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	json.Unmarshal(body, user)
	ctx.Request.Body.Close()
	res, err := db.Exec(sqlQuery, user.Login, user.Password, user.Email, user.Info, user.Updated)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, err
}

func (mhr *MyHander) UpdateTableEntry(ctx *MyContext, baseResp *MyResponse) {
	record := make(map[string]interface{})
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	json.Unmarshal(body, &record)
	ctx.Request.Body.Close()
	table := ctx.GetPathVar("table")
	id, _ := strconv.Atoi(ctx.GetPathVar("id"))
	sqlQuery := fmt.Sprintf("UPDATE %s SET", table)
	for k, v := range record {
		switch v.(type) {
		case float64:
			sqlQuery += fmt.Sprintf(" %s = %f,", k, v)
		case string:
			sqlQuery += fmt.Sprintf(" %s = '%s',", k, v)
		default:
			sqlQuery += fmt.Sprintf(" %s = NULL,", k)
		}
	}
	sqlQuery = sqlQuery[:len(sqlQuery)-1]
	sqlQuery += " WHERE id = ?"
	if _, err := mhr.DB.Exec(sqlQuery, id); err != nil {
		baseResp.ServeError(ctx.Writer, http.StatusNotFound, "error updating entity")
		return
	}
	response := make(map[string]interface{})
	response["updated"] = 1
	baseResp.Body["response"] = response
	baseResp.ServeSuccess(ctx.Writer, http.StatusOK)
}

func (mhr *MyHander) DeleteTableEntry(ctx *MyContext, baseResp *MyResponse) {
	table := ctx.GetPathVar("table")
	id, _ := strconv.Atoi(ctx.GetPathVar("id"))
	sqlQuery := fmt.Sprintf("DELETE FROM %s WHERE id = ?", table)
	res, err := mhr.DB.Exec(sqlQuery, id)
	if err != nil {
		baseResp.ServeError(ctx.Writer, http.StatusNotFound, "error deleting entity")
		return
	}
	n, _ := res.RowsAffected()
	response := make(map[string]interface{})
	response["deleted"] = n
	baseResp.Body["response"] = response
	baseResp.ServeSuccess(ctx.Writer, http.StatusOK)
}
