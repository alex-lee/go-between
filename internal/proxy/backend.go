package proxy

import (
	"log"
	"net"
	"time"
)

// backReader accepts packets on the backend, for a specific session.
func (s *Server) backReader(origAddr *net.UDPAddr, c *net.UDPConn) {
	defer s.wg.Done()

	var err error
	buf := make([]byte, bufSize)

	for {
		select {
		case <-s.done:
			return
		default:
		}

		err = c.SetReadDeadline(time.Now().Add(time.Second))
		if err != nil {
			log.Fatalf("Read error: failed to set deadline: %v", err)
		}

		n, _, err := c.ReadFromUDP(buf)
		if err != nil {
			if _, ok := err.(net.Error); ok {
				continue
			}
			log.Printf("Read error: %v", err)
			continue
		}
		if n == 0 {
			log.Printf("Read error: no bytes read")
			continue
		}

		data := make([]byte, n)
		copy(data, buf[0:n])
		s.backRecv <- packet{origAddr, data}
	}
}
