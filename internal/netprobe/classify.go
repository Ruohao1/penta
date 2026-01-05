package netprobe

import (
	"errors"
	"net"
	"os"
	"syscall"
)

func classifyDialErr(result *Result, err error) *Result {
	// context deadline
	result.OK = false
	result.State = "closed"
	result.Err = err
	if ne := (net.Error)(nil); errors.As(err, &ne) && ne.Timeout() {
		result.Reason = "timeout"
		return result
	}

	// unwrap common syscall errors
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		// opErr.Err can be *os.SyscallError or syscall.Errno
		var se *os.SyscallError
		if errors.As(opErr.Err, &se) {
			if errno, ok := se.Err.(syscall.Errno); ok {
				return fromErrno(result, errno)
			}
		}
		if errno, ok := opErr.Err.(syscall.Errno); ok {
			return fromErrno(result, errno)
		}
	}

	result.Reason = "error"
	return result
}

func fromErrno(result *Result, errno syscall.Errno) *Result {
	switch errno {
	case syscall.ECONNREFUSED:
		result.OK = true
		result.State = "open"
		result.Reason = "refused"
	case syscall.ENETUNREACH, syscall.EHOSTUNREACH:
		result.Reason = "unreachable"
	}
	return result
}
