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

	b1 := byte(0x00)
	b2 := byte(0)

	fmt.Printf("b1:%b\n", b1)
	fmt.Printf("b2:%b\n", b2)

}
