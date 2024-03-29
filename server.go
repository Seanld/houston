package houston

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"golang.org/x/time/rate"
)

type Server struct {
	TLSConfig *tls.Config
	Router    *Router
	Config    ServerConfig
}

type ServerConfig struct {
	// TLS certificate and key file paths.
	CertificatePath string
	KeyPath         string

	// Self-explanatory.
	Hostname string
	Port     uint16

	// Whether connection logging should occur,
	// and where the file should be located.
	EnableLog   bool
	LogFilePath string
	// This can be accessed during the server's lifetime to
	// manually write more info to the config, other than what's
	// built-in to Houston.
	LogFile     *os.File

	// Whether the server should apply rate-limiting.
	EnableLimiting bool
	MaxRate        rate.Limit
	BucketSize     int

	// Example: file exists at `/page.gmi`, and user requests `/page`. Should
	// the server append the `.gmi` and return `/page.gmi`? This assumption will
	// happen before handler functions are checked.
	ImplyExtension bool
}

func NewServer(router *Router, config *ServerConfig) Server {
	if config.CertificatePath == "" || config.KeyPath == "" {
		log.Fatal("Must provide TLS certificate and key path to server config!")
	}

	cer, err := tls.LoadX509KeyPair(config.CertificatePath, config.KeyPath)

	if err != nil {
		log.Fatalf("Error when loading key and certificate: %v", err)
	}

	// Enforce hostname and port number defaults if not set.
	if len(config.Hostname) == 0 {
		config.Hostname = "localhost"
	}
	if config.Port == 0 {
		config.Port = 1965
	}

	// Rate-limiting defaults.
	if config.MaxRate == 0 {
		config.MaxRate = 2
	}
	if config.BucketSize == 0 {
		config.BucketSize = 2
	}

	if len(config.LogFilePath) == 0 {
		config.LogFilePath = "houston.log"
	}

	return Server{
		Router:    router,
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cer}},
		Config:    *config,
	}
}

func (s *Server) Start() {
	var f *os.File
	var fileErr error

	if s.Config.EnableLog {
		f, fileErr = os.OpenFile(s.Config.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		defer f.Close()

		s.Config.LogFile = f
		log.SetOutput(f)
	}

	if fileErr != nil {
		log.Println("Failed to open log file!")
	}

	ln, _ := tls.Listen("tcp", fmt.Sprintf("%s:%d", s.Config.Hostname, s.Config.Port), s.TLSConfig)

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(s, conn)
	}
}
