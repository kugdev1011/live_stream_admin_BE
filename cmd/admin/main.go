package main

import (
	"context"
	"fmt"
	"gitlab/live/be-live-api/cmd/admin/handler"
	"gitlab/live/be-live-api/conf"
	"gitlab/live/be-live-api/datasource"
	"gitlab/live/be-live-api/repository"
	"gitlab/live/be-live-api/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

// Validate method to perform validation using the validator library
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ds, err := datasource.NewDataSource()
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(ds.DB)

	appConfig := conf.GetApplicationConfig()

	srv := service.NewService(repo, ds.RClient, appConfig)

	e := echo.New()
	e.Server.MaxHeaderBytes = 10 << 20 //10MB

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v := validator.New()
	// Register custom validator with Echo
	e.Validator = &CustomValidator{validator: v}

	root := e.Group("/")

	handler := handler.NewHandler(root, srv)
	handler.Register()

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println("Server gracefully stopped")

}
