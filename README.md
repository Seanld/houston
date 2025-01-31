Houston is an Express-like Gemini server, written in Go. Primarily because I want to. Not because it&rsquo;s necessary. There are plenty other Gemini servers that have all important functionality covered, and several are written in Go themselves. I just want to make my own.

That said, this aims to be lightweight and easy to use. This was initially developed to suite my own purposes, for my own capsules, as I wanted a server that I understood. Not somebody else&rsquo;s.


# Roadmap

My initial goal for Houston was to serve static files. It&rsquo;s moved beyond that already. Here are my goals for the project:

-   [X] Serve static directories
-   [X] Handle requests with user-defined functions
-   [X] Basic logging (a toggle for major events like incoming connections, nothing fancy &#x2013; to keep it simple).
-   [X] Rate-limiting capabilities, to prevent DOS attacks and spam.
-   [X] Intuitive support for templates.
-   [ ] [Titan protocol](https://transjovian.org:1965/titan/page/The%20Titan%20Specification) integration.
    -   [Lagrange](https://github.com/skyjake/lagrange) is currently the only client I know of that implements this on the user-side, so not super high priority.


# Basic Usage Example

Easiest way to learn, for me, is by reading an example. Here you go:

```go
package main

import (
    "strconv"

    "git.sr.ht/~seanld/houston"
)

func main() {
    r := houston.BlankRouter()

    // Serve static files from directory `./static` at URL `/`.
    // For example:
    // * Request to `/` yields file `./static/index.gmi`
    // * Request to `/something/example.txt` yields file `./static/something/example.txt`
    r.Sandbox("/", "static")

    // You can handle requests more dynamically and programatically
    // by using `Router.Handle()`, which takes a callback function.
    // The callback is given a `Context` instance, which contains
    // information about the client and the request that can be
    // operated on however you want.
    r.Handle("/input", func(ctx houston.Context) {
        ctx.InputAndDo("Guess a number from 1 to 10", func(s string) {
            asInt, err := strconv.ParseInt(s, 10, 32)
            if err != nil {
                ctx.BadRequest("Please enter an integer!")
                return
            }
            if asInt == 3 {
                ctx.SendString("text/plain", "You got it right!")
            } else {
                ctx.SendString("text/plain", "Sorry, you're incorrect.")
            }
        })
    })

    newServer := houston.NewServer(&r, &houston.ServerConfig{
        // Set the server up with a TLS certificate and key.
        // Self-signed is acceptable and normal in Gemini.
        CertificatePath: "cert/localhost.crt",
        KeyPath: "cert/localhost.key",

        // If you're trying to put your server into production,
        // you'll want to specify a hostname different from
        // `localhost`. Port defaults to 1965.
        Hostname: "0.0.0.0",
        Port: 1965,

        // You can enable connection logging with Houston.
        // Toggle the boolean flag, and give it the log file's
        // path, and it will record connections to your capsule.
        EnableLog: true,
        LogFilePath: "houston.log",

        // Houston comes with a rate-limiter (token-bucket algorithm),
        // and can be enabled and configured. By default `MaxRate` and
        // `BucketSize` are set to 2. These are good defaults for most.
        EnableLimiting: true,
        MaxRate: 2,
        BucketSize: 2,
    })

    // With our router's URL endpoints set up, and the server
    // options configured, the server can be started up. Just SIGTERM
    // to stop it (Ctrl+C for most machines/terminals).
    newServer.Start()
}
```


# More Info

`Router` structs are used by `Server` structs to provide functionality for handling request-to-response. `Routers` can have `Route` and `Sandbox` instances. They can be added to a router by doing `Router.Handle(url, func (net.Conn) {})` or `Router.Sandbox(url, sandboxDirPath)`.

`Route` instances connect a URL path to a function that is executed when it&rsquo;s visited.

`Sandbox` instances connect a URL path to a local directory that holds static files. For example, if you connect `/hello` to local dir `/hello-static`, and `/hello-static` has a file named `index.gmi` in it, and someone visits `/hello`, it will attempt to load the file `/hello-static/index.gmi`. Or any other file specified from that dir.

`Context` instances provide the URL of a connection, the actual `net.Conn` of the connection, some methods for conveniently sending responses, and other features.

