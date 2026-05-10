package main

import (
	"log"
	"media-server/config"
	cacheinfra "media-server/infrastructure/cache"
	"media-server/logger"
	"media-server/router"
)

var _config config.Config

func init() {
	var err error
	_config, err = config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	logger.InitLogger(_config)
	logger.Logger.Info("Logger initialized")
	config.ConnectDB()
	logger.Logger.Info("DB connection established")
	config.SyncDB()
	logger.Logger.Info("DB sync completed")
	defer logger.Logger.Sync()
}

func main() {
	_ = cacheinfra.NewInMemoryCache(50)
	_ = cacheinfra.NewInMemoryCache(50)

	router.InitRouter()
	log.Fatal(router.Start("0.0.0.0:" + _config.PORT))
}
