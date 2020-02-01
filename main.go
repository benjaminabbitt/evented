package main

import (
	"fmt"
	"github.com/benjaminabbitt/evented/framework"
	memoryRepository "github.com/benjaminabbitt/evented/repository/memory"
	console "github.com/benjaminabbitt/evented/transport/console"
)

func main() {
	framework.NewServer(memoryRepository.NewMemoryRepository(), console.NewConsoleSender(), 8080, "localhost:8081")
	fmt.Println("Test")
}
