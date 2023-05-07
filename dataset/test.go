package dataset

import (
	"reflect"
	"unsafe"
)

/*
	type GenHeader[T any] struct {
		List []T
	}

	func (g *GenHeader[T]) Gens() []T {
		return g.List
	}

	func (g *GenHeader[T]) SetGI(i int, x T) {
		g.List[i] = x
	}

	func (g *GenHeader[T]) GetGI(i int) T {
		return g.List[i]
	}
*/
func (a array) Data() interface{} {
	// build a type of []T
	shdr := reflect.SliceHeader{
		Data: a.Uintptr(),
		Len:  a.Len(),
		Cap:  a.Cap(),
	}
	sliceT := reflect.SliceOf(a.t.Type)
	ptr := unsafe.Pointer(&shdr)
	val := reflect.Indirect(reflect.NewAt(sliceT, ptr))
	return val.Interface()

}
