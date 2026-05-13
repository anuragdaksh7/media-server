package main

import (
	"fileserver/internal/auth/handler"
	authMiddleware "fileserver/internal/auth/middleware"
	"fileserver/internal/auth/repository"
	"fileserver/internal/auth/service"
	"fileserver/internal/config"
	"fileserver/internal/db"
	fsHandler "fileserver/internal/filesystem/handler"
	fsRoutes "fileserver/internal/filesystem/routes"
	fsService "fileserver/internal/filesystem/service"
	jobHandler "fileserver/internal/jobs/handler"
	jobManager "fileserver/internal/jobs/manager"
	jobRoutes "fileserver/internal/jobs/routes"
	jobService "fileserver/internal/jobs/service"
	"fileserver/internal/middleware"
	realtimeHandler "fileserver/internal/realtime/handler"
	realtimeHub "fileserver/internal/realtime/hub"
	realtimeRoutes "fileserver/internal/realtime/routes"
	"fileserver/internal/storage/provider"
	streamHandler "fileserver/internal/streaming/handler"
	streamRoutes "fileserver/internal/streaming/routes"
	streamService "fileserver/internal/streaming/service"
	torrentHandler "fileserver/internal/torrent/handler"
	torrentManager "fileserver/internal/torrent/manager"
	torrentRoutes "fileserver/internal/torrent/routes"
	torrentService "fileserver/internal/torrent/service"
	uploadHandler "fileserver/internal/upload/handler"
	uploadRoutes "fileserver/internal/upload/routes"
	uploadService "fileserver/internal/upload/service"
	"fileserver/logger"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var conf config.Config

func init() {
	var err error
	conf, err = config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	logger.InitLogger(conf)

	storageProvider := provider.NewLocalProvider(conf.StoragePath)
	tManager := torrentManager.NewTorrentManager()

	hub := realtimeHub.NewHub()

	go hub.Run()

	realtimeHdl := realtimeHandler.NewRealtimeHandler(
		hub,
		tManager,
	)

	database := db.InitDB()

	authRepo := repository.NewAuthRepository(database)
	jobsManager := jobManager.NewJobManager()

	authService := service.NewAuthService(authRepo)
	filesystemService := fsService.NewFilesystemService(storageProvider)
	streamingService := streamService.NewStreamingService(
		storageProvider,
	)
	uploadSvc := uploadService.NewUploadService(
		storageProvider,
		conf.MaxUploadSize,
	)
	jobsService := jobService.NewJobService(
		jobsManager,
	)
	torrentSvc, err := torrentService.NewTorrentService(
		tManager,
		realtimeHdl,
	)
	if err != nil {
		panic(err)
	}

	authHandler := handler.NewAuthHandler(authService)
	filesystemHandler := fsHandler.NewFilesystemHandler(filesystemService)
	streamingHandler := streamHandler.NewStreamingHandler(
		streamingService,
	)
	uploadHdl := uploadHandler.NewUploadHandler(
		uploadSvc,
	)
	jobsHandler := jobHandler.NewJobHandler(
		jobsService,
	)
	torrentHdl := torrentHandler.NewTorrentHandler(
		torrentSvc,
	)

	r := gin.New()
	r.MaxMultipartMemory = 8 << 20 // 8 mb

	r.Use(
		middleware.RequestID(),
		middleware.Recovery(),
		middleware.RequestLogger(),
	)

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://localhost:5173",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the File Server API"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	api := r.Group("/api")

	protected := api.Group("")
	protected.Use(authMiddleware.JWTMiddleware())
	{
		protected.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"user_id": c.GetString("user_id"),
				"role":    c.GetString("role"),
			})
		})
	}

	fsRoutes.RegisterFilesystemRoutes(
		protected,
		filesystemHandler,
	)

	streamRoutes.RegisterStreamingRoutes(
		api,
		streamingHandler,
	)

	uploadRoutes.RegisterUploadRoutes(
		protected,
		uploadHdl,
	)

	realtimeRoutes.RegisterRealtimeRoutes(
		api,
		realtimeHdl,
	)

	jobRoutes.RegisterJobRoutes(
		api,
		jobsHandler,
	)

	torrentRoutes.RegisterTorrentRoutes(
		api,
		torrentHdl,
	)

	r.Run(":8080")
}
