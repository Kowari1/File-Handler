package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kowari1/File-Handler/internal/config"
	"github.com/Kowari1/File-Handler/internal/db"
	"github.com/Kowari1/File-Handler/internal/handler/device"
	"github.com/Kowari1/File-Handler/internal/parser"
	"github.com/Kowari1/File-Handler/internal/pdf"
	"github.com/Kowari1/File-Handler/internal/repository/postgres"
	"github.com/Kowari1/File-Handler/internal/service"
	"github.com/Kowari1/File-Handler/internal/service/file"
	"github.com/Kowari1/File-Handler/internal/worker"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("config load failed", zap.Error(err))
	}

	db, err := db.NewPostgresPool(
		ctx,
		cfg.DatabaseURL(),
		cfg.MaxDBConns,
		cfg.MinDBConns,
		cfg.MaxDBConnLifetime,
	)
	if err != nil {
		log.Fatal("db connect failed", zap.Error(err))
	}
	defer db.Close()

	pdfGenerator := pdf.NewPDFGenerator()
	deviceRepo := postgres.NewDeviceRepository(db)
	parseErrRepo := postgres.NewParseErrorRepository(db)
	processedFileRepo := postgres.NewProcessedFileRepository(db)
	deviceService := service.NewDeviceService(deviceRepo)
	deviceHandler := device.NewDeviceHandler(deviceService)

	tsvParser := parser.NewTSVParser()

	processor := file.NewProcessor(
		deviceRepo,
		parseErrRepo,
		tsvParser,
	)

	fileService := file.NewFileService(
		processor,
		processedFileRepo,
		*pdfGenerator,
		cfg.OutputDir,
	)

	jobs := make(chan string, 100)

	scanner := file.NewScanner(
		cfg.InputDir,
		cfg.ScanInterval,
		processedFileRepo,
		jobs,
		log,
	)

	pool := worker.New(
		4,
		fileService,
		jobs,
	)

	pool.Start(ctx)

	go scanner.Start(ctx)

	r := gin.Default()

	r.GET("/devices", deviceHandler.GetDevices)
	r.GET("/devices/:guid", deviceHandler.GetByUnitGUID)

	r.Static("/web", "./web")

	r.Run(":8080")

	<-ctx.Done()

	log.Info("shutting down...")

	close(jobs)

	pool.Wait()

	log.Info("application stopped")
}
