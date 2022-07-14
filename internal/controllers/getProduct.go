package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joyouscob/ginapi/internal"
)

//create abinding for the guid uri
// we will be using a uri binding
type guidBinding struct {
	GUID string `uri:"guid" binding:"required,uuid4"`
}

func GetProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//create a variable to bind the struct
		var binding guidBinding
		//variable to get data from the request
		var ctx = c.Request.Context()
		//lets bind uri parameter to pointer binding
		//if it fails an error is returned ortherwise its nil
		if e := c.ShouldBindUri(&binding); e != nil {
			//remember response is from internal module/folder
			var res = internal.NewHTTPResponse(http.StatusBadRequest, e)
			c.JSON(http.StatusBadRequest, res)
			return
		}

		//we using the QueryRowContext to query the database and getting the uri data from the binded data binding.GUID
		var row = db.QueryRowContext(ctx, "SELECT guid,name,price,description,createdAt FROM products WHERE guid=?", binding.GUID)
		var product Product
		if e := row.Scan(&product.GUID, &product.Name, &product.Price, &product.Description, &product.CreatedAt); e != nil {
			if e == sql.ErrNoRows {
				var res = internal.NewHTTPResponse(http.StatusNotFound, e)
				c.JSON(http.StatusNotFound, res)
				return
			}
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}
		var res = internal.NewHTTPResponse(http.StatusOK, product)
		c.JSON(http.StatusOK, res)
	}
}
