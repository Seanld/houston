package houston


import (
	"net"
	"crypto/tls"
	"fmt"
	"strings"
	"path"
)


// Drop protocol portion of request URL string, and return.
func stripProtocol(request string) string {
	if strings.Contains(request, "://") {
		return strings.Split(request, "://")[1]
	}
	return request
}


// Strips CRLF terminator, and formats to string.
func requestAsString(request []byte, noProtocol bool) string {
	requestCopy := make([]byte, 1024)
	copy(requestCopy, request)

	for i, elem := range requestCopy {
		if elem == 10 || elem == 13 {
			requestCopy[i] = 0
		}
	}

	if noProtocol {
		return stripProtocol(string(requestCopy))
	} else {
		return string(requestCopy)
	}
}


type Server struct {
	TLSConfig   *tls.Config
	Router      Router
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
			dataStr := requestAsString(data, true)
			dataStrCleaned := path.Clean(dataStr)

			// Get and call the handler that matches the requested URL.
			handler := s.Router.GetRouteHandler(dataStrCleaned)
			handler(c)

			c.Close()
		}(conn)
	}
}
