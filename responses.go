package houston


import (
	"net"
	"fmt"
)


func formatResponse(code int, meta string) {
	return fmt.Sprintf("%d %s\r\n", code, meta)
}


func Success(conn net.Conn, mimeType string) {
	conn.Write(formatResponse(20, mimeType))
}


///////////////
// 1X INPUTS //
///////////////


func Input(conn net.Conn, prompt string) {
	conn.Write(formatResponse(10, prompt))
}


func SensitiveInput(conn net.Conn, prompt string) {
	conn.Write(formatResponse(11, prompt))
}


//////////////////
// 3X REDIRECTS //
//////////////////


func RedirectTemp(conn net.Conn, url string) {
	conn.Write(formatResponse(30, url))
}

func RedirectPerm(conn net.Conn, url string) {
	conn.Write(formatResponse(31, url))
}


///////////////////////////
// 4X TEMPORARY FAILURES //
///////////////////////////


func TempFail(conn net.Conn, info string) {
	conn.Write(formatResponse(40, info))
}

func ServerUnavailable(conn net.Conn, info string) {
	conn.Write(formatResponse(41, info))
}

func CGIError(conn net.Conn, info string) {
	conn.Write(formatResponse(42, info))
}

func ProxyError(conn net.Conn, info string) {
	conn.Write(formatResponse(43, info))
}

func SlowDown(conn net.Conn, waitSeconds int) {
	conn.Write(Sprintf("44 %d\r\n", waitSeconds))
}


///////////////////////////
// 5X PERMANENT FAILURES //
///////////////////////////


func PermFailure(conn net.Conn, info string) {
	conn.Write(formatResponse(50, info))
}

func NotFound(conn net.Conn, info string) {
	conn.Write(formatResponse(51, info))
}

func Gone(conn net.Conn, info string) {
	conn.Write(formatResponse(52, info))
}

func ProxyRequestRefused(conn net.Conn, info string) {
	conn.Write(formatResponse(53, info))
}

func BadRequest(conn net.Conn, info string) {
	conn.Write(formatResponse(59, info))
}


////////////////////////////////////
// 6X CLIENT CERTIFICATE REQUIRED //
////////////////////////////////////


func ClientCertRequired(conn net.Conn, info string) {
	conn.Write(formatResponse(60, info))
}

func CertNotAuthorized(conn net.Conn, info string) {
	conn.Write(formatResponse(61, info))
}

func CertNotValid(conn net.Conn, info string) {
	conn.Write(formatResponse(62, info))
}
