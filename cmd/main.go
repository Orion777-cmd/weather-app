package main

import (
	"fmt"
	"github.com/Orion777-cmd/weather-app/initiator"
)

func main() {
    fmt.Println("Hello, world")
	initiator.Init()
	fmt.Println("Initiator completed successfully")
}