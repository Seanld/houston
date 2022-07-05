package houston


import (
	"crypto/tls"
	"log"
	"fmt"
	"strings"
	"os"
)


type Server struct {
	TLSConfig    *tls.Config
	Router       Router
	EnableLog    bool
	LogFilePath  string
	LogFile      *os.File
}


func NewServer(router Router, certificatePath string, keyPath string, args ...interface{}) Server {
	cer, err := tls.LoadX509KeyPair(certificatePath, keyPath)

	if err != nil {
		log.Fatalf("Error when loading key and certificate: %v", err)
	}

	enableLog := false
	if args != nil && len(args) >= 1 {
		enableLog = args[0].(bool)
	}

	logFilePath := "houston.log"
	if args != nil && len(args) == 2 {
		logFilePath = args[1].(string)
	}

	return Server{
		Router: router,
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cer}},
		EnableLog: enableLog,
		LogFilePath: logFilePath,
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

	var f *os.File
	var fileErr error

	if s.EnableLog {
		f, fileErr = os.OpenFile(s.LogFilePath, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
		defer f.Close()

		s.LogFile = f
	}

	if fileErr != nil {
		log.Println("Failed to open log file!")
	}

	ln, _ := tls.Listen("tcp", fmt.Sprintf("%s:%d", hostName, port), s.TLSConfig)

	for {
		conn, err := ln.Accept()

		if s.EnableLog {
			clientIp := strings.Split(conn.RemoteAddr().String(), ":")[0]
			message := fmt.Sprintf("New request from %s", clientIp)
			log.SetOutput(f)
			log.Output(1, message)
		}

		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(s, conn)
	}
}
