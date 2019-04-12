package proxy

import (
	"log"
	"net"
	"time"
)

// frontReader accepts packets on the frontend.
func (s *Server) frontReader(conn *net.UDPConn) {
	defer s.wg.Done()

	var err error
	buf := make([]byte, bufSize)

	for {
		select {
		case <-s.done:
			return
		default:
		}

		err = conn.SetReadDeadline(time.Now().Add(time.Second))
		if err != nil {
			log.Fatalf("Frontend: read error: failed to set deadline: %v", err)
		}

		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			if _, ok := err.(net.Error); ok {
				continue
			}
			log.Printf("Frontend: read error: %v", err)
			continue
		}
		if n == 0 {
			log.Printf("Frontend: read error: no bytes read")
			continue
		}

		data := make([]byte, n)
		copy(data, buf[0:n])
		s.frontRecv <- packet{addr, data}
	}
}
