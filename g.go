package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Person struct {
	Id        int    `json:"id" form:"id"`
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
}

const getCustomerWLastName = "SELECT * FROM customer WHERE lastname = ?"

func main() {
	//fmt.Println("vim-go")
	//where root is user, promise is password, /test is database name 127.0.0.1:3306 is mysql location
	db, err := sql.Open("mysql", "root:promise@tcp(127.0.0.1:3306)/test?parseTime=true")
	defer db.Close()
	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	//rows, err := db.Query("SELECT * FROM customer")
	//rows, err := db.Query("SELECT * FROM customer WHERE id = ?", 1)
	//rows, err := db.Query("SELECT * FROM customer WHERE lastname = ?", "Bill")
	rows, err := db.Query(getCustomerWLastName, "Bill")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var person Person

		if err := rows.Scan(&person.Id, &person.FirstName, &person.LastName); err != nil {
			log.Fatal(err)
		}
		fmt.Println("show:", person.Id, person.FirstName, person.LastName)
	}

	result, err := db.Exec(
		"INSERT INTO customer (firstname, lastname) VALUES (?, ?)",
		"syhlion",
		"bill",
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	/*
			err := db.QueryRow(`
		    SELECT
		        db1.users.username
		    FROM
		        db1.users
		    JOIN
		        db2.comments
		        ON db1.users.id = db2.comments.username_id
		`).Scan(&username)
	*/

}
