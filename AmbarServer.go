package main

import (
	"log"
	"youOrHell/src/config"
	"youOrHell/src/handler"
	redis "youOrHell/src/storage"
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
