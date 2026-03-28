package main

import "fmt"

type widget struct {
	name string
}

var globalCounter = 1

func rootErr() error {
	if globalCounter > 0 {
		switch globalCounter {
		case 1:
			globalCounter++
		default:
			globalCounter = 0
		}
	}

	return fmt.Errorf("root error")
}
