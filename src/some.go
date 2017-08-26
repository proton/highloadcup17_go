package main

import (
	"unsafe"
)

func bstring(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// from https://github.com/valyala/fasthttp/blob/master/bytesconv.go
func ParseUint32(buf []byte) (*uint32, bool) {
	v, n, ok := parseUint32Buf(buf)
	if ok {
		if n == len(buf) {
			return &v, true
		}
	}
	return nil, false
}

func parseUint32Buf(b []byte) (uint32, int, bool) {
	n := len(b)
	if n == 0 {
		return 0, 0, false
	}
	v := uint32(0)
	for i := 0; i < n; i++ {
		c := b[i]
		k := c - '0'
		if k > 9 {
			if i == 0 {
				return 0, i, false
			}
			return v, i, true
		}
		if i >= 18 {
			return 0, i, false
		}
		v = 10*v + uint32(k)
	}
	return v, n, true
}
