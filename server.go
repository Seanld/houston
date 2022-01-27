package houston


import (
	"crypto/tls"
	"log"
	"fmt"
)


type Server struct {
	TLSConfig  *tls.Config
	Router     Router
}


func NewServer(router Router, certificatePath string, keyPath string) Server {
	cer, err := tls.LoadX509KeyPair(certificatePath, keyPath)

	if err != nil {
		log.Fatalf("Error when loading key and certificate: %v", err)
	}

	return Server{
		Router: router,
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cer}},
	}
}


// Start a configured server.
// Arguments:
// 1st?: Hostname
// 2nd?: Port
func (s *Server) Start(args ...interface{}) {
	// Use first argument for hostname, otherwise `localhost`.
	var hostName string
	if len(args) >= 1 {
		hostName = args[0].(string)
	} else {
		hostName = "localhost"
	}

	// Use second argument for port #, otherwise `1965`.
	var port int
	if len(args) == 2 {
		port = args[1].(int)
	} else {
		port = 1965
	}

	ln, _ := tls.Listen("tcp", fmt.Sprintf("%s:%d", hostName, port), s.TLSConfig)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go HandleConnection(s, conn)
	}
}
