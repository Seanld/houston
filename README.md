# Houston

Houston is an Express-like Gemini server, written in Go. Primarily because
I want to. Not because it's necessary. There are plenty other Gemini servers
that have all important functionality covered, and several are written in Go
themselves. I just want to make my own.

That said, this aims to be lightweight and easy to use. This was initially
developed to suite my own purposes, for my own capsules, as I wanted a server
that I understood. Not somebody else's.


## Basic Usage Example

Easiest way to learn, for me, is by reading an example. Here you go:

```go
package main
    
import (
    "fmt"
    "git.sr.ht/~seanld/houston"
    "net"
)
    
func main() {
    mainRouter := houston.BlankRouter()
    
    mainRouter.AddSandbox("/", "sandbox")
    
    mainRouter.AddRoute("/other", func(c net.Conn) {
        houston.SendString(c, "text/plain", "Hello, world!")
    })
    
    mainRouter.AddRoute("/hello", func(c net.Conn) {
        houston.SendString(c, "text/gemini", "# Why, hello to you!")
    })
    
    newServer := houston.NewServer(mainRouter, "certificate.crt", "private.key")
    
    fmt.Println("Starting server...")
    newServer.Start("0.0.0.0")
}
```

## More Info

`Router` structs are used by `Server` structs to provide functionality for handling
request-to-response. `Routers` can have `Route` and `Sandbox` instances. They can be
added to a router by doing `Router.AddRoute(url, func (net.Conn) {})` or
`Router.AddSandbox(url, sandboxDirPath)`.

`Route` instances connect a URL path to a function that is executed when it's visited.
You can use the responses from `responses.go` to make responses easier.

`Sandbox` instances connect a URL path to a local directory that holds static files.
For example, if you connect `/hello` to local dir `/hello-static`, and `/hello-static`
has a file named `index.gmi` in it, and someone visits `/hello`, it will attempt
to load the file `/hello-static/index.gmi`. Or any other file specified from that dir.

