package main

import (
	"log"

	"github.com/reche13/habitum/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg.Server.Port)
}