package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func main() {
	err := test2()
	if err != nil {
		// fmt.Println("err:", err)
		fmt.Printf("err :%+v\n", err)
	}

	fmt.Println("Hello world")
}

func test2() error {
	err := test1()
	if err != nil {
		return errors.WithMessage(err, "test2")
		// return errors.Wrap(err, "test2")
		// return errors.New("hello")
	
	return nil
}

func test1() error {
	_, err := os.Open("helloworld.txt")
	if err != nil {
		return errors.Wrap(err, "test1")
	}

	return nil
}
