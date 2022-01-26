package houston


import (
	"net"
	"net/url"
	"crypto/tls"
	"fmt"
	"path"
	"mime"
	"os"
	"log"
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


func GetMimetypeFromPath(targetPath string) string {
	extension := path.Ext(targetPath)
	if extension == ".gmi" || extension == ".gemini" {
		return "text/gemini"
	} else {
		return mime.TypeByExtension(extension)
	}
}


// If the path is a directory, append `index.gmi` to the end.
// Otherwise, keep the path.
// TODO Clean this up. Looks nasty.
func CompletePath(targetPath string) string {
	pathClean := path.Clean(targetPath)
	fileInfo, err := os.Stat(pathClean)

	if err == nil && fileInfo.IsDir() {
		return path.Join(pathClean, "index.gmi")
	} else {
		if path.Ext(pathClean) != "" {
			return pathClean
		} else {
			return pathClean + ".gmi"
		}
	}

	log.Fatalf("Error when opening %s: %v", pathClean, err)
	return ""
}


func isAllZeroes(bytes []byte) bool {
	for _, v := range bytes {
		if v != 0 {
			return false
		}
	}
	return true
}


func HandleConnection(s *Server, c net.Conn) {
	data := make([]byte, 1024)
	c.Read(data)

	if isAllZeroes(data) {
		return
	}

	dataStr := requestAsString(data)
	requestParsed, err := url.Parse(dataStr)

	if err != nil {
		fmt.Println("Error occurred when parsing URL!")
	}

	context := NewContext(dataStr, c)

	cleanedPath := path.Clean(requestParsed.Path)
	handledAsSandbox := false

	// First, see if there is a static file to serve
	// from a sandbox.
	for _, elem := range s.Router.Sandboxes {
		dir := path.Dir(cleanedPath)
		if dir == "." {
			dir = "/"
		}
		if elem.Path == dir {
			cleanedSandboxPath := path.Clean(elem.LocalPath)
			fullLocalPath := CompletePath(path.Join(cleanedSandboxPath, path.Base(cleanedPath)))
			mimeType := GetMimetypeFromPath(fullLocalPath)

			if (context.SendFile(mimeType, fullLocalPath) == nil) {
				handledAsSandbox = true
			}
		}
	}

	if !handledAsSandbox {
		// Get and call the handler that matches the requested URL.
		handler := s.Router.GetRouteHandler(requestParsed.Path)
		handler(context)
	}

	c.Close()
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
