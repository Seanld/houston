package houston


import (
	"crypto/tls"
	"log"
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
