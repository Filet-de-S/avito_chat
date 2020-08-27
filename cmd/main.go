package main

import (
	"avito-chat_service/internal/api/service"
	"log"
)

func main() {
	defer log.Println("exiting")

	s, err := service.New()
	if err != nil {
		log.Fatal("Service init failed: ", err)
	}

	err = s.Run()
	if err != nil {
		log.Fatal("Service failed: ", err)
	}
}
