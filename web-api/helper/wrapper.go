package helper

import (
	"net/http"
	"time"

	"go-unit-test/web-api/module/model/web"

	"github.com/gin-gonic/gin"
)

type response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type paginationResponse struct {
	Success bool         `json:"success"`
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data"`
	Meta    web.MetaData `json:"meta"`
}

type ErrorString struct {
	code    int
	message string
}

func (e ErrorString) Code() int {
	return e.code
}

func (e ErrorString) Error() string {
	return e.message
}

func (e ErrorString) Message() string {
	return e.message
}

func getErrorStatusCode(err error) int {
	errString, ok := err.(*ErrorString)
	if ok {
		return errString.Code()
	}

	return http.StatusInternalServerError
}

func CustomError(code int, msg string) error {
	return &ErrorString{
		code:    code,
		message: msg,
	}
}

type Meta struct {
	Method        string    `json:"method"`
	Url           string    `json:"url"`
	Code          string    `json:"code"`
	ContentLength int64     `json:"content_length"`
	Date          time.Time `json:"date"`
	Ip            string    `json:"ip"`
}

func RespSuccess(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, response{
		Message: message,
		Data:    data,
		Code:    http.StatusOK,
		Success: true,
	})
}

func RespError(c *gin.Context, err error) {
	c.JSON(getErrorStatusCode(err), response{
		Message: err.Error(),
		Data:    nil,
		Code:    getErrorStatusCode(err),
		Success: false,
	})
}

func RespPagination(c *gin.Context, data interface{}, metadata web.MetaData, message string) {
	c.JSON(http.StatusOK, paginationResponse{
		Message: message,
		Meta:    metadata,
		Data:    data,
		Code:    http.StatusOK,
		Success: true,
	})
}
