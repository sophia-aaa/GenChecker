//go:build !appengine
// +build !appengine

package bigcache

import (
	"fmt"
	"unsafe"
)

func bytesToString(b []byte) string {
	fmt.Println(b)
	return string(b[:])
}