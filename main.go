package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joyouscob/ginapi/internal/controllers"
	"github.com/mattn/go-sqlite3"
)

//GET /product 200
//GET /product/:guid 200
//POST /products 201
//PUT /products/guid 200
//DELETE /product/guid 204

func main() {
	//set gins default address
	var router = gin.Default()
	var address = ":3000"

	sql.Register("sqlite3_with_extensions",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				"sqlite3_mod_regexp",
			},
		})
	//open connection
	var db *sql.DB
	var e error
	if db, e = sql.Open("sqlite3", "./data.db"); e != nil {
		log.Fatalf("Error: %v", e)
	}
	defer db.Close()

	if e := db.Ping(); e != nil {
		log.Fatalf("Error: %v", e)

	}

	router.GET("/products", controllers.GetProducts(db))
	router.GET("/product", controllers.GetProducts(db))

	router.GET("/products/:guid", controllers.GetProduct(db))
	router.POST("/products", controllers.PostProduct(db))
	router.DELETE("/products/:guid", controllers.DeleteProduct(db))
	router.PUT("/products/:guid", controllers.PutProduct(db))

	//log any error from the router
	log.Fatalln(router.Run(address))

}
