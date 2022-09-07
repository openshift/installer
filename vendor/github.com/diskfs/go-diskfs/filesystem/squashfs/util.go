package squashfs

import (
	"strings"
)

const (
	// KB represents one KB
	KB int64 = 1024
	// MB represents one MB
	MB int64 = 1024 * KB
	// GB represents one GB
	GB int64 = 1024 * MB
	// TB represents one TB
	TB int64 = 1024 * GB
	// max value of uint32
	uint32max uint64 = 0xffffffff
)

func universalizePath(p string) (string, error) {
	// globalize the separator
	ps := strings.Replace(p, `\`, "/", -1)
	//if ps[0] != '/' {
	//return "", errors.New("Must use absolute paths")
	//}
	return ps, nil
}
func splitPath(p string) ([]string, error) {
	ps, err := universalizePath(p)
	if err != nil {
		return nil, err
	}
	// we need to split such that each one ends in "/", except possibly the last one
	parts := strings.Split(ps, "/")
	// eliminate empty parts
	ret := make([]string, 0)
	for _, sub := range parts {
		if sub != "" {
			ret = append(ret, sub)
		}
	}
	return ret, nil
}
