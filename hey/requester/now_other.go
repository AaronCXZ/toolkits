// +build !windows

package requester

import "time"

var startTime = time.Now()

func now() time.Duration {
	return time.Since(startTime)
}
