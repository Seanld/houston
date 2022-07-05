package houston


import (
	"net"
	"net/url"
	"fmt"
	"path"
	"path/filepath"
	"mime"
	"os"
	"strings"
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


// Get the common part of two paths.
func getSharedPath(path1 string, path2 string) string {
	if path1 == "/" && path2 == "/" {
		return "/"
	}

	path1 = path.Clean(path1)
	path2 = path.Clean(path2)

	var longerPath, shorterPath string
	if len(path1) >= len(path2) {
		longerPath = path1
		shorterPath = path2
	} else {
		longerPath = path2
		shorterPath = path1
	}

	sharedStr := ""

	for i:=0; i<len(shorterPath); i++ {
		if shorterPath[i] == longerPath[i] {
			sharedStr += string(shorterPath[i])
		} else {
			if sharedStr == "/" {
				sharedStr = ""
			}
			return sharedStr
		}
	}

	return sharedStr
}


func cleanURLPath(targetUrl string) string {
	parsed, _ := url.Parse(targetUrl)
	cleaned := path.Clean(parsed.Path)
	if cleaned == "." || cleaned == "" || cleaned == " " {
		cleaned = "/"
	}
	return cleaned
}


// Match a URL path to a local path.
func urlToSandboxPath(targetUrl string, sandbox Sandbox) (string, error) {
	parsed := cleanURLPath(targetUrl)

	localPath := path.Clean(sandbox.LocalPath)
	if string(localPath[len(localPath)-1]) != "/" {
		localPath += "/"
	}

	fullLocalPath := strings.Replace(parsed, sandbox.Path, localPath, 1)
	fullLocalPath = path.Clean(fullLocalPath)
	fileInfo, fileErr := os.Stat(fullLocalPath)

	if fileErr == nil && fileInfo.IsDir() {
		return filepath.Join(fullLocalPath, "index.gmi"), nil
	} else {
		return fullLocalPath, nil
	}

	return "", fileErr
}


func isAllZeroes(bytes []byte) bool {
	for _, v := range bytes {
		if v != 0 {
			return false
		}
	}
	return true
}


func handleConnection(s *Server, c net.Conn) {
	data := make([]byte, 1024)
	c.Read(data)

	if isAllZeroes(data) {
		return
	}

	dataStr := requestAsString(data)
	requestParsed, err := url.Parse(dataStr)

	if s.EnableLog {
		log.Output(1, fmt.Sprintf("%s -> %s", c.RemoteAddr().String(), dataStr))
	}

	if err != nil {
		fmt.Println("Error occurred when parsing URL!")
	}

	// Usually happens when a bot probes the server with a
	// blank HTTP request. This saves it from crashing.
	if (requestParsed == nil) {
		c.Close()
		log.Output(1, fmt.Sprintf("Ignored request to `%s` and closed connection.", requestParsed))
		return
	}

	context := NewContext(dataStr, c)

	// This if statement handles rate-limiting. There's
	// a lot more depth to it, but without this if statement,
	// there is no rate-limiting at all.
	if !allowConnection(context, 1) {
		return
	}

	cleanedPath := cleanURLPath(requestParsed.Path)
	handledAsSandbox := false

	// First, see if there is a static file to serve
	// from a sandbox.
	for _, sandbox := range s.Router.Sandboxes {
		if getSharedPath(sandbox.Path, cleanedPath) == sandbox.Path {
			fullLocalPath, _ := urlToSandboxPath(dataStr, sandbox)
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
