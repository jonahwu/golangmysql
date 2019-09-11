package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	//	"strings"
)

type Person struct {
	Id        int    `json:"id" form:"id"`
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
}

type GPSLoc struct {
	Lon float64 `json:"lon" form:"lon"`
	Lat float64 `json:"lat" form:"lat"`
}

const SQLgetCustomerWLastName = `SELECT * FROM customer 
                              WHERE lastname = ?`

const SQLgetProduct = `SELECT customer.id,customer.firstname,customer.lastname,product.productname 
                    FROM customer 
					JOIN product 
					WHERE customer.id = product.id`

const SQLUpdateProduct = `UPDATE product 
                      SET productname =? 
					  WHERE id = ?`
const SQLInsertCustomer = `
           INSERT INTO customer (firstname, lastname) VALUES (?, ?)
          `

func SomeHandler(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Your handler code goes in here - e.g.
		rows, err := db.Query(SQLgetCustomerWLastName, "Bill")
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

		c.String(200, "ret")
	}
	return gin.HandlerFunc(fn)
}

//curl http://localhost:8080/qs?test=23&&test1=lala
//curl "http://localhost:8080/qs?test=23&test1=lala"  need "
func QueryString(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// Your handler code goes in here - e.g.
		q1 := c.Query("test")
		fmt.Println(q1)
		q2 := c.Query("test1")
		fmt.Println(q2)
		c.String(200, "ret")
	}

	return gin.HandlerFunc(fn)
}

func main() {
	//fmt.Println("vim-go")
	//where root is user, promise is password, /test is database name 127.0.0.1:3306 is mysql location
	db, err := sql.Open("mysql", "root:promise123@tcp(127.0.0.1:3306)/gps?parseTime=true")
	defer db.Close()
	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxIdleConns(20000)
	db.SetMaxOpenConns(20000)
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	//rows, err := db.Query("SELECT * FROM customer")
	//rows, err := db.Query("SELECT * FROM customer WHERE id = ?", 1)
	//rows, err := db.Query("SELECT * FROM customer WHERE lastname = ?", "Bill")
	rows, err := db.Query(SQLGetGpsLocFromName, "Eason")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var loc GPSLoc

		if err := rows.Scan(&loc.Lat, &loc.Lon); err != nil {
			log.Fatal(err)
		}
		fmt.Println("show loc:", loc.Lat, loc.Lon)
	}

	router := gin.Default()
	router.GET("/test", SomeHandler(db))
	router.GET("/qs", QueryString(db))
	router.GET("/gpsloc", GetLoc(db))
	router.POST("/gpsloc", POSTLoc(db))
	router.GET("/customer", GetCustomer(db))
	router.Run(":8080")
}
