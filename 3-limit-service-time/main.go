//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync/atomic"
	"time"
)

var FreeTierLimit int64 = 10 // in seconds

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if process == nil || u == nil {
		return false
	}

	if u.IsPremium {
		process()
		return true
	}

	ticker := time.NewTicker(time.Second)
	done := make(chan bool)
	success := make(chan bool)

	go func() {
		process()
		close(done)
	}()

	go func() {
		defer func() {
			close(success)
		}()
		for {
			select {
			case <-done:
				success <- true
				return
			case <-ticker.C:
				if used := atomic.AddInt64(&u.TimeUsed, 1); used >= FreeTierLimit {
					return
				}
			}
		}
	}()

	return <-success
}

func main() {
	RunMockServer()
}
