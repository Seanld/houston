package houston


import (
	"net"
	"crypto/tls"
	"fmt"
)


type Server struct {
	SandboxDir string
	Config     *tls.Config
}


func NewServer(sandboxDirPath string, certificatePath string, keyPath string) Server {
	cer, _ := tls.LoadX509KeyPair(certificatePath, keyPath)
	return Server{
		SandboxDir: sandboxDirPath,
		Config: &tls.Config{Certificates: []tls.Certificate{cer}},
	}
}


// Start a configured server.
func (s Server) Start(args ...interface{}) {
	var hostName string
	if len(args) >= 1 {
		hostName = args[0].(string)
	} else {
		hostName = "localhost"
	}

	var port int
	if len(args) == 2 {
		port = args[1].(int)
	} else {
		port = 1965
	}

	ln, _ := tls.Listen("tcp", fmt.Sprintf("%s:%d", hostName, port), s.Config)

	for {
		conn, _ := ln.Accept()

		go func(c net.Conn) {
			data := make([]byte, 1024)
			c.Read(data)
			fmt.Println(data)
			c.Close()
		}(conn)
	}
}
