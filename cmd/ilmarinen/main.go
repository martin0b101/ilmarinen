package main

import (
	"fmt"

	"github.com/martin0b101/ilmarinen/config"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config, error: %v\n", err)
	}

	fmt.Printf("Ilmarinen is starting, on address: %s port: %v\n", config.Address, config.Port)
}
