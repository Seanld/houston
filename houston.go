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
func GetSharedPath(path1 string, path2 string) string {
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

	if sharedStr == "/" {
		sharedStr = ""
	}
	return sharedStr
}


// Match a URL path to a local path.
func URLToSandboxPath(targetUrl string, sandbox Sandbox) (string, error) {
	parsed, _ := url.Parse(targetUrl)
	fullLocalPath := strings.Replace(parsed.Path, sandbox.Path, sandbox.LocalPath, 1)
	fullLocalPath = path.Clean(fullLocalPath)
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

	if s.EnableLog {
		log.Output(1, fmt.Sprintf("%s -> %s", c.RemoteAddr().String(), dataStr))
	}

	if err != nil {
		fmt.Println("Error occurred when parsing URL!")
	}

	context := NewContext(dataStr, c)

	cleanedPath := path.Clean(requestParsed.Path)
	handledAsSandbox := false

	// First, see if there is a static file to serve
	// from a sandbox.
	for _, sandbox := range s.Router.Sandboxes {
		if GetSharedPath(sandbox.Path, cleanedPath) != "" {
			fullLocalPath, _ := URLToSandboxPath(dataStr, sandbox)
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
