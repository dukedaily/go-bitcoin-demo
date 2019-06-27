package main

import (
	"fmt"
	"os"
)

func main() {

	// var Args []string
	inputs := os.Args

	for i, v := range inputs {
		fmt.Printf("i: %d, %s\n", i, v)
	}
}
