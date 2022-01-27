package houston


import (
	"net"
	"net/url"
	"fmt"
	"path"
	"path/filepath"
	"mime"
	"os"
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


// Match a URL path to a local path.
func URLToSandboxPath(targetUrl string, sandboxBasePath string) (string, error) {
	parsed, _ := url.Parse(targetUrl)
	fullLocalPath := filepath.Join(sandboxBasePath, parsed.Path)
	fileInfo, fileErr := os.Stat(fullLocalPath)

	if fileErr == nil && fileInfo.IsDir() {
		return filepath.Join(fullLocalPath, "index.gmi"), nil
	} else {
		if filepath.Ext(fullLocalPath) != "" {
			return fullLocalPath, nil
		} else {
			return fullLocalPath + ".gmi", nil
		}
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
	for _, sandbox := range s.Router.Sandboxes {
		dir := path.Dir(cleanedPath)
		if dir == "." {
			dir = "/"
		}
		if sandbox.Path == dir {
			fullLocalPath, _ := URLToSandboxPath(dataStr, sandbox.LocalPath)
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
