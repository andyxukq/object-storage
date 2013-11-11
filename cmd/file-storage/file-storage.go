package main

import (
	"flag"
	"fmt"

	"file-storage/handlers"
)

func main() {
	port := flag.String("port", "8001", "Listening port number")

	flag.Parse()

	handlers.Start(fmt.Sprint(":", *port))

	select {}
}
