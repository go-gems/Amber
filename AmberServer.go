package main

import (
	"amber/src/config"
	"amber/src/handler"
	redis "amber/src/storage"
	"log"
)

func main() {
	configuration, err := config.FromFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	service, err := redis.New(configuration.Redis.Host, configuration.Redis.Port)
	if err != nil {
		log.Fatal(err)
	}

	router := handler.New(configuration.Options.Schema, configuration.Options.Prefix, service)

	log.Fatal(router.Run())
}
