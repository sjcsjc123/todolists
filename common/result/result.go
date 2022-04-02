package result

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

//通用返回成功结果
func Success(data interface{}) interface{} {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Slice && value.IsNil() {
		data = gin.H{}
	}
	return &gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	}
}

//通用返回错误结果
func Error(errorCode int, errorMsg string) interface{} {
	return &gin.H{
		"code": errorCode,
		"msg":  errorMsg,
	}
}
