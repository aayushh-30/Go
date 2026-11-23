// Package main contains small operating-system signal helpers for the stage0 server.
package main

// Import os so the helper can use the signal API's signal type.
import "os"

// Import os/signal so the process can receive Ctrl+C.
import "os/signal"

// signalNotify is a variable so the main code remains easy to replace in tests.
var signalNotify = func(channel chan<- os.Signal) {
	// Register interrupt signals with the channel.
	signal.Notify(channel, os.Interrupt)
}
