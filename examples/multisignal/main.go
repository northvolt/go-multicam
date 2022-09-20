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
	primary          = flag.Int("primary", 0, "board number to use as primary for capture (default 0)")
	secondary        = flag.Int("secondary", 1, "board number to use as secondary for capture (default 1)")
	numberSurfaces   = flag.Int("number-surfaces", 10, "number of surfaces for each channel/board. defaults to 10")
	height           = flag.Int("height", 1000, "frame height. defaults to 1000")
	width            = flag.Int("width", 7320, "width of single grabber. defaults to 7320")

	saveCh chan pair
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
	g1, err := createGrabber(*primary, *camfilePrimary)
	if err != nil {
		fmt.Println(err)
		return
	}
	g1.primary = true

	g2, err := createGrabber(*secondary, *camfileSecondary)
	if err != nil {
		fmt.Println(err)
		return
	}

	saveCh = make(chan pair, *numberSurfaces)

	g1.createBuffers()
	defer g1.deleteBuffers()

	g1.createSurfaces()
	defer g1.deleteSurfaces()

	g2.createSurfaces()
	defer g2.deleteSurfaces()

	go func() {
		g2.start()
		defer g2.stop()

		g1.start()
		defer g1.stop()

		defer func() {
			close(saveCh)
		}()

		for {
			sig1info, err := g1.ch.WaitSignal(mc.AnySignal, 1000)
			if err != nil {
				g1.Println("WaitSignal", err)
				return
			}
			surface1, err := g1.handleSignal(sig1info)
			if err != nil {
				fmt.Println("surface1 handleSignal", err)
				return
			}

			sig2info, err := g2.ch.WaitSignal(mc.AnySignal, 1000)
			if err != nil {
				g2.Println("WaitSignal", err)
				return
			}
			surface2, err := g2.handleSignal(sig2info)
			if surface2 == nil {
				fmt.Println("surface2 handleSignal", err)
				return
			}

			saveCh <- pair{s1: surface1, s2: surface2}
		}
	}()

	go func() {
		img := image.NewGray(image.Rect(0, 0, *width*2, *height))

		for p := range saveCh {
			ptr, err := p.s1.Ptr(*width*2, *height)
			if err != nil {
				fmt.Println("saveCh", err)
				return
			}

			copy(img.Pix, ptr)

			// free the surfaces while we save the JPG data for concurrency
			p.s1.SetParamInt(mc.SurfaceStateParam, mc.SurfaceStateFree)
			p.s2.SetParamInt(mc.SurfaceStateParam, mc.SurfaceStateFree)

			saveImage(img)
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
