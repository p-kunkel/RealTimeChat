package main

import (
	"RealTimeChat/config"
	dict "RealTimeChat/dictionaries"
	"RealTimeChat/mappings"
	"RealTimeChat/models"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var err error

	autoMigrate := flag.Bool("automigrate", false, "")
	flag.Parse()

	if err = config.LoadEnvToCache(); err != nil {
		log.Fatalf("load env to chache failed, err: %s", err)
	}

	if err = config.ConnectToDB(); err != nil {
		log.Fatalf("failed connection to db, err: %s", err)
	}

	if *autoMigrate {
		if err = models.DBAutoMigrate(); err != nil {
			log.Fatalf("models migration error: %s", err)
		}
		fmt.Println("\nthe migration was successful")
		os.Exit(0)
	}

	go models.ListenDatabase("message_notify", config.DBAddres(), models.ListenMessageNotify)
	go models.ChatHub.Run()

	if err = dict.Dicts.LoadFromDB(); err != nil {
		log.Fatalf("loading dicts failed: %s", err)
	}

	if err = mappings.RunServer(); err != nil {
		log.Fatalf("the server failed to start: %s", err)
	}
}
