package houston

import (
	"errors"

	"golang.org/x/time/rate"
)

// Rate-limited connection. Stores the IP address, and
// the pointer to the rate-limiter it's associated with.
type LimitedConn struct {
	IP string
	Limiter *rate.Limiter
}


var connPool []*LimitedConn = []*LimitedConn{}


func PoolFindEntry(addr string) (*LimitedConn, error) {
	for _, limitedConn := range connPool {
		if limitedConn.IP == addr {
			return limitedConn, nil
		}
	}
	return nil, errors.New("Connection not found in pool!")
}


// TODO Implement request rate-limit-checking.
// func PoolCheckEntry(addr string)


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


func PoolCreateEntry(addr string, rateLimit rate.Limit, burst int) {
	newEntry := LimitedConn{
		IP: addr,
		Limiter: rate.NewLimiter(rateLimit, burst),
	}

	connPool = append(connPool, &newEntry)
}
