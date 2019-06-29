package main

import "fmt"

func main() {

	num := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

LABEL:
	for i, v := range num {
		if i == 5 {
			// goto LABEL
			continue LABEL
		}
		fmt.Printf("%d:%d\n", i, v)
	}

	fmt.Println("over")
}
