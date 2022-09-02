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

type pair struct {
	s1 *mc.Surface
	s2 *mc.Surface
}

var (
	camfilePrimary   = flag.String("camfile-primary", "", "CAM file to use for primary board capture")
	camfileSecondary = flag.String("camfile-secondary", "", "CAM file to use for secondary board capture")
	numberSurfaces   = flag.Int("number-surfaces", 10, "number of surfaces for each channel/board. defaults to 10")

	saveCh chan pair
)

const (
	fullWidth  = 14000
	fullHeight = 1000
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
	g2, err := createGrabber(1, *camfileSecondary)
	if err != nil {
		fmt.Println(err)
		return
	}

	// primary last
	g1, err := createGrabber(0, *camfilePrimary)
	if err != nil {
		fmt.Println(err)
		return
	}
	g1.primary = true

	saveCh = make(chan pair, *numberSurfaces)

	g1.createBuffers()
	defer g1.deleteBuffers()

	g1.createSurfaces()
	defer g1.deleteSurfaces()

	g2.createSurfaces()
	defer g2.deleteSurfaces()

	go func() {
		g2.start()
		g1.start()

		defer g1.stop()
		defer g2.stop()

		for {
			sig1info, err := g1.ch.WaitSignal(mc.AnySignal, 1000)
			if err != nil {
				g1.Println("WaitSignal", err)
				return
			}
			surface1 := g1.handleSignal(sig1info)

			sig2info, err := g2.ch.WaitSignal(mc.AnySignal, 1000)
			if err != nil {
				g2.Println("WaitSignal", err)
				return
			}
			surface2 := g2.handleSignal(sig2info)

			saveCh <- pair{s1: surface1, s2: surface2}
		}
	}()

	go func() {
		for p := range saveCh {
			ptr, err := p.s1.Ptr(fullWidth, fullHeight)
			if err != nil {
				fmt.Println(err)
				return
			}

			img := image.NewGray(image.Rect(0, 0, fullWidth, fullHeight))
			img.Pix = ptr

			saveImage(img)

			// free the surfaces when finished
			p.s1.SetParamInt(mc.SurfaceStateParam, mc.SurfaceStateFree)
			p.s2.SetParamInt(mc.SurfaceStateParam, mc.SurfaceStateFree)
		}
	}()

	for {
		time.Sleep(time.Second)
	}
}

func saveImage(img *image.Gray) {
	filename := fmt.Sprintf("%s.jpg", filetime(time.Now()))
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
