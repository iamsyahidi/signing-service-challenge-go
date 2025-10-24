package main

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

const (
	ListenAddress = ":8080"
)

func main() {
	// Initialize dependencies
	repository := persistence.NewInMemoryDeviceRepository()
	deviceService := domain.NewDeviceService(repository)
	server := api.NewServer(ListenAddress, deviceService)

	log.Printf("Starting signature service on %s", ListenAddress)
	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
