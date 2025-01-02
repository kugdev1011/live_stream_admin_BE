package main

import (
	"context"
	"fmt"
	"gitlab/live/be-live-api/cmd/admin/handler"
	"gitlab/live/be-live-api/conf"
	"gitlab/live/be-live-api/datasource"
	cmiddleware "gitlab/live/be-live-api/middleware"
	"gitlab/live/be-live-api/model"
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

	if err := ds.DB.AutoMigrate(&model.Role{}, &model.User{}, &model.AdminLog{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	repo := repository.NewRepository(ds.DB)
	streamServerConfig := conf.GetStreamServerConfig()
	streamServer := service.NewStreamServerService(streamServerConfig.HTTPURL, streamServerConfig.RTMPURL)
	//roleService := service.NewRoleService(repo, ds.RClient)
	srv := service.NewService(repo, ds.RClient, streamServer)
	conf.SeedRoles(srv.Role)
	conf.SeedSuperAdminUser(srv.User, srv.Role)

	log.Println("Seeding completed successfully")

	// conf.SeedRoles(srv.Role)
	// appConfig := conf.GetApplicationConfig()

	// srv := service.NewService(repo, ds.RClient)

	e := echo.New()
	e.Server.MaxHeaderBytes = 10 << 20 //10MB

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// it would be messed up if config change to other paths
	e.Use(cmiddleware.ExcludePathMiddleware("/api/file/recordings/", "/api/file/scheduled_videos/"))

	// Use CORS middleware, for local run
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     conf.GetApplicationConfig().AllowedOrigins,                                                       // Allow all origins (use specific origins for production)
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}, // Allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Access-Token"},                  // Allowed headers
		AllowCredentials: true,                                                                                             // Allow credentials like cookies
	}))

	v := validator.New()
	// Register custom validator with Echo
	e.Validator = &CustomValidator{validator: v}
	log.Println(conf.GetFileStorageConfig().RootFolder)

	root := e.Group("/")
	handler := handler.NewHandler(root, srv)

	fileH := e.Group("/api/file")
	fileH.GET("/avatar/:filename", func(c echo.Context) error {
		avatarPath := conf.GetFileStorageConfig().AvatarFolder + c.Param("filename")

		return c.File(avatarPath)
	})

	fileH.Use(handler.JWTMiddleware())
	fileH.Static("/", conf.GetFileStorageConfig().RootFolder)

	handler.Register()

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", conf.GetApplicationConfig().Port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
