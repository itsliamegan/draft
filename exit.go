package main

import "log"

func ExitIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
