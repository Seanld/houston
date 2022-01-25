package houston


import (
	"net"
	"net/url"
	"crypto/tls"
	"fmt"
)


// Strips all NULL chars, and formats to string.
func requestAsString(request []byte) string {
	var endIndex int

	// Once we hit a NULL char, we know where to
	// cut the string off at, and where the CRLF
	// will be.
	for i, elem := range request {
		if elem == 0 {
			endIndex = i
			break
		}
	}

	newRequest := make([]byte, endIndex-2)

	for i:=0; i<endIndex-2; i++ {
		newRequest[i] = request[i]
	}

	return string(newRequest)
}


type Server struct {
	TLSConfig  *tls.Config
	Router     Router
}


func NewServer(router Router, certificatePath string, keyPath string) Server {
	cer, _ := tls.LoadX509KeyPair(certificatePath, keyPath)
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
		conn, _ := ln.Accept()

		go func(c net.Conn) {
			data := make([]byte, 1024)
			c.Read(data)
			dataStr := requestAsString(data)
			requestParsed, err := url.Parse(dataStr)

			if err != nil {
				fmt.Println("Error occurred when parsing URL!")
			}

			// Get and call the handler that matches the requested URL.
			handler := s.Router.GetRouteHandler(requestParsed.Path)
			handler(c)

			c.Close()
		}(conn)
	}
}
