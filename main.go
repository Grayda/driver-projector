package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {

	_, err := NewProjectorDriver()

	if err != nil {
		log.Fatalf("Failed to create Projector driver: %s", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	s := <-c
	fmt.Println("Got signal:", s)

}
