/**
 * @Author : nopsky
 * @Email : cnnopsky@gmail.com
 * @Date : 2021/9/18 11:02
 */
package parser

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func GetCreateTableFromDB(dsn string, tables []string) ([]string, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.WithMessage(err, "open db error")
	}
	defer db.Close()

	if len(tables) == 0 {
		tables, err = getTables(db)
		if err != nil {
			return nil, err
		}
	}

	var rows *sql.Rows
	var createSQL []string
	for _, tableName := range tables {
		rows, err = db.Query("SHOW CREATE TABLE " + tableName)

		if err != nil {
			rows.Close()
			return nil, errors.WithMessage(err, "query show create table error")
		}
		if !rows.Next() {
			rows.Close()
			return nil, errors.Errorf("table(%s) not found", tableName)
		}
		var table string
		var createSql string
		err = rows.Scan(&table, &createSql)
		if err != nil {
			rows.Close()
			return nil, err
		}
		createSQL = append(createSQL, createSql)
		rows.Close()
	}
	return createSQL, nil
}

func getTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW TABLES ")

	if err != nil {
		return nil, errors.WithMessage(err, "query show create table error")
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}
	return tables, nil
}
