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
func (h *Header) Bools() []bool {
	return (*(*[]bool)(unsafe.Pointer(&h.Raw)))[:h.TypedLen(bType):h.TypedLen(bType)]
}
func (h *Header) SetB(i int, x bool) { h.Bools()[i] = x }
func (h *Header) GetB(i int) bool    { return h.Bools()[i] }

/* int */

func (h *Header) Ints() []int {
	return (*(*[]int)(unsafe.Pointer(&h.Raw)))[:h.TypedLen(iType):h.TypedLen(iType)]
}
func (h *Header) SetI(i int, x int) { h.Ints()[i] = x }
func (h *Header) GetI(i int) int    { return h.Ints()[i] }

/* int8 */

func (h *Header) Int8s() []int8 {
	return (*(*[]int8)(unsafe.Pointer(&h.Raw)))[:h.TypedLen(i8Type):h.TypedLen(i8Type)]
}
func (h *Header) SetI8(i int, x int8) { h.Int8s()[i] = x }
func (h *Header) GetI8(i int) int8    { return h.Int8s()[i] }

/* int16 */

func (h *Header) Int16s() []int16 {
	return (*(*[]int16)(unsafe.Pointer(&h.Raw)))[:h.TypedLen(i16Type):h.TypedLen(i16Type)]
}
func (h *Header) SetI16(i int, x int16) { h.Int16s()[i] = x }
func (h *Header) GetI16(i int) int16    { return h.Int16s()[i] }

/* int32 */

func (h *Header) Int32s() []int32 {
	return (*(*[]int32)(unsafe.Pointer(&h.Raw)))[:h.TypedLen(i32Type):h.TypedLen(i32Type)]
}
func (h *Header) SetI32(i int, x int32) { h.Int32s()[i] = x }
func (h *Header) GetI32(i int) int32    { return h.Int32s()[i] }

/* int64 */

func (h *Header) Int64s() []int64 {
	return (*(*[]int64)(unsafe.Pointer(&h.Raw)))[:h.TypedLen(i64Type):h.TypedLen(i64Type)]
}
func (h *Header) SetI64(i int, x int64) { h.Int64s()[i] = x }
func (h *Header) GetI64(i int) int64    { return h.Int64s()[i] }

/* uint */

func (h *Header) Uints() []uint {
	return (*(*[]uint)(unsafe.Pointer(&h.Raw)))[:h.TypedLen(uType):h.TypedLen(uType)]
}
func (h *Header) SetU(i int, x uint) { h.Uints()[i] = x }
func (h *Header) GetU(i int) uint    { return h.Uints()[i] }

/* uint8 */

func (h *Header) Uint8s() []uint8 {
	return (*(*[]uint8)(unsafe.Pointer(&h.Raw)))[:h.TypedLen(u8Type):h.TypedLen(u8Type)]
}
func (h *Header) SetU8(i int, x uint8) { h.Uint8s()[i] = x }
func (h *Header) GetU8(i int) uint8    { return h.Uint8s()[i] }

/* uint16 */

func (h *Header) Uint16s() []uint16 {
	return (*(*[]uint16)(unsafe.Pointer(&h.Raw)))[:h.TypedLen(u16Type):h.TypedLen(u16Type)]
}
func (h *Header) SetU16(i int, x uint16) { h.Uint16s()[i] = x }
func (h *Header) GetU16(i int) uint16    { return h.Uint16s()[i] }

func (e E) ReduceFirst(t reflect.Type, data *storage.Header, retVal *storage.Header, split int, size int, fn interface{}) (err error) {
	switch t {
	case Bool:
		dt := data.Bools()
		rt := retVal.Bools()
		switch f := fn.(type) {
		case func([]bool, []bool):
			reduceFirstB(dt, rt, split, size, f)
		case func(bool, bool) bool:
			genericReduceFirstB(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Int:
		dt := data.Ints()
		rt := retVal.Ints()
		switch f := fn.(type) {
		case func([]int, []int):
			reduceFirstI(dt, rt, split, size, f)
		case func(int, int) int:
			genericReduceFirstI(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Int8:
		dt := data.Int8s()
		rt := retVal.Int8s()
		switch f := fn.(type) {
		case func([]int8, []int8):
			reduceFirstI8(dt, rt, split, size, f)
		case func(int8, int8) int8:
			genericReduceFirstI8(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Int16:
		dt := data.Int16s()
		rt := retVal.Int16s()
		switch f := fn.(type) {
		case func([]int16, []int16):
			reduceFirstI16(dt, rt, split, size, f)
		case func(int16, int16) int16:
			genericReduceFirstI16(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Int32:
		dt := data.Int32s()
		rt := retVal.Int32s()
		switch f := fn.(type) {
		case func([]int32, []int32):
			reduceFirstI32(dt, rt, split, size, f)
		case func(int32, int32) int32:
			genericReduceFirstI32(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Int64:
		dt := data.Int64s()
		rt := retVal.Int64s()
		switch f := fn.(type) {
		case func([]int64, []int64):
			reduceFirstI64(dt, rt, split, size, f)
		case func(int64, int64) int64:
			genericReduceFirstI64(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Uint:
		dt := data.Uints()
		rt := retVal.Uints()
		switch f := fn.(type) {
		case func([]uint, []uint):
			reduceFirstU(dt, rt, split, size, f)
		case func(uint, uint) uint:
			genericReduceFirstU(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Uint8:
		dt := data.Uint8s()
		rt := retVal.Uint8s()
		switch f := fn.(type) {
		case func([]uint8, []uint8):
			reduceFirstU8(dt, rt, split, size, f)
		case func(uint8, uint8) uint8:
			genericReduceFirstU8(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Uint16:
		dt := data.Uint16s()
		rt := retVal.Uint16s()
		switch f := fn.(type) {
		case func([]uint16, []uint16):
			reduceFirstU16(dt, rt, split, size, f)
		case func(uint16, uint16) uint16:
			genericReduceFirstU16(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Uint32:
		dt := data.Uint32s()
		rt := retVal.Uint32s()
		switch f := fn.(type) {
		case func([]uint32, []uint32):
			reduceFirstU32(dt, rt, split, size, f)
		case func(uint32, uint32) uint32:
			genericReduceFirstU32(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Uint64:
		dt := data.Uint64s()
		rt := retVal.Uint64s()
		switch f := fn.(type) {
		case func([]uint64, []uint64):
			reduceFirstU64(dt, rt, split, size, f)
		case func(uint64, uint64) uint64:
			genericReduceFirstU64(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Uintptr:
		dt := data.Uintptrs()
		rt := retVal.Uintptrs()
		switch f := fn.(type) {
		case func([]uintptr, []uintptr):
			reduceFirstUintptr(dt, rt, split, size, f)
		case func(uintptr, uintptr) uintptr:
			genericReduceFirstUintptr(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Float32:
		dt := data.Float32s()
		rt := retVal.Float32s()
		switch f := fn.(type) {
		case func([]float32, []float32):
			reduceFirstF32(dt, rt, split, size, f)
		case func(float32, float32) float32:
			genericReduceFirstF32(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Float64:
		dt := data.Float64s()
		rt := retVal.Float64s()
		switch f := fn.(type) {
		case func([]float64, []float64):
			reduceFirstF64(dt, rt, split, size, f)
		case func(float64, float64) float64:
			genericReduceFirstF64(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Complex64:
		dt := data.Complex64s()
		rt := retVal.Complex64s()
		switch f := fn.(type) {
		case func([]complex64, []complex64):
			reduceFirstC64(dt, rt, split, size, f)
		case func(complex64, complex64) complex64:
			genericReduceFirstC64(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Complex128:
		dt := data.Complex128s()
		rt := retVal.Complex128s()
		switch f := fn.(type) {
		case func([]complex128, []complex128):
			reduceFirstC128(dt, rt, split, size, f)
		case func(complex128, complex128) complex128:
			genericReduceFirstC128(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case String:
		dt := data.Strings()
		rt := retVal.Strings()
		switch f := fn.(type) {
		case func([]string, []string):
			reduceFirstStr(dt, rt, split, size, f)
		case func(string, string) string:
			genericReduceFirstStr(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case UnsafePointer:
		dt := data.UnsafePointers()
		rt := retVal.UnsafePointers()
		switch f := fn.(type) {
		case func([]unsafe.Pointer, []unsafe.Pointer):
			reduceFirstUnsafePointer(dt, rt, split, size, f)
		case func(unsafe.Pointer, unsafe.Pointer) unsafe.Pointer:
			genericReduceFirstUnsafePointer(dt, rt, split, size, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	default:
		return errors.Errorf("Unsupported type %v for ReduceFirst", t)
	}
}

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
