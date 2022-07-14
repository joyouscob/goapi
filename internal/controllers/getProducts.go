package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joyouscob/ginapi/internal"
)

// type Product struct {
// 	GUID        string  `json:"guid" binding:"required"`
// 	Name        string  `json:"name" binding:"required"`
// 	Price       float64 `json:"price" binding:"required,gt=0"`
// 	Description string  `json:"description" binding:"omitempty,max=250"`
// 	CreatedAt   string  `json:"createdAt"`
// }

func GetProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rows *sql.Rows
		var e error
		if rows, e = db.Query("SELECT guid,name,price,description,createdAt FROM products"); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		defer rows.Close()
		var products []Product
		for rows.Next() {
			var product Product
			if e := rows.Scan(&product.GUID, &product.Name, &product.Price, &product.Description, &product.CreatedAt); e != nil {
				var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
				c.JSON(http.StatusInternalServerError, res)
				return
			}
			products = append(products, product)
		}

		if len(products) == 0 {
			var res = internal.NewHTTPResponse(http.StatusNotFound, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		var res = internal.NewHTTPResponse(http.StatusOK, products)
		c.JSON(http.StatusOK, res)
	}
}
