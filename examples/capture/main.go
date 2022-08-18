package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"strings"
	"time"

	mc "github.com/northvolt/go-multicam"
)

var (
	camfile = flag.String("camfile", "", "CAM file to use for capture")

	x, y, pitch int
	ch          *mc.Channel
)

func main() {
	flag.Parse()
	if *camfile == "" {
		fmt.Println("camfile flag is required in order to capture")
		return
	}

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

	// Get board
	brd := mc.BoardForIndex(1)

	//  Create a channel for board.
	var err error
	ch, err = brd.CreateChannel()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ch.Delete()

	SetupCamera()

	ch.RegisterCallback(cbhandler)

	// MC_SIG_SURFACE_PROCESSING: acquisition done and locked for processing
	if err := ch.SetParamInt(mc.SignalEnableParam+mc.SurfaceProcessingSignal, mc.SignalEnableOn); err != nil {
		fmt.Println(err)
		return
	}

	// MC_SIG_ACQUISITION_FAILURE: acquisition failed.
	if err := ch.SetParamInt(mc.SignalEnableParam+mc.AcquisitionFailureSignal, mc.SignalEnableOn); err != nil {
		fmt.Println(err)
		return
	}

	// Start Acquisitions for this channel.
	if err := ch.SetParamInt(mc.ChannelStateParam, int(mc.ChannelStateActive)); err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := ch.SetParamInt(mc.ChannelStateParam, int(mc.ChannelStateIdle)); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Done.")
	}()

	for {
		time.Sleep(time.Second)
	}
}

func SetupCamera() {
	// For all GrabLink boards but Grablink DualBase
	if err := ch.SetParamStr(mc.ConnectorParam, "M"); err != nil {
		fmt.Println(err)
		return
	}

	// Choose the CAM file
	if err := ch.SetParamStr(mc.CamFileParam, *camfile); err != nil {
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

func cbhandler(info *mc.SignalInfo) {
	switch mc.ParamID(info.Signal) {
	case mc.SurfaceProcessingSignal:
		s := mc.SurfaceForHandle(mc.Handle(info.SignalInfo))
		ptr, err := s.Ptr(x, y)
		if err != nil {
			fmt.Println(err)
			return
		}

		img := image.NewGray(image.Rect(0, 0, x, y))
		img.Pix = ptr

		saveImage(img)
	case mc.AcquisitionFailureSignal:
		fmt.Println("frame error")
	default:
		fmt.Println("other error")
	}
}

func saveImage(img *image.Gray) {
	f, err := os.Create(filename(time.Now()) + ".jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	opt := jpeg.Options{
		Quality: 90,
	}
	err = jpeg.Encode(f, img, &opt)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func filename(t time.Time) string {
	id := t.UTC().Format(time.RFC3339Nano)
	id = strings.ReplaceAll(id, ":", "-")
	return id
}
