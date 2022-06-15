package internal

import (
	"fmt"
	"github.com/Dann-Go/InnoTaxiUserService/internal/handler"
	"github.com/Dann-Go/InnoTaxiUserService/internal/repository"
	"github.com/Dann-Go/InnoTaxiUserService/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

type Server struct {
	server *http.Server
}

type DbPostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func initLogger() {
	logger := log.New()
	logger.Out = os.Stdout
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
}

func envsCheck() {
	requiredEnvs := []string{"HOST", "DBPORT", "USERNAME", "PASSWORD",
		"DBNAME", "SSLMODE", "SERVPORT"}
	var msg []string
	for _, el := range requiredEnvs {
		val, exists := os.LookupEnv(el)
		if !exists || len(val) == 0 {
			msg = append(msg, el)
		}
	}
	if len(msg) > 0 {
		log.Fatal(strings.Join(msg, ", "), " env(s) not set")
	}
}

func Inject() *gin.Engine {
	cfg := DbPostgresConfig{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("DBPORT"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
		DBName:   os.Getenv("DBNAME"),
		SSLMode:  os.Getenv("SSLMODE"),
	}

	connection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)
	db, err := sqlx.Open("postgres", connection)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf(err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	router := handlers.InitRoutes()

	router.Group("/").GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "alive"})
		return
	})

	return router
}

func (s *Server) Run(port string) error {
	initLogger()
	envsCheck()

	router := Inject()
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
