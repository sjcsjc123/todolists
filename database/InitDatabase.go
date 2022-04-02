package database

import (
	"TodoLists/common/config"
	"TodoLists/common/constant"
	"TodoLists/common/model"
	"TodoLists/utils"
	"context"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
)

var (
	RedisClient *redis.Client
	DB          *gorm.DB
)

var logger = utils.Logger

func InitDatabase() {
	go func() {
		InitMysql()
	}()
	go func() {
		InitRedis()
	}()
}

func InitMysql() {
	username := config.Conf.GetString(constant.MysqlUsername)
	password := config.Conf.GetString(constant.MysqlPassword)
	addr := config.Conf.GetString(constant.MysqlAddr)
	dbname := config.Conf.GetString(constant.MysqlDbName)
	db, err := gorm.Open(mysql.Open(username+":"+password+"@tcp("+addr+")/"+dbname+"?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		Logger: logger2.Default.LogMode(logger2.Info),
	})
	if err != nil {
		logger.Error("connect mysql error:" + err.Error())
		panic("application start fail")
	}
	logger.Info("mysql connect success")
	DB = db
	err2 := db.AutoMigrate(&model.User{}, &model.Task{})
	if err2 != nil {
		panic(err2)
	}
	logger.Info("mysql init database success")
}

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Conf.GetString(constant.RedisAddr),
		Password: config.Conf.GetString(constant.RedisPassword),
		DB:       config.Conf.GetInt(constant.RedisDbName),
	})
	_, err := client.Ping(context.Background()).Result()
	logger.Info("redis connect success")
	if err != nil {
		logger.Error("connect Redis error: ", err.Error())
		panic("application start fail")
	}
	RedisClient = client
	logger.Info("redis init database success")
}
