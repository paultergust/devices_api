package main

import (
    "context"
    "os"

    "dev.paultergust/devices-api/internal/db"
    "dev.paultergust/devices-api/internal/handler"
    "dev.paultergust/devices-api/internal/logger"
    "dev.paultergust/devices-api/internal/middleware"
    "dev.paultergust/devices-api/internal/repository"
    "dev.paultergust/devices-api/internal/service"

    "github.com/gin-gonic/gin"
)

func main() {
    log := logger.New()

    pool, err := db.NewPoolWithRetry(context.Background())
    if err != nil {
        log.Fatal().Err(err).Msg("db connection failed")
    }

    dbURL := os.Getenv("DATABASE_URL")

    if err := db.RunMigrations(dbURL); err != nil {
        log.Fatal().Err(err).Msg("migration failed")
    }

    repo := repository.NewDeviceRepository(pool)
    svc := service.NewDeviceService(repo)
    h := handler.NewHandler(svc)

    r := gin.New()
    r.Use(middleware.Logger(log))
    r.Use(gin.Recovery())

    r.POST("/devices", h.CreateDevice)
    r.GET("/devices", h.ListDevices)
    r.GET("/devices/:id", h.GetDevice)
    r.PATCH("/devices/:id", h.UpdateDevice)
    r.DELETE("/devices/:id", h.DeleteDevice)

    r.Run(":8080")
}
