package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lmxdawn/wallet/config"
	"github.com/lmxdawn/wallet/engine"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Start 启动服务
func Start(isSwag bool, configPath string) {

	conf, err := config.NewConfig(configPath)
	fmt.Println(conf)
	if err != nil || len(conf.Engines) == 0 {
		panic("Failed to load configuration")
	}

	var engines []*engine.ConCurrentEngine
	for _, engineConfig := range conf.Engines {
		eth, err := engine.NewEngine(engineConfig)
		if err != nil {
			panic(fmt.Sprintf("eth run err：%v", err))
		}
		engines = append(engines, eth)
	}

	for _, currentEngine := range engines {
		currentEngine.Run()
	}

	server := gin.Default()

	// 中间件
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	server.Use(SetDB(engines...))

	auth := server.Group("/api", AuthRequired())
	{
		auth.GET("/createWallet", CreateWallet)
		auth.GET("/delWallet", DelWallet)
		auth.GET("/withdraw", Withdraw)
		auth.GET("/getTransactionReceipt", GetTransactionReceipt)
	}

	if isSwag {
		swagHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
		server.GET("/swagger/*any", swagHandler)
	}

	err = server.Run(fmt.Sprintf(":%v", conf.App.Port))
	if err != nil {
		panic("start error")
	}

}
