package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/joyouscob/ginapi/internal"

	"time"

	"github.com/gin-gonic/gin"
)

//this is the struct we are using
//for the request
//notice it does not have createdat and guid
//cos we dont need these info coming from the request
type postProduct struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Description string  `json:"description" binding:"omitempty,max=250"`
}

//the reason for another struct is to
//use it to retrieve the products
type Product struct {
	GUID        string  `json:"guid" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Description string  `json:"description" binding:"omitempty,max=250"`
	CreatedAt   string  `json:"createdAt"`
}

func PostProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//assign payload to the postProduct struct
		var payload postProduct
		//we assign the request data to ctx, could be anything
		var ctx = c.Request.Context()

		//bind payload pointer to e. and if e is not nil then
		//we have an error, cos maybe nothing should be returned
		//returning does not automatically return a response
		//The .bindjson method is what will handle sending the error response.
		if e := c.ShouldBindJSON(&payload); e != nil {
			//remeber response is from internal module/folder
			var res = internal.NewHTTPResponse(http.StatusBadRequest, e)
			c.JSON(http.StatusBadRequest, res)
			return
		}
		//guid from the installed google guid
		var guid = uuid.New().String()
		//created time from the time package
		var createdAt = time.Now().Format(time.RFC3339)
		//here we are excuting the database query
		//if e is not null then there is an error, internalserver error
		if _, e := db.ExecContext(ctx, "INSERT INTO products(guid, name, price, description, createdAt) VALUES(?,?,?,?,?)", guid, payload.Name, payload.Price, payload.Description, createdAt); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}
		//the reason why are assigning this to the product struct is to get back the product
		//so assign product to product struct
		var product Product
		//We running a query using db.QueryRow with sql query using the product guid
		var row = db.QueryRow("SELECT guid,name,price,description,createdAt FROM products WHERE guid=?", guid)
		//AND SCAN THE ROW, i dont know the shit how to explain this
		//if e is not nil, then an error is returned
		//Use scan to map query to a particular go struct
		if e := row.Scan(&product.GUID, &product.Name, &product.Price, &product.Description, &product.CreatedAt); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		//if none of the Es returned error, then we are successful
		//our response can then be the created status
		var res = internal.NewHTTPResponse(http.StatusCreated, product)
		//we also can write the product location to the header
		//said its a good practice
		c.Writer.Header().Add("Location", fmt.Sprintf("/products/%s", guid))
		//
		c.JSON(http.StatusCreated, res)
	}
}
