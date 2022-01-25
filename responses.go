package houston


import (
	"net"
	"fmt"
)


func formatResponse(code int, meta string) string {
	return fmt.Sprintf("%d %s\r\n", code, meta)
}


func Success(c net.Conn, mimeType string) string {
	return formatResponse(20, mimeType)
}


///////////////
// 1X INPUTS //
///////////////


func Input(c net.Conn, prompt string) string {
	return formatResponse(10, prompt)
}


func SensitiveInput(c net.Conn, prompt string) string {
	return formatResponse(11, prompt)
}


//////////////////
// 3X REDIRECTS //
//////////////////


func RedirectTemp(c net.Conn, url string) string {
	return formatResponse(30, url)
}

func RedirectPerm(c net.Conn, url string) string {
	return formatResponse(31, url)
}


///////////////////////////
// 4X TEMPORARY FAILURES //
///////////////////////////


func TempFail(c net.Conn, info string) string {
	return formatResponse(40, info)
}

func ServerUnavailable(c net.Conn, info string) string {
	return formatResponse(41, info)
}

func CGIError(c net.Conn, info string) string {
	return formatResponse(42, info)
}

func CGIError(c net.Conn, info string) string {
	return formatResponse(43, info)
}

func CGIError(c net.Conn, info string) string {
	return formatResponse(44, info)
}


///////////////////////////
// 5X PERMANENT FAILURES //
///////////////////////////


func PermFailure(c net.Conn, info string) string {
	return formatResponse(50, info)
}

func NotFound(c net.Conn, info string) string {
	return formatResponse(51, info)
}

func Gone(c net.Conn, info string) string {
	return formatResponse(52, info)
}

func ProxyRequestRefused(c net.Conn, info string) string {
	return formatResponse(53, info)
}

func BadRequest(c net.Conn, info string) string {
	return formatResponse(59, info)
}


////////////////////////////////////
// 6X CLIENT CERTIFICATE REQUIRED //
////////////////////////////////////


func ClientCertRequired(c net.Conn, info string) string {
	return formatResponse(60, info)
}

func CertNotAuthorized(c net.Conn, info string) string {
	return formatResponse(61, info)
}

func CertNotValid(c net.Conn, info string) string {
	return formatResponse(62, info)
}
