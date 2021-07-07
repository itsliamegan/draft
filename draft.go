package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var port uint

	flag.Usage = func() {
		fmt.Println("usage: draft [options] [<directory>]")
		flag.PrintDefaults()
	}

	flag.UintVar(&port, "port", 4000, "port to serve files from; port + 1 must also be available")
	flag.Parse()

	var root string
	var err error

	args := flag.Args()

	if len(args) > 1 {
		root = args[1]
	} else {
		root, err = os.Getwd()
		ExitIf(err)
	}

	changes := make(chan Change)

	go Serve(root, port)
	go Watch(root, changes)
	go Announce(changes, port+1)

	// Idle the process until manually cancelled.
	<-make(chan bool)
}
