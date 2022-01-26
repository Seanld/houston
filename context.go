package houston


import (
	"net"
	"net/url"
	"log"
)


type Context struct {
	URL         string
	Connection  net.Conn
}


func NewContext(newUrl string, newConn net.Conn) Context {
	return Context{newUrl, newConn}
}


func (c *Context) GetQuery() string {
	parsed, err := url.Parse(c.URL)
	if err != nil {
		log.Fatalf("Error when gettig query string: %v", err)
	}
	return parsed.RawQuery
}
