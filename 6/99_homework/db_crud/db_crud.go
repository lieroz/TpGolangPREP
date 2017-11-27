package main

import (
	"database/sql"
	"net/http"
)

var validEntites []string

func isValidEntity(table string) bool {
	for _, entity := range validEntites {
		if entity == table {
			return true
		}
	}
	return false
}

func NewDbCRUD(db *sql.DB) (http.Handler, error) {
	router := NewRouter()
	handler := &MyHander{DB: db}

	router.RegisterHandler("/", http.MethodGet, handler.GetTables)
	router.RegisterHandler("/$table", http.MethodGet, handler.GetTableEntries)
	router.RegisterHandler("/$table/$id", http.MethodGet, handler.GetTableEntry)
	router.RegisterHandler("/$table", http.MethodPut, handler.CreateTableEntry)
	router.RegisterHandler("/$table/$id", http.MethodPost, handler.UpdateTableEntry)
	router.RegisterHandler("/$table/$id", http.MethodDelete, handler.DeleteTableEntry)

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var table string
		rows.Scan(&table)
		validEntites = append(validEntites, table)
	}

	handler.Router = router
	return handler, nil
}
