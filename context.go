package houston


import (
	"net"
	"net/url"
	"log"
	"fmt"
	"io/ioutil"
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
		log.Fatalf("Error when getting query string: %v", err)
	}
	return parsed.RawQuery
}


///////////////
// RESPONSES //
///////////////


func formatResponse(code int, meta string) []byte {
	return []byte(fmt.Sprintf("%d %s\r\n", code, meta))
}


func (ctx *Context) Success(mimeType string) {
	ctx.Connection.Write(formatResponse(20, mimeType))
}


// Send content data to the client.
func (ctx *Context) SendBytes(mimeType string, content []byte) {
	ctx.Success(mimeType)
	ctx.Connection.Write(content)
}


func (ctx *Context) SendString(mimeType string, str string) {
	ctx.SendBytes(mimeType, []byte(str))
}


func (ctx *Context) SendStringf(mimeType string, str string, values ...interface{}) {
	formatted := fmt.Sprintf(str, values...)
	ctx.SendBytes(mimeType, []byte(formatted))
}


func (ctx *Context) SendFile(mimeType string, path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	ctx.SendBytes(mimeType, content)
	return nil
}


func (ctx *Context) SendTemplate(mimeType string, path string, data interface{}) error {
	rendered, err := Template(path, data)
	if err != nil {
		return err
	}
	ctx.SendString(mimeType, rendered)
	return nil
}


///////////////
// 1X INPUTS //
///////////////


type InputHandler func(string, Context)


func (ctx *Context) Input(prompt string) {
	ctx.Connection.Write(formatResponse(10, prompt))
}


func (ctx *Context) InputAndDo(prompt string, handler InputHandler) {
	queryString := ctx.GetQuery()
	if queryString != "" {
		handler(queryString, *ctx)
	} else {
		ctx.Input(prompt)
	}
}


func (ctx *Context) SensitiveInput(prompt string) {
	ctx.Connection.Write(formatResponse(11, prompt))
}


func (ctx *Context) SensitiveInputAndDo(prompt string, handler InputHandler) {
	queryString := ctx.GetQuery()
	if queryString != "" {
		handler(queryString, *ctx)
	} else {
		ctx.SensitiveInput(prompt)
	}
}


//////////////////
// 3X REDIRECTS //
//////////////////


func (ctx *Context) RedirectTemp(url string) {
	ctx.Connection.Write(formatResponse(30, url))
}

func (ctx *Context) RedirectPerm(url string) {
	ctx.Connection.Write(formatResponse(31, url))
}


///////////////////////////
// 4X TEMPORARY FAILURES //
///////////////////////////


func (ctx *Context) TempFail(info string) {
	ctx.Connection.Write(formatResponse(40, info))
}

func (ctx *Context) ServerUnavailable(info string) {
	ctx.Connection.Write(formatResponse(41, info))
}

func (ctx *Context) CGIError(info string) {
	ctx.Connection.Write(formatResponse(42, info))
}

func (ctx *Context) ProxyError(info string) {
	ctx.Connection.Write(formatResponse(43, info))
}

func (ctx *Context) SlowDown(waitSeconds int) {
	ctx.Connection.Write([]byte(fmt.Sprintf("44 %d\r\n", waitSeconds)))
}


///////////////////////////
// 5X PERMANENT FAILURES //
///////////////////////////


func (ctx *Context) PermFailure(info string) {
	ctx.Connection.Write(formatResponse(50, info))
}

func (ctx *Context) NotFound(info string) {
	ctx.Connection.Write(formatResponse(51, info))
}

func (ctx *Context) Gone(info string) {
	ctx.Connection.Write(formatResponse(52, info))
}

func (ctx *Context) ProxyRequestRefused(info string) {
	ctx.Connection.Write(formatResponse(53, info))
}

func (ctx *Context) BadRequest(info string) {
	ctx.Connection.Write(formatResponse(59, info))
}


////////////////////////////////////
// 6X CLIENT CERTIFICATE REQUIRED //
////////////////////////////////////


func (ctx *Context) ClientCertRequired(info string) {
	ctx.Connection.Write(formatResponse(60, info))
}

func (ctx *Context) CertNotAuthorized(info string) {
	ctx.Connection.Write(formatResponse(61, info))
}

func (ctx *Context) CertNotValid(info string) {
	ctx.Connection.Write(formatResponse(62, info))
}
