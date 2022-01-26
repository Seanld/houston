package houston


import (
	"fmt"
	"io/ioutil"
)


func formatResponse(code int, meta string) []byte {
	return []byte(fmt.Sprintf("%d %s\r\n", code, meta))
}


func Success(ctx Context, mimeType string) {
	ctx.Connection.Write(formatResponse(20, mimeType))
}


// Send content data to the client.
func SendBytes(ctx Context, mimeType string, content []byte) {
	Success(ctx, mimeType)
	ctx.Connection.Write(content)
}


func SendString(ctx Context, mimeType string, str string) {
	SendBytes(ctx, mimeType, []byte(str))
}


func SendStringf(ctx Context, mimeType string, template string, values ...string) {
	formatted := fmt.Sprintf(template, values)
	SendBytes(ctx, mimeType, []byte(formatted))
}


func SendFile(ctx Context, mimeType string, path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	SendBytes(ctx, mimeType, content)
	return nil
}


///////////////
// 1X INPUTS //
///////////////


type InputHandler func(string, Context)


func Input(ctx Context, prompt string) {
	ctx.Connection.Write(formatResponse(10, prompt))
}


func InputAndDo(ctx Context, prompt string, handler InputHandler) {
	queryString := ctx.GetQuery()
	if queryString != "" {
		handler(queryString, ctx)
	} else {
		Input(ctx, prompt)
	}
}


func SensitiveInput(ctx Context, prompt string) {
	ctx.Connection.Write(formatResponse(11, prompt))
}


//////////////////
// 3X REDIRECTS //
//////////////////


func RedirectTemp(ctx Context, url string) {
	ctx.Connection.Write(formatResponse(30, url))
}

func RedirectPerm(ctx Context, url string) {
	ctx.Connection.Write(formatResponse(31, url))
}


///////////////////////////
// 4X TEMPORARY FAILURES //
///////////////////////////


func TempFail(ctx Context, info string) {
	ctx.Connection.Write(formatResponse(40, info))
}

func ServerUnavailable(ctx Context, info string) {
	ctx.Connection.Write(formatResponse(41, info))
}

func CGIError(ctx Context, info string) {
	ctx.Connection.Write(formatResponse(42, info))
}

func ProxyError(ctx Context, info string) {
	ctx.Connection.Write(formatResponse(43, info))
}

func SlowDown(ctx Context, waitSeconds int) {
	ctx.Connection.Write([]byte(fmt.Sprintf("44 %d\r\n", waitSeconds)))
}


///////////////////////////
// 5X PERMANENT FAILURES //
///////////////////////////


func PermFailure(ctx Context, info string) {
	ctx.Connection.Write(formatResponse(50, info))
}

func NotFound(ctx Context, info string) {
	ctx.Connection.Write(formatResponse(51, info))
}

func Gone(ctx Context, info string) {
	ctx.Connection.Write(formatResponse(52, info))
}

func ProxyRequestRefused(ctx Context, info string) {
	ctx.Connection.Write(formatResponse(53, info))
}

func BadRequest(ctx Context, info string) {
	ctx.Connection.Write(formatResponse(59, info))
}


////////////////////////////////////
// 6X CLIENT CERTIFICATE REQUIRED //
////////////////////////////////////


func ClientCertRequired(ctx Context, info string) {
	ctx.Connection.Write(formatResponse(60, info))
}

func CertNotAuthorized(ctx Context, info string) {
	ctx.Connection.Write(formatResponse(61, info))
}

func CertNotValid(ctx Context, info string) {
	ctx.Connection.Write(formatResponse(62, info))
}
