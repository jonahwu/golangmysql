package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	//  "strings"
)

const SQLGetGpsLocFromName = `
                              SELECT ST_Longitude(loc),ST_Latitude(loc) FROM gps_data
                              WHERE customer_id = ?`

//curl -X POST  "http://localhost:8080/gpsloc?lat=23.3&lon=123.6"
//curl -X POST  "http://localhost:8080/gpsloc?lat=23.3&lon=123.6&Name=Mary"
func POSTLoc(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		Lat := c.Query("lat")
		Lon := c.Query("lon")
		cid := c.Query("cid")
		fmt.Println(cid)
		if cid == "" {
			log.Println("can not store gps data")
			//	c.String(403, "Need Name information in querystring")
			c.JSON(http.StatusMisdirectedRequest, gin.H{"status": "can not store data"})
			c.Abort()
			return
		}
		//Name := "Mary"
		// can not write to const with not Simply question mark
		// we must to store lat(23), lon(123)
		StoreGPS := fmt.Sprintf("INSERT INTO gps_data (customer_id,loc) VALUES(\"%s\",ST_GeomFromText('POINT(%s %s)',4326))", cid, Lat, Lon)
		fmt.Println(StoreGPS)
		resultu, err := db.Exec(StoreGPS)
		if err != nil {
			fmt.Println("show error")
			log.Fatal(err)
		}
		fmt.Println(resultu)
		fmt.Println("show data")
		//c.String(200, "ret")
		c.JSON(http.StatusOK, gin.H{"status": "store data finished"})
	}
	return gin.HandlerFunc(fn)
}

//curl http://localhost:8080/gpsloc
func GetLoc(db *sql.DB) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var loc GPSLoc
		var locs []GPSLoc
		cid := c.Query("cid")
		rows, err := db.Query(SQLGetGpsLocFromName, cid)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&loc.Lat, &loc.Lon); err != nil {
				log.Fatal(err)
			}
			locs = append(locs, loc)
			fmt.Println("show loc:", loc.Lat, loc.Lon)
		}
		//c.String(200, "ret")
		//c.JSON(http.StatusOK, gin.H{"status": "GetLoc"})
		//c.JSON(http.StatusOK, gin.H{"status": loc})
		c.JSON(http.StatusOK, gin.H{"status": locs})
	}
	return gin.HandlerFunc(fn)
}
