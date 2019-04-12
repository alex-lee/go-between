package proxy

import (
	"log"
	"net"
	"sync"
)

const bufSize = 4096

// Server is a UDP proxy.
type Server struct {
	front, back *net.UDPAddr
	sessions    map[string]*session // map originator address to a session

	frontRecv chan packet
	backRecv  chan packet

	wg   sync.WaitGroup
	done chan struct{}
}

type session struct {
	origAddr *net.UDPAddr // address of originator
	conn     *net.UDPConn // connection on the backend
}

type packet struct {
	origAddr *net.UDPAddr
	data     []byte
}

// New returns a new proxy server.
func New(front, back *net.UDPAddr) *Server {
	s := &Server{
		front:     front,
		back:      back,
		sessions:  make(map[string]*session),
		frontRecv: make(chan packet),
		backRecv:  make(chan packet),
		done:      make(chan struct{}),
	}

	s.wg.Add(1)
	go s.run()

	return s
}

// Stop shuts down the proxy.
func (s *Server) Stop() {
	close(s.done)
	s.wg.Wait()
}

// run handles the main proxy logic.
// The frontend accepts packets and forwards them to the backend.
// Sessions are created to track return traffic.
func (s *Server) run() {
	defer s.wg.Done()

	var err error

	conn, err := net.ListenUDP("udp", s.front)
	if err != nil {
		log.Fatalf("Listen failed: %v", err)
	}

	s.wg.Add(1)
	go s.frontReader(conn)

	var p packet
	for {
		select {
		case <-s.done:
			return

		case p = <-s.frontRecv:
			s.sendToBackend(p)

		case p = <-s.backRecv:
			s.sendToOriginator(conn, p)
		}
	}
}

func (s *Server) sendToBackend(p packet) {
	sess, err := s.getSession(p.origAddr)
	if err != nil {
		log.Printf("Session access error: %v", err)
		return
	}
	n, err := sess.conn.Write(p.data)
	if err != nil {
		log.Printf("Write error: %v", err)
		return
	}
	if n != len(p.data) {
		log.Printf("Write error: incomplete: %d of %d bytes", n, len(p.data))
		return
	}
}

func (s *Server) sendToOriginator(conn *net.UDPConn, p packet) {
	n, err := conn.WriteToUDP(p.data, p.origAddr)
	if err != nil {
		log.Printf("Write error: %v", err)
		return
	}
	if n != len(p.data) {
		log.Printf("Write error: incomplete: %d of %d bytes", n, len(p.data))
		return
	}
}

// getSession gets the session for the given address.
// If there is no existing session, a new one is created.
func (s *Server) getSession(origAddr *net.UDPAddr) (*session, error) {
	addr := origAddr.String()
	sess, ok := s.sessions[addr]
	if ok {
		return sess, nil
	}

	conn, err := net.DialUDP("udp", nil, s.back)
	if err != nil {
		return nil, err
	}

	sess = &session{origAddr, conn}
	s.sessions[addr] = sess
	log.Printf("New session created for address %s", addr)

	s.wg.Add(1)
	go s.backReader(origAddr, conn)

	return sess, nil
}
