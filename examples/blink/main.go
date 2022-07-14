package main

import (
	"fmt"
	"time"

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
		brd := mc.BoardForIndex(i)
		pci, err := brd.GetParamStr(mc.BoardPCIPositionParam)
		if err != nil {
			fmt.Println(err)
			return
		}

		bn, err := brd.GetParamStr(mc.BoardNameParam)
		if err != nil {
			fmt.Println(err)
			return
		}

		bi, err := brd.GetParamStr(mc.BoardIdentifierParam)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(pci, bn, bi)

		brd.SetParamStr(mc.OutputConfigParam+mc.LED, "SOFT")
		for i := 0; i < 5; i++ {
			brd.SetParamStr(mc.OutputStateParam+mc.LED, "ON")
			time.Sleep(1 * time.Second)
			brd.SetParamStr(mc.OutputStateParam+mc.LED, "OFF")
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Println("Done.")
}
