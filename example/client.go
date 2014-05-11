package main

import pulseaudio ".."
import (
	"github.com/guelfey/go.dbus" // imported as dbus.

	"log"
	"strconv"
)

//
//-------------------------------------------------------------------[ CLIENT ]--

// AppPulse is a client that connects 6 callbacks.
//
type AppPulse struct{}

func (ap *AppPulse) NewSink(path dbus.ObjectPath) {
	log.Println("one: NewSink", path)
}

func (ap *AppPulse) SinkRemoved(path dbus.ObjectPath) {
	log.Println("one: SinkRemoved", path)
}

func (ap *AppPulse) NewPlaybackStream(path dbus.ObjectPath) {
	log.Println("one: NewPlaybackStream", path)
}

func (ap *AppPulse) PlaybackStreamRemoved(path dbus.ObjectPath) {
	log.Println("one: PlaybackStreamRemoved", path)
}

func (ap *AppPulse) DeviceVolumeUpdated(path dbus.ObjectPath, values []uint32) {
	log.Println("one: device volume", path, values)
}

func (ap *AppPulse) StreamVolumeUpdated(path dbus.ObjectPath, values []uint32) {
	log.Println("one: stream volume", path, values)
}

// ClientTwo is a client that connects only one callback.
//
type ClientTwo struct {
	*pulseaudio.Client
}

func (two *ClientTwo) DeviceVolumeUpdated(path dbus.ObjectPath, values []uint32) {
	log.Println("two: volume updated", path)
}

// Show is an example to show how to get properties.
func (two *ClientTwo) Show() {
	// Get the list of streams from the Core and show some informations about them.
	// You better handle errors that were not checked here for code clarity.

	// Get the list of playback streams from the core.
	streams, _ := two.Core().ListPath("PlaybackStreams") // []ObjectPath
	for _, stream := range streams {

		// Get the device to query properties for the stream referenced by his path.
		dev := two.Device(stream)

		// Get some informations about this stream.
		mute, _ := dev.Bool("Mute")         // bool
		vols, _ := dev.ListUint32("Volume") // []uint32
		println("stream", volumeText(mute, vols))
	}

	// Same with sinks.
	sinks, _ := two.Core().ListPath("Sinks")
	for _, sink := range sinks {
		dev := two.Device(sink)
		name, _ := dev.String("Name") // string
		mute, _ := dev.Bool("Mute")
		vols, _ := dev.ListUint32("Volume")
		println("sink  ", volumeText(mute, vols), name)
	}
}

func volumeText(mute bool, vals []uint32) string {
	if mute {
		return "muted"
	}
	vol := int(volumeAverage(vals)) * 100 / 65535
	return " " + strconv.Itoa(vol) + "% "
}

func volumeAverage(vals []uint32) uint32 {
	var vol uint32
	if len(vals) > 0 {
		for _, cur := range vals {
			vol += cur
		}
		vol /= uint32(len(vals))
	}
	return vol
}

// Create a pulse dbus service with 2 clients.
func main() {
	pulse, e := pulseaudio.New()
	if e != nil {
		log.Panicln("connect", e)
	}

	app := &AppPulse{}
	pulse.Register(app)

	two := &ClientTwo{pulse}
	pulse.Register(two)

	two.Show()

	// Mute all playback streams.
	streams, _ := two.Core().ListPath("PlaybackStreams")
	for _, stream := range streams {
		dev := two.Stream(stream)
		e = dev.Set("Mute", true)
		if e != nil {
			log.Println(e)
		}
	}

	pulse.Listen()
}

//

//

// introspect

// import "github.com/guelfey/go.dbus/introspect"

// s, ei := introspect.Call(pulse.core)
// log.Err(ei, "intro")
// for _, interf := range s.Interfaces {
// 	log.Println(interf.Name)
// 	for _, sig := range interf.Methods {
// 		log.Println("  method", sig)
// 	}

// 	log.Println(interf.Name)
// 	for _, sig := range interf.Signals {
// 		log.Println("  signal", sig)
// 	}
// }
