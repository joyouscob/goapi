package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joyouscob/ginapi/internal"
)

func DeleteProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//create a variable to bind the struct
		//remember we have access to the guidBindgin from getProduct
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

		var result sql.Result
		var e error
		if result, e = db.ExecContext(ctx, "DELETE FROM products WHERE guid=?", binding.GUID); e != nil {
			var res = internal.NewHTTPResponse(http.StatusInternalServerError, e)
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		//lets check if the product was deleted by getting the rows affected
		if nProducts, _ := result.RowsAffected(); nProducts == 0 {
			var res = internal.NewHTTPResponse(http.StatusNotFound, sql.ErrNoRows)
			c.JSON(http.StatusNotFound, res)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
