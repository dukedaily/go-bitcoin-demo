package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func main() {
	err := test2()
	if err != nil {
		fmt.Printf("err +v :%+v\n", err)
		fmt.Printf("err :%v\n", err)
	}
}

func test2() error {
	err := test1()
	if err != nil {
		// return errors.WithMessage(err, "test2")
		return errors.WithStack(err)
		// return errors.Wrap(err, "test2")
		return err

	}
	return err
}

func test1() error {
	_, err := os.Open("hello.txt")
	if err != nil {
		return errors.WithMessage(err, "test1")
		// return err
	}
	return nil
}
