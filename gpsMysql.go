package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	//  "strings"
)

const SQLGetGpsLocFromName = `
                              SELECT ST_Longitude(loc),ST_Latitude(loc) FROM gps_data
                              WHERE name = ?`

//curl -X POST  http://localhost:8080/gpsloc
func POSTLoc(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		Lon := "63.3"
		Lat := "123.3"
		Name := "Mary"
		// can not write to const with not Simply question mark
		StoreGPS := fmt.Sprintf("INSERT INTO gps_data (name,loc) VALUES(\"%s\",ST_GeomFromText('POINT(%s %s)',4326))", Name, Lon, Lat)
		fmt.Println(StoreGPS)
		resultu, err := db.Exec(StoreGPS)
		if err != nil {
			fmt.Println("show error")
			log.Fatal(err)
		}
		fmt.Println(resultu)
		fmt.Println("show data")
	}
	return gin.HandlerFunc(fn)
}

//curl http://localhost:8080/gpsloc
func GetLoc(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
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
		c.String(200, "ret")
	}
	return gin.HandlerFunc(fn)
}
