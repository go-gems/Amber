package main

import (
	"amber/src/config"
	"amber/src/handler"
	"amber/src/storage"
	"flag"
	"log"
)

func main() {
	inmemoryFlag := flag.Bool("inmemory", false, "Use in-memory storage instead of Redis")
	flag.Parse()

	configuration, err := config.FromFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	var service storage.Service
	if *inmemoryFlag {
		log.Println("Using in-memory storage")
		service = storage.NewInMemory()
	} else {
		log.Println("Using Redis storage")
		service, err = storage.New(configuration.Redis.Host, configuration.Redis.Port)
		if err != nil {
			log.Fatal(err)
		}
	}

	router := handler.New(configuration.Options.Schema, configuration.Options.Prefix, service, configuration.Admin.Password)

	log.Fatal(router.Run())
}
