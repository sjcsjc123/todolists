package api

import (
	"TodoLists/common/errorCode"
	"TodoLists/common/errorMsg"
	"TodoLists/common/model"
	"TodoLists/common/result"
	"TodoLists/database"
	"TodoLists/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type LoginRequestBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequestBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

func Login(c *gin.Context) {
	var requestBody LoginRequestBody
	err := c.ShouldBindWith(&requestBody, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.ParseJsonError, errorMsg.ParseJsonErrorMsg))
		logger.Error("json error:" + err.Error())
		return
	}
	var user model.User
	var count int64
	err = database.DB.Model(&model.User{}).
		Where("username = ? and password = ?", requestBody.Username, requestBody.Password).
		Count(&count).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.MysqlError, errorMsg.MysqlErrorMsg))
		logger.Error("mysql error:" + err.Error())
		return
	}
	token, err := utils.GenerateToken(user.UserId)
	userId := strconv.Itoa(user.UserId)
	database.RedisClient.Set(context.Background(), "token"+userId, token, time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.JwtToken, errorMsg.JwtTokenMsg))
		logger.Error("jwt error:" + err.Error())
	}
	c.JSON(http.StatusOK, result.Success("login success"))
}

func Register(c *gin.Context) {

	//绑定请求体
	var requestBody RegisterRequestBody
	err := c.ShouldBindWith(&requestBody, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.ParseJsonError, errorMsg.ParseJsonErrorMsg))
		logger.Error("json error:" + err.Error())
		return
	}

	//查找用户名是否唯一
	var countUsername int64
	err = database.DB.Model(&model.User{}).Where("username = ?", requestBody.Username).Count(&countUsername).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.MysqlError, errorMsg.MysqlErrorMsg))
		logger.Error("mysql error:" + err.Error())
		return
	}
	if countUsername >= 1 {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.UsernameRepeat, errorMsg.UsernameRepeatMsg))
		logger.Error("service error:username has been used,please change")
		return
	}

	//查找电话是否唯一
	var countPhone int64
	err = database.DB.Model(&model.User{}).Where("phone = ?", requestBody.Phone).Count(&countPhone).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.MysqlError, errorMsg.MysqlErrorMsg))
		logger.Error("mysql error:" + err.Error())
		return
	}
	if countPhone >= 1 {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.PhoneRepeat, errorMsg.PhoneRepeatMsg))
		logger.Error("service error:phone has been used,please change")
		return
	}

	//判断电话号码格式是否正确
	matchString, err := regexp.
		MatchString("^1(3[0-9]|4[57]|5[0-35-9]|7[0135678]|8[0-9])\\d{8}$", requestBody.Phone)
	if err != nil || !matchString {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.PhoneError, errorMsg.PhoneErrorMsg))
		logger.Error("service error:phone is no right")
		return
	}

	//校验密码
	matched, err := regexp.MatchString("^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{8,20}$", requestBody.Password)
	if err != nil || !matched {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.PasswordError, errorMsg.PasswordErrorMsg))
		logger.Error("service error:phone is no right")
		return
	}

	//插入数据库
	user := model.User{
		Username: requestBody.Username,
		Password: requestBody.Password,
		Phone:    requestBody.Phone,
	}
	err = database.DB.Model(&model.User{}).Create(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.MysqlError, errorMsg.MysqlErrorMsg))
		logger.Error("mysql error:" + err.Error())
		return
	}

	c.JSON(http.StatusOK, result.Success("register success"))
}
