package response

import (
	"godp/pkg/errorcode"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Details []string    `json:"details"`
}

func NewResponse() *Response {
	return &Response{}
}

func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}

	c.JSON(http.StatusOK, Response{
		Code:    errorcode.Success.Code(),
		Message: errorcode.Success.Msg(),
		Data:    data,
		Details: []string{},
	})
}

func Error(c *gin.Context, err error) {
	if err != nil {
		if v, ok := err.(*errorcode.Error); ok {
			response := Response{
				Code:    v.Code(),
				Message: v.Msg(),
				Data:    gin.H{},
				Details: []string{},
			}
			details := v.Details()
			if len(details) > 0 {
				response.Details = details
			}
			c.JSON(v.StatusCode(), response)
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
			Data:    gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    errorcode.Success.Code(),
		Message: errorcode.Success.Msg(),
		Data:    gin.H{},
	})
}

func ErrorCustom(c *gin.Context, errCode int, errMessage string) {
	c.JSON(http.StatusOK, Response{
		Code:    errCode,
		Message: errMessage,
		Data:    gin.H{},
	})
}
