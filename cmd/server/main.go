package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	gintrace "github.com/DataDog/dd-trace-go/contrib/gin-gonic/gin/v2"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/DataDog/dd-trace-go/v2/profiler"

	"github.com/yeison2020/subway-routing-service/internal/config"
	"github.com/yeison2020/subway-routing-service/internal/mbta"
	"github.com/yeison2020/subway-routing-service/internal/middlerware"
	"github.com/yeison2020/subway-routing-service/internal/routes"
)

func main() {

	// Load .env
	_ = godotenv.Load()

	// Initialize Zap logger
	middlerware.InitLogger() // returns *zap.Logger

	// Init MBTA local cache
	mbta.InitCache(300) // TTL 5 minutes

	cfg := config.LoadConfig()

	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.Use(gin.Recovery(), gin.Logger())
	server.Use(gintrace.Middleware(cfg.Service))

	tracer.Start(
		tracer.WithEnv(cfg.Env),
		tracer.WithService(cfg.Service),
		tracer.WithServiceVersion(cfg.Version),
	)
	
	err := profiler.Start(
		profiler.WithService(cfg.Service),
		profiler.WithEnv(cfg.Env),
		profiler.WithVersion(cfg.Version),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Datadog instrumentation defers
	defer profiler.Stop()
	defer tracer.Stop()

	// Plug middlerware
	server.Use(middlerware.LoggerMiddleware(middlerware.Logger))
	server.Use(middlerware.RequestIdMiddleware())
	server.Use(middlerware.PanicHandler())

	// Routes
	routes.RegisterRoutes(server)


	// App startup logs
	middlerware.Logger.Info("Application started")
	middlerware.Logger.Info("Server running",
		zap.String("port", cfg.ServerPort),
		zap.String("mode", "release"),
	)

	// Run server
	log.Fatal(server.Run(":" + cfg.ServerPort))

}
