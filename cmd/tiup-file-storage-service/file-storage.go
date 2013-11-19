package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"file-storage/handlers"
)

func main() {
	log.Println("Current GOMAXPROCS:", runtime.GOMAXPROCS(0))
	port := flag.String("port", "8001", "Listening port number")

	flag.Parse()

	handlers.Start(fmt.Sprint(":", *port))

	select {}
}
