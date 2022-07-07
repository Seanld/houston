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
	Router       *Router
	Config       ServerConfig
}


type ServerConfig struct {
	// TLS certificate and key file paths.
	CertificatePath  string
	KeyPath          string

	// Self-explanatory.
	Hostname         string
	Port             uint16

	// Whether connection logging should occur,
	// and where the file should be located.
	EnableLog        bool
	LogFilePath      string
	// This can be accessed during the server's lifetime to
	// manually write more info to the config, other than what's
	// built-in to Houston.
	LogFile          *os.File

	// Whether the server should apply rate-limiting.
	EnableLimiting   bool
}


func NewServer(router *Router, config *ServerConfig) Server {
	if config.CertificatePath == "" || config.KeyPath == "" {
		log.Fatal("Must provide TLS certificate and key path to server config!")
	}

	cer, err := tls.LoadX509KeyPair(config.CertificatePath, config.KeyPath)

	if err != nil {
		log.Fatalf("Error when loading key and certificate: %v", err)
	}

	return Server{
		Router: router,
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cer}},
		Config: *config,
	}
}


func (s *Server) Start() {
	// If IP or hostname is not given, set to `localhost`.
	if len(s.Config.Hostname) == 0 {
		s.Config.Hostname = "localhost"
	}

	// If port number is not given (is 0), default to 1965 (standard Gemini port).
	if s.Config.Port == 0 {
		s.Config.Port = 1965
	}

	var f *os.File
	var fileErr error

	if s.Config.EnableLog {
		f, fileErr = os.OpenFile(s.Config.LogFilePath, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
		defer f.Close()

		s.Config.LogFile = f
	}

	if fileErr != nil {
		log.Println("Failed to open log file!")
	}

	ln, _ := tls.Listen("tcp", fmt.Sprintf("%s:%d", s.Config.Hostname, s.Config.Port), s.TLSConfig)

	for {
		conn, err := ln.Accept()

		if s.Config.EnableLog {
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
