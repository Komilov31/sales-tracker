package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Komilov31/sales-tracker/internal/config"
	"github.com/Komilov31/sales-tracker/internal/handler"
	"github.com/Komilov31/sales-tracker/internal/repository"
	"github.com/Komilov31/sales-tracker/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		zlog.Logger.Info().Msgf("recieved shutting signal %v. Shuting down", sig)
		cancel()
	}()

	zlog.Init()

	dbString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Cfg.Postgres.Host,
		config.Cfg.Postgres.Port,
		config.Cfg.Postgres.User,
		config.Cfg.Postgres.Password,
		config.Cfg.Postgres.Name,
	)
	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}

	db, err := dbpg.New(dbString, []string{}, opts)
	if err != nil {
		log.Fatal("could not init db: " + err.Error())
	}

	repository := repository.New(db)
	service := service.New(repository)
	handler := handler.New(ctx, service)

	router := ginext.New()
	registerRoutes(router, handler)

	zlog.Logger.Info().Msg("succesfully started server on " + config.Cfg.HttpServer.Address)
	return router.Run(config.Cfg.HttpServer.Address)
}

func registerRoutes(engine *ginext.Engine, handler *handler.Handler) {
	// Register static files
	engine.LoadHTMLFiles("static/index.html")
	engine.Static("/static", "static")

	// POST requests
	engine.POST("/items", handler.CreateItem)

	// GET requests
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.GET("/", handler.GetMainPage)
	engine.GET("/items", handler.GetAllItems)
	engine.GET("/analytics", handler.GetAggregated)
	engine.GET("/analytics/csv", handler.GetAggregatedCSV)
	engine.GET("/items/csv", handler.GetFilteredCSV)

	// PUT request
	engine.PUT("/items/:id", handler.UpdateItem)

	// DELETE request
	engine.DELETE("/items/:id", handler.DeleteItem)
}
