package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/joyouscob/ginapi/internal"
)

type putProduct struct {
	Name        string  `json:"name" binding:"required_without_all=Price Descrption"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Description string  `json:"description" binding:omitempty,max=250"`
}

func PutProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var binding guidBinding
		var payload putProduct
		var ctx = c.Request.Context()

		//get uri parameter
		if e := c.ShouldBindUri(&binding); e != nil {
			var res = internal.NewHTTPResponse(http.StatusBadRequest, e)
			c.JSON(http.StatusBadRequest, res)
			return

		}

		//bind to json
		if e := c.ShouldBindJSON(&payload); e != nil {
			var res = internal.NewHTTPResponse(http.StatusBadRequest, e)
			c.JSON(http.StatusBadRequest, res)
			return
		}

		var row = db.QueryRowContext(ctx, "SELECT name,price,description FROM products WHERE guid=?", binding.GUID)
		var currentProduct Product
		if e := row.Scan(&currentProduct.Name, &currentProduct.Price, &currentProduct.Description); e != nil {
			//if no row was found
			if e == sql.ErrNoRows {
				var res = internal.NewHTTPResponse(http.StatusNotFound, e)
				c.JSON(http.StatusNotFound, res)
				return
			}
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}
		//setting option for copier package
		var option = copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
		}
		//copying current product to payload
		if e := copier.CopyWithOption(&currentProduct, &payload, option); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.JSON(http.StatusInternalServerError, res)
			return

		}

		//update after copy
		if _, e := db.ExecContext(ctx, "UPDATE products SET name=?,price=?,description=? WHERE guid=?", currentProduct.Name, currentProduct.Price, currentProduct.Description, binding.GUID); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		//get back the product
		//we using the QueryRowContext to query the database and getting the uri data from the binded data binding.GUID
		var updatedRow = db.QueryRowContext(ctx, "SELECT guid,name,price,description,createdAt FROM products WHERE guid=?", binding.GUID)
		var product Product
		if e := updatedRow.Scan(&product.GUID, &product.Name, &product.Price, &product.Description, &product.CreatedAt); e != nil {

			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}
		var res = internal.NewHTTPResponse(http.StatusOK, product)
		c.JSON(http.StatusOK, res)

		fmt.Println("Product updated", currentProduct)
		// c.String(http.StatusOK, "put ptoduct")

	}
}
