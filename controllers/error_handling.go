package controllers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	StatusCode  int    `json:"status_code"`
	Path        string `json:"path"`
	ErrorSource string `json:"error_source,omitempty"`
	Message     string `json:"message"`
	Err         error  `json:"-"`
}

func HandleErrResponse(c *gin.Context, errResp ErrorResponse) {
	errResp.Path = c.FullPath()
	if errResp.Message != "" {
		errResp.Message = errResp.Err.Error()
	}

	c.JSON(errResp.StatusCode, errResp)
}

func MakeErrResponse(err error, statusCode ...int) ErrorResponse {
	errResp := ErrorResponse{
		Err: err,
	}

	errResp.setStatusOrGetDeafult()

	return errResp
}

func (e *ErrorResponse) setStatusOrGetDeafult(statusCode ...int) {
	e.StatusCode = 400
	if len(statusCode) > 0 {
		e.StatusCode = statusCode[0]
	}
}

func (e *ErrorResponse) Error() string {
	b, _ := json.MarshalIndent(e, "", "\t")
	return string(b)
}
