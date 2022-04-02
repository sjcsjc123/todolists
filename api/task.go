package api

import (
	"TodoLists/common/constant"
	"TodoLists/common/errorCode"
	"TodoLists/common/errorMsg"
	"TodoLists/common/model"
	"TodoLists/common/result"
	"TodoLists/database"
	"TodoLists/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"time"
)

var logger = utils.Logger

type TaskRequestBody struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func ListTask(c *gin.Context) {
	userId, _ := c.Get("userId")
	var tasks []model.Task
	err := database.DB.Model(&model.Task{}).Where("user_id = ?", userId).Order("start_time desc").Find(&tasks).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.MysqlError, errorMsg.MysqlErrorMsg))
		logger.Error("mysql error:" + err.Error())
		return
	}
	c.JSON(http.StatusOK, result.Success(tasks))
}

func CreateTask(c *gin.Context) {
	fromContext, _ := c.Get("userId")
	userId := fromContext.(int)
	var requestBody TaskRequestBody
	err := c.ShouldBindWith(&requestBody, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.ParseJsonError, errorMsg.ParseJsonErrorMsg))
		logger.Error("json error:" + err.Error())
		return
	}
	task := model.Task{
		UserId:            userId,
		Title:             requestBody.Title,
		Content:           requestBody.Content,
		StartTime:         time.Now(),
		StatusDescription: constant.NoDownMsg,
		Status:            constant.NoDown,
	}
	err = database.DB.Model(&model.Task{}).Where("user_id = ?", userId).Create(&task).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.MysqlError, errorMsg.MysqlErrorMsg))
		logger.Error("mysql error:" + err.Error())
		return
	}
	c.JSON(http.StatusOK, result.Success("create success"))
}

func FinishTask(c *gin.Context) {
	fromContext, _ := c.Get("userId")
	query := c.Query("taskId")
	userId := fromContext.(int)
	var task model.Task
	err := database.DB.Model(&model.Task{}).Where("user_id = ? and task_id = ?", userId, query).First(&task).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.MysqlError, errorMsg.MysqlErrorMsg))
		logger.Error("mysql error:" + err.Error())
		return
	}
	task.Status = constant.Finish
	task.StatusDescription = constant.FinishMsg
	err = database.DB.Model(&model.Task{}).Where("user_id = ? and task_id = ?", userId, query).Save(&task).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.MysqlError, errorMsg.MysqlErrorMsg))
		logger.Error("mysql error:" + err.Error())
		return
	}
	c.JSON(http.StatusOK, result.Success("finish success"))
}

func DeleteTask(c *gin.Context) {
	fromContext, _ := c.Get("userId")
	query := c.Query("taskId")
	userId := fromContext.(int)
	var task model.Task
	err := database.DB.Model(&model.Task{}).Where("user_id = ? and task_id = ?", userId, query).Delete(&task).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.MysqlError, errorMsg.MysqlErrorMsg))
		logger.Error("mysql error:" + err.Error())
		return
	}
	c.JSON(http.StatusOK, result.Success("delete success"))
}

func ListNoFinishTask(c *gin.Context) {
	fromContext, _ := c.Get("userId")
	userId := fromContext.(int)
	var tasks []model.Task
	err := database.DB.Model(&model.Task{}).Where("user_id = ? and status = ?", userId, 1).Find(&tasks).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.MysqlError, errorMsg.MysqlErrorMsg))
		logger.Error("mysql error:" + err.Error())
		return
	}
	c.JSON(http.StatusOK, result.Success(tasks))
}

func ListFinishTask(c *gin.Context) {
	fromContext, _ := c.Get("userId")
	userId := fromContext.(int)
	var tasks []model.Task
	err := database.DB.Model(&model.Task{}).Where("user_id = ? and status = ?", userId, 0).Find(&tasks).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, result.Error(errorCode.MysqlError, errorMsg.MysqlErrorMsg))
		logger.Error("mysql error:" + err.Error())
		return
	}
	c.JSON(http.StatusOK, result.Success(tasks))
}
