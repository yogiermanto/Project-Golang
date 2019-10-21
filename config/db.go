package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DbConn() (Db *sql.DB) {
	DbDriver := "mysql"
	DbUser := "root"
	DbPass := ""
	DbName := "project_antrian"
	Db, err := sql.Open(DbDriver, DbUser+":"+DbPass+"@/"+DbName)
	if err != nil {
		panic(err.Error())
	}
	return Db
}
