package houston

import (
	"errors"
	"fmt"
	"time"
	"strings"

	"golang.org/x/time/rate"
)

// Rate-limited connection. Stores the IP address, and
// the pointer to the rate-limiter it's associated with.
type LimitedConn struct {
	IP string
	Limiter *rate.Limiter
	Reservation *rate.Reservation
}


var connPool []*LimitedConn = []*LimitedConn{}


func PoolFindEntry(addr string) (*LimitedConn, error) {
	for _, limitedConn := range connPool {
		fmt.Println(limitedConn.IP, addr)
		if limitedConn.IP == addr {
			return limitedConn, nil
		}
	}
	return nil, errors.New("Connection not found in pool!")
}


func PoolCheckEntry(addr string) (bool, error) {
	limitedConn, err := PoolFindEntry(addr)
	if err != nil {
		return false, err
	}
	return limitedConn.Limiter.Allow(), nil
}


func PoolReserveEntry(addr string, tokens int) (*rate.Reservation, error) {
	limitedConn, err := PoolFindEntry(addr)
	if err != nil {
		return nil, err
	}
	return limitedConn.Limiter.ReserveN(time.Now(), tokens), nil
}


func PoolDeleteEntry(addr string) {
	idx := 0
	found := false

	for _, limitedConn := range connPool {
		if limitedConn.IP == addr {
			found = true
			break
		}
	}

	if found == true {
		newPool := make([]*LimitedConn, 0)
		newPool = append(newPool, connPool[:idx]...)
		connPool = append(newPool, connPool[idx+1:]...)
	}
}


func PoolCreateEntry(addr string, rateLimit rate.Limit, burst int) *LimitedConn {
	fmt.Println("Entry created!")
	newEntry := LimitedConn{
		IP: addr,
		Limiter: rate.NewLimiter(rateLimit, burst),
	}
	connPool = append(connPool, &newEntry)
	return &newEntry
}


// Checks if the connection's request will be allowed
// to continue, or if it will be denied due to rate limits.
// Returns a "slow down" response, and closes connection if condition
// is denied. If allowed, add token to rate limiter bucket.
func allowConnection(ctx Context, tokens int) bool {
	var entryPtr *LimitedConn
	addrWithPort := ctx.Connection.RemoteAddr().String()
	addr := strings.Split(addrWithPort, ":")[0]

	entryPtr, err := PoolFindEntry(addr)
	if err != nil {
		// Entry doesn't exist in pool yet, so create one.
		// TODO Allow customizing of the rate and bucket size
		// via config.
		entryPtr = PoolCreateEntry(addr, 2, 2)
	}
	fmt.Println(entryPtr)

	if entryPtr.Limiter.AllowN(time.Now(), tokens) {
		// Connection allowed, so use up some of the client's tokens.
		entryPtr.Reservation = entryPtr.Limiter.ReserveN(time.Now(), tokens)
		return true
	} else {
		ctx.SlowDown(int(entryPtr.Reservation.Delay().Seconds()))
		ctx.Connection.Close()
		return false
	}
}
