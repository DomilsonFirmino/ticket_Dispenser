package main

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/glebarez/go-sqlite"
)

var mu sync.RWMutex

// bloqueio para impedir escrita multipla e leitura multipla
var pessoas = map[int]person{
	0: {
		Id:        "0",
		Firstname: "domilson",
		Lastname:  "firmino",
		Age:       10,
	},
}

func connectToDatabase() {
	db, err := sql.Open("sqlite", "./storage/my.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	fmt.Println("Connected to the database successfully!")

	//Create table
	_, err = createTable(db)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}
	fmt.Println("Table created successfully!")
	fmt.Println("Database connection closed.")

	//insert person
	p := &person{
		Firstname: "John",
		Lastname:  "Doe",
		Age:       30,
	}
	result, err := insertPerson(db, p)
	if err != nil {
		fmt.Println("Error inserting person:", err)
		return
	}
	id, pe := result.LastInsertId()
	fmt.Printf("Inserted person with ID %d\n", id, pe)
}

func createTable(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS person (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		firstname TEXT NOT NULL,
		lastname TEXT NOT NULL,
		age INTEGER NOT NULL
	);`

	return db.Exec(sql)
}

func insertPerson(db *sql.DB, p *person) (sql.Result, error) {
	sql := `INSERT INTO person (firstname, lastname, age) VALUES (?, ?, ?)`
	result, err := db.Exec(sql, p.Firstname, p.Lastname, p.Age)
	if err != nil {
		return nil, err
	}
	return result, nil
}

var pessoaCounter int = 1
