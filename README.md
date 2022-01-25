# Houston

Houston is a Gemini server, written in Go. Primarily because I want to.
Not because it's necessary. There are plenty other Gemini servers that
have all important functionality covered, and several are written in Go
themselves. I just want to make my own.

That said, this aims to be lightweight and easy to use. I will be using
Houston to host my own capsule.


## Basic Usage Example

Easiest way to learn, for me, is by reading an example. Here you go:

    package main
    
    import (
      "fmt"
      "git.sr.ht/~seanld/houston"
      "net"
    )
    
    func main() {
      mainRouter := houston.NewRouter(houston.RouterOptions{
        ErrorHandler: func(c net.Conn) {
          houston.PermFailure(c, "Failure!")
        },
      })
    
      mainRouter.AddRoute("/", func(c net.Conn) {
        houston.SendString(c, "text/plain", "Hello, world!")
      })
    
      mainRouter.AddRoute("/hi", func(c net.Conn) {
        houston.SendString(c, "text/gemini", "# Why, hello to you!")
      })
    
      newServer := houston.NewServer(mainRouter, "certificate.crt", "private.key")
    
      fmt.Println("Starting server...")
      newServer.Start()
    }


## Roadmap

The only goal I currently have is to host static files in a sandboxed
directory. Once I achieve that, I'll move on to some cool features.

