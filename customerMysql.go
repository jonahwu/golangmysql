package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	//  "strings"
)

const SQLGetCustomer = `SELECT BIN_TO_UUID(customer_id) FROM customers WHERE mail = ?`

//curl -X POST  http://localhost:8080/gpsloc
func GetCustomer(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// can not write to const with not Simply question mark
		var customerId string
		mail := c.Query("mail")
		//mail := "junmein@hotmail.com"
		rows, err := db.Query(SQLGetCustomer, mail)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&customerId); err != nil {
				log.Fatal(err)
			}
			fmt.Println("show customer's id:", customerId)
		}
		c.String(200, customerId)
	}
	return gin.HandlerFunc(fn)
}
