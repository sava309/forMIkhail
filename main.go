package main

import (
	"database/sql"
	"fmt"
	_"github.com/mattn/go-sqlite3"
)

func main() {
	database, _ := sql.Open("sqllite3", "./savei.db") //создаём таблицу
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT )")
	statement.Exec()
	//добавляем в таблицу новые параметры имена
	statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	statement.Exec("Lorem", "Ipsum")

	rows, _ := database.Query("SELECT id, firstname, lastname From people")
	var id int //вводим переменные
	var firstname string
	var lastname string
	//ставим указатели
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
		fmt.Printf("%d: %s %s\n", id, firstname, lastname)
	}

}
