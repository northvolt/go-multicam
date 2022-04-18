package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"reflect"
	"strings"
	"time"
	"unsafe"

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

	// Create grabber for each board
	g1, err := createGrabber(1, "KD6R309MX_L7296RG_PRIMARY.cam")
	if err != nil {
		fmt.Println(err)
		return
	}

	g2, err := createGrabber(2, "KD6R309MX_L7296RG_SECONDARY.cam")
	if err != nil {
		fmt.Println(err)
		return
	}

	g3, err := createGrabber(3, "KD6R309MX_L7296SP_SECONDARY.cam")
	if err != nil {
		fmt.Println(err)
		return
	}

	go g1.start()
	go g2.start()
	go g3.start()

	for {
		time.Sleep(time.Second)
	}
}

type grabber struct {
	board       int
	camfile     string
	x, y, pitch int
	ch          *mc.Channel
}

func createGrabber(board int, camfile string) (*grabber, error) {
	g := grabber{board: board, camfile: camfile}
	g.ch = mc.NewChannel()
	err := g.ch.Create()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	g.setup()

	return &g, nil
}

func (g *grabber) setup() {
	if err := g.ch.SetParamInt(mc.DriverIndexParam, g.board); err != nil {
		fmt.Println(err)
		return
	}

	// For all GrabLink boards but Grablink DualBase
	if err := g.ch.SetParamStr(mc.ConnectorParam, "M"); err != nil {
		fmt.Println(err)
		return
	}

	// Choose the CAM file
	if err := g.ch.SetParamStr(mc.CamFileParam, g.camfile); err != nil {
		fmt.Println(err)
		return
	}

	// Set the color format.
	if err := g.ch.SetParamInt(mc.ColorFormatParam, mc.ColorFormatY8); err != nil {
		fmt.Println(err)
		return
	}

	var err error

	// Retrieve channel size information.
	g.x, err = g.ch.GetParamInt(mc.ImageSizeXParam)
	if err != nil {
		fmt.Println(err)
		return
	}

	g.y, err = g.ch.GetParamInt(mc.ImageSizeYParam)
	if err != nil {
		fmt.Println(err)
		return
	}

	g.pitch, err = g.ch.GetParamInt(mc.BufferPitchParam)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("x:", g.x, "y:", g.y, "pitch:", g.pitch)

	// The number of images to acquire.
	if err := g.ch.SetParamInt(mc.SeqLengthFrParam, mc.IndeterminateLength); err != nil {
		fmt.Println(err)
		return
	}

	g.ch.RegisterCallback(g.cbhandler)
}

func (g *grabber) start() {
	// MC_SIG_SURFACE_PROCESSING: acquisition done and locked for processing
	if err := g.ch.SetParamInt(mc.SignalEnableParam+mc.SurfaceProcessingSignal, mc.SignalEnableOn); err != nil {
		fmt.Println(err)
		return
	}

	// MC_SIG_ACQUISITION_FAILURE: acquisition failed.
	if err := g.ch.SetParamInt(mc.SignalEnableParam+mc.AcquisitionFailureSignal, mc.SignalEnableOn); err != nil {
		fmt.Println(err)
		return
	}

	// Start Acquisitions for this channel.
	if err := g.ch.SetParamInt(mc.ChannelStateParam, int(mc.ChannelStateActive)); err != nil {
		fmt.Println(err)
		return
	}

	defer g.stop()

	for {
		time.Sleep(time.Second)
	}
}

func (g *grabber) stop() {
	if err := g.ch.SetParamInt(mc.ChannelStateParam, int(mc.ChannelStateIdle)); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Done.")

	g.ch.Delete()
}

func (g *grabber) cbhandler(info *mc.SignalInfo) {
	switch mc.ParamID(info.Signal) {
	case mc.SurfaceProcessingSignal:
		s := mc.SurfaceForHandle(mc.Handle(info.SignalInfo))
		pimg, err := s.GetParamPtr(mc.SurfaceAddrParam)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("frame received at address", pimg)
		h := &reflect.SliceHeader{
			Data: uintptr(pimg),
			Len:  int(g.x * g.y),
			Cap:  int(g.x * g.y),
		}
		ptr := *(*[]byte)(unsafe.Pointer(h))

		img := image.NewGray(image.Rect(0, 0, g.x, g.y))
		img.Pix = ptr

		g.saveImage(img)
	case mc.AcquisitionFailureSignal:
		fmt.Println("frame error")
	default:
		fmt.Println("other error")
	}
}

func (g *grabber) saveImage(img *image.Gray) {
	filename := fmt.Sprintf("%s_%d.jpg", filetime(time.Now()))
	f, err := os.Create(filename)
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

func filetime(t time.Time) string {
	id := t.UTC().Format(time.RFC3339Nano)
	id = strings.ReplaceAll(id, ":", "-")
	return id
}
