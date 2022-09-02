package main

// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"

	mc "github.com/northvolt/go-multicam"
)

var (
	bufferSize  int
	bufferPitch int
	buffers     []unsafe.Pointer
)

type grabber struct {
	board       int
	camfile     string
	x, y, pitch int
	ch          *mc.Channel
	surfaces    []*mc.Surface
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
}

func (g *grabber) createSurfaces() error {
	if err := g.ch.SetParamInt(mc.SurfaceCountParam, *numberSurfaces); err != nil {
		return err
	}

	for i := 0; i < *numberSurfaces; i++ {
		s := mc.NewSurface()
		g.surfaces = append(g.surfaces, s)

		bufferOffset, err := g.ch.GetParamInt(mc.MinBufferPitchParam)
		if err != nil {
			g.Println("MinBufferPitchParam", err)
			return err
		}

		// bufferOffset *= module->moduleIndex;
		bufferOffset *= g.board

		// bufferAddress = m_pImageBuffers[i] + bufferOffset;
		bufferAddress := uintptr(buffers[i]) + uintptr(bufferOffset)

		// McSetParamPtr(surface, MC_SurfaceAddr, bufferAddress);
		s.SetParamPtr(mc.SurfaceAddrParam, unsafe.Pointer(bufferAddress))

		// McSetParamInt(surface, MC_SurfacePitch, module->bufferPitch);
		s.SetParamInt(mc.SurfacePitchParam, bufferPitch)

		// McSetParamInt(surface, MC_SurfaceSize, module->bufferSize);
		s.SetParamInt(mc.SurfaceSizeParam, bufferSize)

		// McSetParamInst(module->channel, MC_Cluster + i, module->surfaces[i]);
		g.ch.SetParamInst(mc.ClusterParam+mc.ParamID(i), s.Handle())
	}

	return nil
}

func (g *grabber) deleteSurfaces() error {
	for _, v := range g.surfaces {
		v.Delete()
	}

	return nil
}

// Allocate user buffers big enough to contain pixel data from all modules
func (g *grabber) createBuffers() (err error) {
	bufferSize, err = g.ch.GetParamInt(mc.BufferSizeParam)
	if err != nil {
		return
	}

	bufferPitch, err = g.ch.GetParamInt(mc.BufferPitchParam)
	if err != nil {
		return
	}

	totalBufferSize := bufferSize + *numberSurfaces*bufferPitch

	for i := 0; i < *numberSurfaces; i++ {
		buffers = append(buffers, C.malloc(C.ulong(totalBufferSize)))
	}

	return nil
}

func (g *grabber) deleteBuffers() error {
	for i := 0; i < *numberSurfaces; i++ {
		C.free(buffers[i])
	}

	return nil
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

	if g.primary {
		err := g.ch.SetParamStr(mc.ForceTrigParam, "TRIG")
		if err != nil {
			g.Println("ForceTrigParam", err)
		}
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

func (g *grabber) handleSignal(info *mc.SignalInfo) *mc.Surface {
	switch mc.ParamID(info.Signal()) {
	case mc.SurfaceProcessingSignal:
		s := mc.SurfaceForHandle((info.SignalInfo()))

		// set surface to reserved so it is not overwritten
		s.SetParamInt(mc.SurfaceStateParam, mc.SurfaceStateReserved)

		return s
	case mc.AcquisitionFailureSignal:
		fmt.Println("frame error")
	default:
		fmt.Println("other error")
	}
	return nil
}

func (g *grabber) Println(name string, err error) {
	fmt.Println(g.board, name, err)
}
