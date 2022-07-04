package main

import (
	"fmt"

	mc "github.com/northvolt/go-multicam"
)

func main() {
	if err := mc.OpenDriver(); err != nil {
		fmt.Println(err)
		return
	}
	defer mc.CloseDriver()

	if err := mc.SetParamStr(mc.ConfigurationHandle, mc.ErrorLogParam, "error.log"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Driver was opened...")
	bc, err := mc.GetParamInt(mc.ConfigurationHandle, mc.BoardCountParam)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Boards detected:", bc)
	for i := 0; i < bc; i++ {
		bn, err := mc.GetParamStr(mc.DefaultBoardHandle+mc.Handle(i), mc.BoardNameParam)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(bn)
	}

	fmt.Println("Done.")
}
