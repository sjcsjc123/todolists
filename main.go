package main

import (
	"TodoLists/common/config"
	"TodoLists/common/constant"
	"TodoLists/database"
	"TodoLists/routes"
	"TodoLists/utils"
)

func init() {
	utils.WriteLogger()
}

func main() {
	database.InitDatabase()
	r := routes.NewRoutes()
	_ = r.Run(config.Conf.GetString(constant.ServerPort))
}
