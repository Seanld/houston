package houston


import (
	"net"
	"fmt"
)


func formatResponse(code int, meta string) string {
	return fmt.Sprintf("%d %s\r\n", code, meta)
}


func Success(conn net.Conn, mimeType string) string {
	return formatResponse(20, mimeType)
}


///////////////
// 1X INPUTS //
///////////////


func Input(conn net.Conn, prompt string) string {
	return formatResponse(10, prompt)
}


func SensitiveInput(conn net.Conn, prompt string) string {
	return formatResponse(11, prompt)
}


//////////////////
// 3X REDIRECTS //
//////////////////


func RedirectTemp(conn net.Conn, url string) string {
	return formatResponse(30, url)
}

func RedirectPerm(conn net.Conn, url string) string {
	return formatResponse(31, url)
}


///////////////////////////
// 4X TEMPORARY FAILURES //
///////////////////////////


func TempFail(conn net.Conn, info string) string {
	return formatResponse(40, info)
}

func ServerUnavailable(conn net.Conn, info string) string {
	return formatResponse(41, info)
}

func CGIError(conn net.Conn, info string) string {
	return formatResponse(42, info)
}

func CGIError(conn net.Conn, info string) string {
	return formatResponse(43, info)
}

func CGIError(conn net.Conn, info string) string {
	return formatResponse(44, info)
}


///////////////////////////
// 5X PERMANENT FAILURES //
///////////////////////////


func PermFailure(conn net.Conn, info string) string {
	return formatResponse(50, info)
}

func NotFound(conn net.Conn, info string) string {
	return formatResponse(51, info)
}

func Gone(conn net.Conn, info string) string {
	return formatResponse(52, info)
}

func ProxyRequestRefused(conn net.Conn, info string) string {
	return formatResponse(53, info)
}

func BadRequest(conn net.Conn, info string) string {
	return formatResponse(59, info)
}


////////////////////////////////////
// 6X CLIENT CERTIFICATE REQUIRED //
////////////////////////////////////


func ClientCertRequired(conn net.Conn, info string) string {
	return formatResponse(60, info)
}

func CertNotAuthorized(conn net.Conn, info string) string {
	return formatResponse(61, info)
}

func CertNotValid(conn net.Conn, info string) string {
	return formatResponse(62, info)
}
