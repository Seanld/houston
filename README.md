Houston is an Express-like Gemini server, written in Go. Primarily because
I want to. Not because it's necessary. There are plenty other Gemini servers
that have all important functionality covered, and several are written in Go
themselves. I just want to make my own.

That said, this aims to be lightweight and easy to use. This was initially
developed to suite my own purposes, for my own capsules, as I wanted a server
that I understood. Not somebody else's.


# Roadmap

My initial goal for Houston was to serve static files. It's moved beyond that
already. Here are my goals for the project:

-   [X] Serve static directories
-   [X] Handle requests with user-defined functions
-   [X] Basic logging (a toggle for major events like incoming connections, nothing
    fancy &#x2013; to keep it simple).
-   [ ] Rate-limiting capabilities, to prevent DOS attacks and spam.
-   [X] Intuitive support for templates.
-   [ ] [Titan protocol](https://transjovian.org:1965/titan/page/The%20Titan%20Specification) integration.
    -   [Lagrange](https://github.com/skyjake/lagrange) is currently the only client I know of that implements this on the
        user-side, so not super high priority.


# Basic Usage Example

Easiest way to learn, for me, is by reading an example. Here you go:

```go
package main
    
import (
    "fmt"
    "git.sr.ht/~seanld/houston"
)
    
func main() {
    r := houston.BlankRouter()

    // Route a URL path to a static file directory, like:
    // gemini://localhost/ -> ./sandbox/index.gmi
    // gemini://localhost/hello.gmi -> ./sandbox/hello.gmi
    r.Sandbox("/", "sandbox")

    // Run a function when gemini://localhost/interact is visited.
    r.Handle("/interact", func(ctx houston.Context) {
        // Send input response to client if no query string is provided,
        // and then run a function with the entered input value.
        ctx.InputAndDo("Enter name", func(s string) {
            ctx.SendStringf("text/gemini", "Hello, %s, you are #%d.", s, 1)
        })
    })

    // Provide a router, certificate file, and private key file, and enable
    // basic connection logging.
    s := houston.NewServer(r, "main.crt", "my.key", true)

    fmt.Println("Starting server...")
    s.Start("localhost")
}
```


# More Info

`Router` structs are used by `Server` structs to provide functionality for handling
request-to-response. `Routers` can have `Route` and `Sandbox` instances. They can be
added to a router by doing `Router.AddRoute(url, func (net.Conn) {})` or
`Router.AddSandbox(url, sandboxDirPath)`.

`Route` instances connect a URL path to a function that is executed when it's visited.

`Sandbox` instances connect a URL path to a local directory that holds static files.
For example, if you connect `/hello` to local dir `/hello-static`, and `/hello-static`
has a file named `index.gmi` in it, and someone visits `/hello`, it will attempt
to load the file `/hello-static/index.gmi`. Or any other file specified from that dir.

`Context` instances provide the URL of a connection, the actual `net.Conn` of the
connection, some methods for conveniently sending responses, and other features.

