package main

import (
	"github.com/guneyin/sbda-api-service/service"
	"log"
)

func main() {
	api, err := service.NewApiService()
	if err != nil {
		log.Fatal(err)
	}

	if err = api.Serve(); err != nil {
		log.Fatal(err)
	}
}
