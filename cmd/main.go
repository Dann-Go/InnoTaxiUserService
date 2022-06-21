package main

import (
	"github.com/Dann-Go/InnoTaxiUserService/internal/config"

	"github.com/Dann-Go/InnoTaxiUserService/internal"
	log "github.com/sirupsen/logrus"
)

// @title           InnoTaxi User Microservice
// @version         1.0
// @description     This is a user microservice for InnoTaxi App.

// @host      localhost:8000
// @BasePath  /

func main() {
	err := config.EnvsCheck()
	serverCfg := config.NewServerConfig()
	if err != nil {
		log.Fatalf("envs are not set %s", err.Error())
	}
	server := new(internal.Server)
	if err := server.Run(serverCfg.Port); err != nil {
		log.Fatalf("error while running server %s", err.Error())
	}
}
