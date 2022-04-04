package result

import (
	"TodoLists/common/errorCode"
	"TodoLists/common/errorMsg"
	"TodoLists/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

var logger = utils.Logger

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

//定义几个常用的错误，抽取重复代码
func JsonError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Error(errorCode.ParseJsonError, errorMsg.ParseJsonErrorMsg))
	logger.Error("json error:" + err.Error())
	return
}

func MysqlError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Error(errorCode.MysqlError, errorMsg.MysqlErrorMsg))
	logger.Error("mysql error:" + err.Error())
	return
}
