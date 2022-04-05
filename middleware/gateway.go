package middleware

import (
	"TodoLists/common/errorCode"
	"TodoLists/common/errorMsg"
	"TodoLists/common/result"
	"TodoLists/database"
	"TodoLists/utils"
	"context"
	_ "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

var logger = utils.Logger

func ValidToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		openid := c.GetHeader("openid")
		if openid != "" {
			c.Set("userId", openid)
			c.Next()
			return
		}
		//获取请求头中的token
		token := c.GetHeader("Authorization")
		logger.Info("parse token ...")
		//解析token
		userId, err := utils.GetUsernameFormToken(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, result.Error(errorCode.JwtToken, errorMsg.JwtTokenMsg))
			logger.Error("jwt token error:" + err.Error())
			c.Abort()
			return
		}
		logger.Info("valid token ...")
		//从redis中获取token比对
		redisToken, err := database.RedisClient.Get(context.Background(), "token"+userId).Result()
		if err != nil {
			c.JSON(http.StatusBadRequest, result.Error(errorCode.RedisError, errorMsg.RedisErrorMsg))
			logger.Error("redis error:" + err.Error())
			c.Abort()
			return
		}
		//比对token
		matched, err1 := regexp.MatchString(token, redisToken)
		if !matched || err1 != nil {
			c.JSON(http.StatusBadRequest, result.Error(errorCode.InvalidToken, errorMsg.InvalidTokenMsg))
			logger.Error("invalid token:" + err1.Error())
			c.Abort()
			return
		}
		c.Set("userId", userId)
		c.Next()
	}
}
