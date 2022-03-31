package main

import (
	"fmt"

	mc "github.com/northvolt/go-multicam"
)

var (
	x, y, pitch int
	ch          *mc.Channel
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

	//  Create a channel.
	ch = mc.NewChannel()
	err := ch.Create()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ch.Delete()

	SetupCamera()

	// TODO: Register our Callback function for the MultiCam asynchronous signals.
	// status = McRegisterCallback(hChannel, McCallback, NULL);

	// TODO: Enable the signals we need:
	// MC_SIG_SURFACE_PROCESSING: acquisition done and locked for processing
	// MC_SIG_ACQUISITION_FAILURE: acquisition failed.
	// status = McSetParamInt(hChannel, MC_SignalEnable + MC_SIG_SURFACE_PROCESSING, MC_SignalEnable_ON);

	// status = McSetParamInt(hChannel, MC_SignalEnable + MC_SIG_ACQUISITION_FAILURE, MC_SignalEnable_ON);

	fmt.Println("Done.")
}

func SetupCamera() {
	// Link the channel to a board. Here we take the first board.
	if err := ch.SetParamInt(mc.DriverIndexParam, 1); err != nil {
		fmt.Println(err)
		return
	}

	// For all GrabLink boards but Grablink DualBase
	if err := ch.SetParamStr(mc.ConnectorParam, "M"); err != nil {
		fmt.Println(err)
		return
	}

	// Choose the CAM file
	if err := ch.SetParamStr(mc.CamFileParam, "KD6R309MX_L7296RG"); err != nil {
		fmt.Println(err)
		return
	}

	// Set the color format.
	if err := ch.SetParamInt(mc.ColorFormatParam, mc.ColorFormatY8); err != nil {
		fmt.Println(err)
		return
	}

	var err error

	// Retrieve channel size information.
	x, err = ch.GetParamInt(mc.ImageSizeXParam)
	if err != nil {
		fmt.Println(err)
		return
	}

	y, err = ch.GetParamInt(mc.ImageSizeYParam)
	if err != nil {
		fmt.Println(err)
		return
	}

	pitch, err = ch.GetParamInt(mc.BufferPitchParam)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("x:", x, "y:", y, "pitch:", pitch)

	// The number of images to acquire.
	if err := ch.SetParamInt(mc.SeqLengthFrParam, mc.IndeterminateLength); err != nil {
		fmt.Println(err)
		return
	}
}
