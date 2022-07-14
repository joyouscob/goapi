package internal

import (
	"net/http"
)

//response struct
type HTTPResponse struct {
	Status  int         `json:"status"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//creating a factory function to create a new instance of the HTTPResponse struct
func NewHTTPResponse(status int, data interface{}) HTTPResponse {
	switch status {
	case http.StatusBadRequest,
		http.StatusInternalServerError,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusRequestTimeout:

		if e, ok := data.(error); ok {

			return HTTPResponse{
				Status:  status,
				Success: false,
				Message: e.Error(),
			}
		}
		return HTTPResponse{
			Status:  status,
			Success: false,
			Message: data.(string),
		}
	default:
		return HTTPResponse{
			Status:  status,
			Success: true,
			Data:    data,
		}
	}
}
