package main

import (
	"RealTimeChat/config"
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
}
