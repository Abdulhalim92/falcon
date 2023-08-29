package main

import (
	"falconapi/api/middlewares"
	"falconapi/api/routes"
	_ "falconapi/docs"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// @title   Сервис админки FalconApi
// @version  1.0
// @description FalconApi

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host   192.168.100.155:3006
// @BasePath  /api/v1
// @schemes http
func main() {
	initViper()

	router := gin.Default()

	middlewares.InitGinMiddlewares(router, routes.InitPublicRoutes, routes.InitProtectedRoutes)

	var listenIp = viper.GetString("ListenIP")
	var listenPort = viper.GetString("ListenPort")

	log.Printf("will listen on %v:%v", listenIp, listenPort)

	err := router.Run(fmt.Sprintf("%v:%v", listenIp, listenPort))
	log.Fatal(err)
}

func initViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("demo")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("unable to initialize viper: %w", err))
	}
	log.Println("viper config initialized")
}
