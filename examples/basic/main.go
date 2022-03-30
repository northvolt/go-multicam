package main

import (
	"fmt"

	mc "github.com/northvolt/go-multicam"
)

func main() {
	err := mc.OpenDriver()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer mc.CloseDriver()

	fmt.Println("Driver was opened, now closing...")
}
