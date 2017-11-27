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
	sqlQuery := fmt.Sprintf("SELECT * FROM %s", table)
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
	cols, _ := rows.Columns()
	var records []interface{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		rows.Scan(columnPointers...)
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		for k, v := range m {
			if v != nil {
				if k == "id" {
					m[k], _ = strconv.Atoi(string(v.([]uint8)))
				} else {
					m[k] = string(v.([]uint8))
				}
			} else {
				m[k] = nil
			}
		}
		records = append(records, m)
	}
	response := make(map[string]interface{})
	response["records"] = records
	baseResp.Body["response"] = response
	baseResp.ServeSuccess(ctx.Writer, http.StatusOK)
}

func (mhr *MyHander) GetTableEntry(ctx *MyContext, baseResp *MyResponse) {
	table := ctx.GetPathVar("table")
	id, _ := strconv.Atoi(ctx.GetPathVar("id"))
	sqlQuery := "SELECT * FROM %s WHERE id = ?"
	row := mhr.DB.QueryRow(fmt.Sprintf(sqlQuery, table), id)
	rows, _ := mhr.DB.Query(fmt.Sprintf(sqlQuery, table), id) // Needed to get columns information, can also be done in NewCrudDB
	cols, _ := rows.Columns()
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}
	if err := row.Scan(columnPointers...); err != nil {
		baseResp.ServeError(ctx.Writer, http.StatusNotFound, "record not found")
		return
	}
	m := make(map[string]interface{})
	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		m[colName] = *val
	}
	for k, v := range m {
		if v != nil {
			if k == "id" {
				m[k], _ = strconv.Atoi(string(v.([]uint8)))
			} else {
				m[k] = string(v.([]uint8))
			}
		} else {
			m[k] = nil
		}
	}
	response := make(map[string]interface{})
	response["record"] = m
	baseResp.Body["response"] = response
	baseResp.ServeSuccess(ctx.Writer, http.StatusOK)
}

func (mhr *MyHander) CreateTableEntry(ctx *MyContext, baseResp *MyResponse) {
	table := ctx.GetPathVar("table")
	data := make(map[string]interface{})
	defer ctx.Request.Body.Close()
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	json.Unmarshal(body, &data)
	sqlQuery := fmt.Sprintf("INSERT INTO %s (", table)
	var args []interface{}
	for k, v := range data {
		if k == "id" {
			continue
		}
		sqlQuery += fmt.Sprintf("%s,", k)
		args = append(args, v)
	}
	sqlQuery = sqlQuery[:len(sqlQuery)-1]
	sqlQuery += ") VALUES ("
	for i := 0; i < len(args); i++ {
		sqlQuery += "?,"
	}
	sqlQuery = sqlQuery[:len(sqlQuery)-1]
	sqlQuery += ")"
	res, err := mhr.DB.Exec(sqlQuery, args...)
	if err != nil {
		baseResp.ServeError(ctx.Writer, http.StatusNotFound, "error creating entity")
		return
	}
	id, _ := res.LastInsertId()
	response := make(map[string]interface{})
	response["id"] = id
	baseResp.Body["response"] = response
	baseResp.ServeSuccess(ctx.Writer, http.StatusOK)
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
		if k == "id" {
			baseResp.ServeError(ctx.Writer, http.StatusBadRequest, fmt.Sprintf("field %s have invalid type", k))
			return
		}
		switch v.(type) {
		case float64:
			if k == "updated" || k == "title" {
				baseResp.ServeError(ctx.Writer, http.StatusBadRequest, fmt.Sprintf("field %s have invalid type", k))
				return
			}
			sqlQuery += fmt.Sprintf(" %s = %f,", k, v)
		case string:
			sqlQuery += fmt.Sprintf(" %s = '%s',", k, v)
		default:
			if k == "title" {
				baseResp.ServeError(ctx.Writer, http.StatusBadRequest, fmt.Sprintf("field %s have invalid type", k))
				return
			}
			sqlQuery += fmt.Sprintf(" %s = NULL,", k)
		}
	}
	sqlQuery = sqlQuery[:len(sqlQuery)-1]
	sqlQuery += " WHERE id = ?"
	if _, err := mhr.DB.Exec(sqlQuery, id); err != nil {
		baseResp.ServeError(ctx.Writer, http.StatusBadRequest, "error updating entry")
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
