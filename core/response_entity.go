package core

import "github.com/kataras/iris/v12"

type ResponseEntity struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func ResponseOK(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(ResponseEntity{
		Code:    iris.StatusOK,
		Data:    nil,
		Message: "",
	})
}

func ResponseError(ctx iris.Context, code int, message string) {
	ctx.StatusCode(code)
	ctx.JSON(ResponseEntity{
		Code:    code,
		Data:    nil,
		Message: message,
	})
}
