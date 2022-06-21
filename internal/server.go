package internal

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Dann-Go/InnoTaxiUserService/internal/config"
	"github.com/Dann-Go/InnoTaxiUserService/internal/handler"
	"github.com/Dann-Go/InnoTaxiUserService/internal/migrations"
	"github.com/Dann-Go/InnoTaxiUserService/internal/repository"
	"github.com/Dann-Go/InnoTaxiUserService/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	server *http.Server
}

func initLogger() {
	logger := log.New()
	logger.Out = os.Stdout
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
}

func Inject() (*gin.Engine, error) {
	dbConfig := config.NewDbConfig()

	connection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.DBName, dbConfig.Password, dbConfig.SSLMode)
	log.Printf(connection)
	db, err := sqlx.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Printf("Start migrating database \n")
	err = migrations.MigrationUp(db)
	if err != nil {
		return nil, err
	}

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	authService := service.NewAuthorizationService(userRepository)
	handlers := handler.NewHandler(userService, authService)

	router := handlers.InitRoutes()

	router.Group("/").GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "alive"})
	})

	return router, nil
}

func (s *Server) Run(port string) error {
	initLogger()

	router, err := Inject()
	if err != nil {
		return err
	}
	s.server = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() {
	log.Println("Shutting down")

	s.server.Close()
}
