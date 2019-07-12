package main

import "log"

func parseError(fmt string, err error) {
	if err != nil {
		log.Printf(fmt, err)
	}
}
