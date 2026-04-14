package response

import (
	github.com/gin-gonic/gin
)

type Response struct {
	Code    int         `json: code`
	Message string      `json: message`
	Data    interface{} `json: data`
	Error   string      `json: error,omitempty`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    200,
		Message: success,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string, errMsg ...string) {
	err := 
	if len(errMsg) > 0 {
		err = errMsg[0]
	}
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Error:   err,
	})
}

func Paginate(c *gin.Context, list interface{}, total int, page int, pageSize int) {
	Success(c, gin.H{
		list:     list,
		total:    total,
		page:     page,
		pageSize: pageSize,
	})
}
