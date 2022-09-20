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
	camfilePrimary   = flag.String("camfile-primary", "", "CAM file to use for primary board capture")
	camfileSecondary = flag.String("camfile-secondary", "", "CAM file to use for secondary board capture")
	primary          = flag.Int("primary", 0, "board number to use as primary for capture (default 0)")
	secondary        = flag.Int("secondary", 1, "board number to use as secondary for capture (default 1)")
)

func main() {
	flag.Parse()
	if *camfilePrimary == "" || *camfileSecondary == "" {
		fmt.Println("camfile-primary and camfile-secondary flags are both required in order to capture")
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

	bc, err := mc.GetParamInt(mc.ConfigurationHandle, mc.BoardCountParam)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Boards detected:", bc)

	// Create grabber for each board
	g2, err := createGrabber(*secondary, *camfileSecondary)
	if err != nil {
		fmt.Println(err)
		return
	}

	// primary last
	g1, err := createGrabber(*primary, *camfilePrimary)
	if err != nil {
		fmt.Println(err)
		return
	}
	g1.primary = true

	go g2.start()
	go g1.start()

	for {
		time.Sleep(time.Second)
	}
}

type grabber struct {
	board       int
	camfile     string
	x, y, pitch int
	ch          *mc.Channel
	primary     bool
}

func createGrabber(board int, camfile string) (*grabber, error) {
	g := grabber{board: board, camfile: camfile}
	brd := mc.BoardForIndex(board)

	var err error
	g.ch, err = brd.CreateChannel()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	g.setup()

	return &g, nil
}

func (g *grabber) setup() {
	// For all GrabLink boards but Grablink DualBase
	if err := g.ch.SetParamStr(mc.ConnectorParam, "M"); err != nil {
		g.Println("ConnectorParam", err)
		return
	}

	// Choose the CAM file
	if err := g.ch.SetParamStr(mc.CamFileParam, g.camfile); err != nil {
		g.Println("CamFileParam", err)
		return
	}

	// Set the color format.
	if err := g.ch.SetParamInt(mc.ColorFormatParam, mc.ColorFormatY8); err != nil {
		g.Println("ColorFormatParam", err)
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
		g.Println("SeqLengthFrParam", err)
		return
	}

	g.ch.RegisterCallback(g.cbhandler)
}

func (g *grabber) start() {
	fmt.Println("Starting grabber", g.board)

	// MC_SIG_SURFACE_PROCESSING: acquisition done and locked for processing
	if err := g.ch.SetParamInt(mc.SignalEnableParam+mc.SurfaceProcessingSignal, mc.SignalEnableOn); err != nil {
		g.Println("SurfaceProcessingSignal", err)
		return
	}

	// MC_SIG_ACQUISITION_FAILURE: acquisition failed.
	if err := g.ch.SetParamInt(mc.SignalEnableParam+mc.AcquisitionFailureSignal, mc.SignalEnableOn); err != nil {
		g.Println("AcquisitionFailureSignal", err)
		return
	}

	// Start Acquisitions for this channel.
	if err := g.ch.SetParamInt(mc.ChannelStateParam, int(mc.ChannelStateActive)); err != nil {
		g.Println("ChannelStateParam", err)
		return
	}

	defer g.stop()

	if g.primary {
		err := g.ch.SetParamStr(mc.ForceTrigParam, "TRIG")
		if err != nil {
			g.Println("ForceTrigParam", err)
		}
	}

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

func (g *grabber) cbhandler(info *mc.CallbackInfo) {
	switch mc.ParamID(info.Signal) {
	case mc.SurfaceProcessingSignal:
		s := mc.SurfaceForHandle(mc.Handle(info.SignalInfo))
		ptr, err := s.Ptr(g.x, g.y)
		if err != nil {
			fmt.Println(err)
			return
		}

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
	filename := fmt.Sprintf("%s_%d.jpg", filetime(time.Now()), g.board)
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

func (g *grabber) Println(name string, err error) {
	fmt.Println(g.board, name, err)
}
