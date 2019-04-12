package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/alex-lee/go-between/internal/proxy"
)

var (
	front, back *net.UDPAddr
)

func init() {
	var (
		err               error
		frontStr, backStr string
	)
	flag.StringVar(&frontStr, "front", "127.0.0.1:8053", "listen address for UDP packets")
	flag.StringVar(&backStr, "back", "10.0.0.1:53", "remote address for UDP packets")
	flag.Parse()

	front, err = net.ResolveUDPAddr("udp", frontStr)
	if err != nil {
		log.Fatalf("Invalid frontend address: %s", frontStr)
	}
	back, err = net.ResolveUDPAddr("udp", backStr)
	if err != nil {
		log.Fatalf("Invalid backend address: %s", backStr)
	}
}

func main() {
	// Listen for shutdown signals.
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initialize and start the proxy.
	p := proxy.New(front, back)
	defer p.Stop()
	log.Printf("Listening on %s ...", front.String())

	// Wait for shutdown.
	<-sigChan
	log.Printf("Exit")
}
