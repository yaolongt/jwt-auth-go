package main

import (
	"go-jwt/configs"
	"go-jwt/routes"
)

func init() {
	configs.InitEnv()
	configs.LoadKeys()
}

func main() {
	configs.GetDBInstance()
	routes.InitRoutes()
}
