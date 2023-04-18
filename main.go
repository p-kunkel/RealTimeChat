package main

import (
	"RealTimeChat/config"
	"RealTimeChat/mappings"
	"RealTimeChat/models"
	"log"
)

func main() {
	var err error

	if err = config.LoadEnvToCache(); err != nil {
		log.Fatalf("load env to chache failed, err: %s", err)
	}

	if err = config.ConnectToDB(); err != nil {
		log.Fatalf("failed connection to db, err: %s", err)
	}

	if err = models.DBAutoMigrate(); err != nil {
		log.Fatalf("models migration error: %s", err)
	}

	if err = mappings.RunServer(); err != nil {
		log.Fatalf("the server failed to start: %s", err)
	}
}
