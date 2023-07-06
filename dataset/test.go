package dataset

import (
	"reflect"
	"unsafe"
)

func (e E) ReduceLast(t reflect.Type, data *storage.Header, retVal *storage.Header, dimSize int, defaultValue interface{}, fn interface{}) (err error) {
	var ok bool
	switch t {
	case Bool:
		var def bool

		if def, ok = defaultValue.(bool); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Bools()
		rt := retVal.Bools()
		switch f := fn.(type) {
		case func([]bool) bool:
			reduceLastB(dt, rt, dimSize, def, f)
		case func(bool, bool) bool:
			genericReduceLastB(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Int:
		var def int

		if def, ok = defaultValue.(int); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Ints()
		rt := retVal.Ints()
		switch f := fn.(type) {
		case func([]int) int:
			reduceLastI(dt, rt, dimSize, def, f)
		case func(int, int) int:
			genericReduceLastI(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Int8:
		var def int8

		if def, ok = defaultValue.(int8); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Int8s()
		rt := retVal.Int8s()
		switch f := fn.(type) {
		case func([]int8) int8:
			reduceLastI8(dt, rt, dimSize, def, f)
		case func(int8, int8) int8:
			genericReduceLastI8(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Int16:
		var def int16

		if def, ok = defaultValue.(int16); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Int16s()
		rt := retVal.Int16s()
		switch f := fn.(type) {
		case func([]int16) int16:
			reduceLastI16(dt, rt, dimSize, def, f)
		case func(int16, int16) int16:
			genericReduceLastI16(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Int32:
		var def int32

		if def, ok = defaultValue.(int32); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Int32s()
		rt := retVal.Int32s()
		switch f := fn.(type) {
		case func([]int32) int32:
			reduceLastI32(dt, rt, dimSize, def, f)
		case func(int32, int32) int32:
			genericReduceLastI32(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Int64:
		var def int64

		if def, ok = defaultValue.(int64); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Int64s()
		rt := retVal.Int64s()
		switch f := fn.(type) {
		case func([]int64) int64:
			reduceLastI64(dt, rt, dimSize, def, f)
		case func(int64, int64) int64:
			genericReduceLastI64(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Uint:
		var def uint

		if def, ok = defaultValue.(uint); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Uints()
		rt := retVal.Uints()
		switch f := fn.(type) {
		case func([]uint) uint:
			reduceLastU(dt, rt, dimSize, def, f)
		case func(uint, uint) uint:
			genericReduceLastU(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Uint8:
		var def uint8

		if def, ok = defaultValue.(uint8); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Uint8s()
		rt := retVal.Uint8s()
		switch f := fn.(type) {
		case func([]uint8) uint8:
			reduceLastU8(dt, rt, dimSize, def, f)
		case func(uint8, uint8) uint8:
			genericReduceLastU8(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Uint16:
		var def uint16

		if def, ok = defaultValue.(uint16); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Uint16s()
		rt := retVal.Uint16s()
		switch f := fn.(type) {
		case func([]uint16) uint16:
			reduceLastU16(dt, rt, dimSize, def, f)
		case func(uint16, uint16) uint16:
			genericReduceLastU16(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Uint32:
		var def uint32

		if def, ok = defaultValue.(uint32); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Uint32s()
		rt := retVal.Uint32s()
		switch f := fn.(type) {
		case func([]uint32) uint32:
			reduceLastU32(dt, rt, dimSize, def, f)
		case func(uint32, uint32) uint32:
			genericReduceLastU32(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Uint64:
		var def uint64

		if def, ok = defaultValue.(uint64); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Uint64s()
		rt := retVal.Uint64s()
		switch f := fn.(type) {
		case func([]uint64) uint64:
			reduceLastU64(dt, rt, dimSize, def, f)
		case func(uint64, uint64) uint64:
			genericReduceLastU64(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Uintptr:
		var def uintptr

		if def, ok = defaultValue.(uintptr); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Uintptrs()
		rt := retVal.Uintptrs()
		switch f := fn.(type) {
		case func([]uintptr) uintptr:
			reduceLastUintptr(dt, rt, dimSize, def, f)
		case func(uintptr, uintptr) uintptr:
			genericReduceLastUintptr(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Float32:
		var def float32

		if def, ok = defaultValue.(float32); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Float32s()
		rt := retVal.Float32s()
		switch f := fn.(type) {
		case func([]float32) float32:
			reduceLastF32(dt, rt, dimSize, def, f)
		case func(float32, float32) float32:
			genericReduceLastF32(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Float64:
		var def float64

		if def, ok = defaultValue.(float64); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Float64s()
		rt := retVal.Float64s()
		switch f := fn.(type) {
		case func([]float64) float64:
			reduceLastF64(dt, rt, dimSize, def, f)
		case func(float64, float64) float64:
			genericReduceLastF64(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Complex64:
		var def complex64

		if def, ok = defaultValue.(complex64); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Complex64s()
		rt := retVal.Complex64s()
		switch f := fn.(type) {
		case func([]complex64) complex64:
			reduceLastC64(dt, rt, dimSize, def, f)
		case func(complex64, complex64) complex64:
			genericReduceLastC64(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case Complex128:
		var def complex128

		if def, ok = defaultValue.(complex128); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Complex128s()
		rt := retVal.Complex128s()
		switch f := fn.(type) {
		case func([]complex128) complex128:
			reduceLastC128(dt, rt, dimSize, def, f)
		case func(complex128, complex128) complex128:
			genericReduceLastC128(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case String:
		var def string

		if def, ok = defaultValue.(string); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.Strings()
		rt := retVal.Strings()
		switch f := fn.(type) {
		case func([]string) string:
			reduceLastStr(dt, rt, dimSize, def, f)
		case func(string, string) string:
			genericReduceLastStr(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	case UnsafePointer:
		var def unsafe.Pointer

		if def, ok = defaultValue.(unsafe.Pointer); !ok {
			return errors.Errorf(defaultValueErrMsg, def, defaultValue, defaultValue)
		}
		dt := data.UnsafePointers()
		rt := retVal.UnsafePointers()
		switch f := fn.(type) {
		case func([]unsafe.Pointer) unsafe.Pointer:
			reduceLastUnsafePointer(dt, rt, dimSize, def, f)
		case func(unsafe.Pointer, unsafe.Pointer) unsafe.Pointer:
			genericReduceLastUnsafePointer(dt, rt, dimSize, def, f)
		default:
			return errors.Errorf(reductionErrMsg, fn)
		}
		return nil
	default:
		return errors.Errorf("Unsupported type %v for ReduceLast", t)
	}
}
func (e E) Map(t reflect.Type, fn interface{}, a *storage.Header, incr bool) (err error) {
	as := isScalar(a, t)
	switch t {
	case Bool:
		var f0 func(bool) bool
		var f1 func(bool) (bool, error)

		switch f := fn.(type) {
		case func(bool) bool:
			f0 = f
		case func(bool) (bool, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Bools()
		if incr {
			return errors.Errorf("Cannot perform increment on t of %v", t)
		}
		switch {
		case as && f0 != nil:
			at[0] = f0(at[0])
		case as && f0 == nil:
			at[0], err = f1(at[0])
		case !as && f0 == nil:
			err = MapErrB(f1, at)
		default:
			MapB(f0, at)
		}
	case Int:
		var f0 func(int) int
		var f1 func(int) (int, error)

		switch f := fn.(type) {
		case func(int) int:
			f0 = f
		case func(int) (int, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Ints()
		switch {
		case as && incr && f0 != nil:
			at[0] += f0(at[0])
		case as && incr && f0 == nil:
			var tmp int
			if tmp, err = f1(at[0]); err != nil {
				return
			}
			at[0] += tmp
		case as && !incr && f0 != nil:
			at[0] = f0(at[0])
		case as && !incr && f0 == nil:
			at[0], err = f1(at[0])
		case !as && incr && f0 != nil:
			MapIncrI(f0, at)
		case !as && incr && f0 == nil:
			err = MapIncrErrI(f1, at)
		case !as && !incr && f0 == nil:
			err = MapErrI(f1, at)
		default:
			MapI(f0, at)
		}
	case Int8:
		var f0 func(int8) int8
		var f1 func(int8) (int8, error)

		switch f := fn.(type) {
		case func(int8) int8:
			f0 = f
		case func(int8) (int8, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Int8s()
		switch {
		case as && incr && f0 != nil:
			at[0] += f0(at[0])
		case as && incr && f0 == nil:
			var tmp int8
			if tmp, err = f1(at[0]); err != nil {
				return
			}
			at[0] += tmp
		case as && !incr && f0 != nil:
			at[0] = f0(at[0])
		case as && !incr && f0 == nil:
			at[0], err = f1(at[0])
		case !as && incr && f0 != nil:
			MapIncrI8(f0, at)
		case !as && incr && f0 == nil:
			err = MapIncrErrI8(f1, at)
		case !as && !incr && f0 == nil:
			err = MapErrI8(f1, at)
		default:
			MapI8(f0, at)
		}
	case Int16:
		var f0 func(int16) int16
		var f1 func(int16) (int16, error)

		switch f := fn.(type) {
		case func(int16) int16:
			f0 = f
		case func(int16) (int16, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Int16s()
		switch {
		case as && incr && f0 != nil:
			at[0] += f0(at[0])
		case as && incr && f0 == nil:
			var tmp int16
			if tmp, err = f1(at[0]); err != nil {
				return
			}
			at[0] += tmp
		case as && !incr && f0 != nil:
			at[0] = f0(at[0])
		case as && !incr && f0 == nil:
			at[0], err = f1(at[0])
		case !as && incr && f0 != nil:
			MapIncrI16(f0, at)
		case !as && incr && f0 == nil:
			err = MapIncrErrI16(f1, at)
		case !as && !incr && f0 == nil:
			err = MapErrI16(f1, at)
		default:
			MapI16(f0, at)
		}
	case Int32:
		var f0 func(int32) int32
		var f1 func(int32) (int32, error)

		switch f := fn.(type) {
		case func(int32) int32:
			f0 = f
		case func(int32) (int32, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Int32s()
		switch {
		case as && incr && f0 != nil:
			at[0] += f0(at[0])
		case as && incr && f0 == nil:
			var tmp int32
			if tmp, err = f1(at[0]); err != nil {
				return
			}
			at[0] += tmp
		case as && !incr && f0 != nil:
			at[0] = f0(at[0])
		case as && !incr && f0 == nil:
			at[0], err = f1(at[0])
		case !as && incr && f0 != nil:
			MapIncrI32(f0, at)
		case !as && incr && f0 == nil:
			err = MapIncrErrI32(f1, at)
		case !as && !incr && f0 == nil:
			err = MapErrI32(f1, at)
		default:
			MapI32(f0, at)
		}
	case Int64:
		var f0 func(int64) int64
		var f1 func(int64) (int64, error)

		switch f := fn.(type) {
		case func(int64) int64:
			f0 = f
		case func(int64) (int64, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Int64s()
		switch {
		case as && incr && f0 != nil:
			at[0] += f0(at[0])
		case as && incr && f0 == nil:
			var tmp int64
			if tmp, err = f1(at[0]); err != nil {
				return
			}
			at[0] += tmp
		case as && !incr && f0 != nil:
			at[0] = f0(at[0])
		case as && !incr && f0 == nil:
			at[0], err = f1(at[0])
		case !as && incr && f0 != nil:
			MapIncrI64(f0, at)
		case !as && incr && f0 == nil:
			err = MapIncrErrI64(f1, at)
		case !as && !incr && f0 == nil:
			err = MapErrI64(f1, at)
		default:
			MapI64(f0, at)
		}
	case Uint:
		var f0 func(uint) uint
		var f1 func(uint) (uint, error)

		switch f := fn.(type) {
		case func(uint) uint:
			f0 = f
		case func(uint) (uint, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Uints()
		switch {
		case as && incr && f0 != nil:
			at[0] += f0(at[0])
		case as && incr && f0 == nil:
			var tmp uint
			if tmp, err = f1(at[0]); err != nil {
				return
			}
			at[0] += tmp
		case as && !incr && f0 != nil:
			at[0] = f0(at[0])
		case as && !incr && f0 == nil:
			at[0], err = f1(at[0])
		case !as && incr && f0 != nil:
			MapIncrU(f0, at)
		case !as && incr && f0 == nil:
			err = MapIncrErrU(f1, at)
		case !as && !incr && f0 == nil:
			err = MapErrU(f1, at)
		default:
			MapU(f0, at)
		}
	case Uint8:
		var f0 func(uint8) uint8
		var f1 func(uint8) (uint8, error)

		switch f := fn.(type) {
		case func(uint8) uint8:
			f0 = f
		case func(uint8) (uint8, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Uint8s()
		switch {
		case as && incr && f0 != nil:
			at[0] += f0(at[0])
		case as && incr && f0 == nil:
			var tmp uint8
			if tmp, err = f1(at[0]); err != nil {
				return
			}
			at[0] += tmp
		case as && !incr && f0 != nil:
			at[0] = f0(at[0])
		case as && !incr && f0 == nil:
			at[0], err = f1(at[0])
		case !as && incr && f0 != nil:
			MapIncrU8(f0, at)
		case !as && incr && f0 == nil:
			err = MapIncrErrU8(f1, at)
		case !as && !incr && f0 == nil:
			err = MapErrU8(f1, at)
		default:
			MapU8(f0, at)
		}
	case Uint16:
		var f0 func(uint16) uint16
		var f1 func(uint16) (uint16, error)

		switch f := fn.(type) {
		case func(uint16) uint16:
			f0 = f
		case func(uint16) (uint16, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Uint16s()
		switch {
		case as && incr && f0 != nil:
			at[0] += f0(at[0])
		case as && incr && f0 == nil:
			var tmp uint16
			if tmp, err = f1(at[0]); err != nil {
				return
			}
			at[0] += tmp
		case as && !incr && f0 != nil:
			at[0] = f0(at[0])
		case as && !incr && f0 == nil:
			at[0], err = f1(at[0])
		case !as && incr && f0 != nil:
			MapIncrU16(f0, at)
		case !as && incr && f0 == nil:
			err = MapIncrErrU16(f1, at)
		case !as && !incr && f0 == nil:
			err = MapErrU16(f1, at)
		default:
			MapU16(f0, at)
		}
	case Uint32:
		var f0 func(uint32) uint32
		var f1 func(uint32) (uint32, error)

		switch f := fn.(type) {
		case func(uint32) uint32:
			f0 = f
		case func(uint32) (uint32, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Uint32s()
		switch {
		case as && incr && f0 != nil:
			at[0] += f0(at[0])
		case as && incr && f0 == nil:
			var tmp uint32
			if tmp, err = f1(at[0]); err != nil {
				return
			}
			at[0] += tmp
		case as && !incr && f0 != nil:
			at[0] = f0(at[0])
		case as && !incr && f0 == nil:
			at[0], err = f1(at[0])
		case !as && incr && f0 != nil:
			MapIncrU32(f0, at)
		case !as && incr && f0 == nil:
			err = MapIncrErrU32(f1, at)
		case !as && !incr && f0 == nil:
			err = MapErrU32(f1, at)
		default:
			MapU32(f0, at)
		}
	case Uint64:
		var f0 func(uint64) uint64
		var f1 func(uint64) (uint64, error)

		switch f := fn.(type) {
		case func(uint64) uint64:
			f0 = f
		case func(uint64) (uint64, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Uint64s()
		switch {
		case as && incr && f0 != nil:
			at[0] += f0(at[0])
		case as && incr && f0 == nil:
			var tmp uint64
			if tmp, err = f1(at[0]); err != nil {
				return
			}
			at[0] += tmp
		case as && !incr && f0 != nil:
			at[0] = f0(at[0])
		case as && !incr && f0 == nil:
			at[0], err = f1(at[0])
		case !as && incr && f0 != nil:
			MapIncrU64(f0, at)
		case !as && incr && f0 == nil:
			err = MapIncrErrU64(f1, at)
		case !as && !incr && f0 == nil:
			err = MapErrU64(f1, at)
		default:
			MapU64(f0, at)
		}
	case Uintptr:
		var f0 func(uintptr) uintptr
		var f1 func(uintptr) (uintptr, error)

		switch f := fn.(type) {
		case func(uintptr) uintptr:
			f0 = f
		case func(uintptr) (uintptr, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Uintptrs()
		if incr {
			return errors.Errorf("Cannot perform increment on t of %v", t)
		}
		switch {
		case as && f0 != nil:
			at[0] = f0(at[0])
		case as && f0 == nil:
			at[0], err = f1(at[0])
		case !as && f0 == nil:
			err = MapErrUintptr(f1, at)
		default:
			MapUintptr(f0, at)
		}
	case Float32:
		var f0 func(float32) float32
		var f1 func(float32) (float32, error)

		switch f := fn.(type) {
		case func(float32) float32:
			f0 = f
		case func(float32) (float32, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Float32s()
		switch {
		case as && incr && f0 != nil:
			at[0] += f0(at[0])
		case as && incr && f0 == nil:
			var tmp float32
			if tmp, err = f1(at[0]); err != nil {
				return
			}
			at[0] += tmp
		case as && !incr && f0 != nil:
			at[0] = f0(at[0])
		case as && !incr && f0 == nil:
			at[0], err = f1(at[0])
		case !as && incr && f0 != nil:
			MapIncrF32(f0, at)
		case !as && incr && f0 == nil:
			err = MapIncrErrF32(f1, at)
		case !as && !incr && f0 == nil:
			err = MapErrF32(f1, at)
		default:
			MapF32(f0, at)
		}
	case Float64:
		var f0 func(float64) float64
		var f1 func(float64) (float64, error)

		switch f := fn.(type) {
		case func(float64) float64:
			f0 = f
		case func(float64) (float64, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Float64s()
		switch {
		case as && incr && f0 != nil:
			at[0] += f0(at[0])
		case as && incr && f0 == nil:
			var tmp float64
			if tmp, err = f1(at[0]); err != nil {
				return
			}
			at[0] += tmp
		case as && !incr && f0 != nil:
			at[0] = f0(at[0])
		case as && !incr && f0 == nil:
			at[0], err = f1(at[0])
		case !as && incr && f0 != nil:
			MapIncrF64(f0, at)
		case !as && incr && f0 == nil:
			err = MapIncrErrF64(f1, at)
		case !as && !incr && f0 == nil:
			err = MapErrF64(f1, at)
		default:
			MapF64(f0, at)
		}
	case Complex64:
		var f0 func(complex64) complex64
		var f1 func(complex64) (complex64, error)

		switch f := fn.(type) {
		case func(complex64) complex64:
			f0 = f
		case func(complex64) (complex64, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Complex64s()
		switch {
		case as && incr && f0 != nil:
			at[0] += f0(at[0])
		case as && incr && f0 == nil:
			var tmp complex64
			if tmp, err = f1(at[0]); err != nil {
				return
			}
			at[0] += tmp
		case as && !incr && f0 != nil:
			at[0] = f0(at[0])
		case as && !incr && f0 == nil:
			at[0], err = f1(at[0])
		case !as && incr && f0 != nil:
			MapIncrC64(f0, at)
		case !as && incr && f0 == nil:
			err = MapIncrErrC64(f1, at)
		case !as && !incr && f0 == nil:
			err = MapErrC64(f1, at)
		default:
			MapC64(f0, at)
		}
	case Complex128:
		var f0 func(complex128) complex128
		var f1 func(complex128) (complex128, error)

		switch f := fn.(type) {
		case func(complex128) complex128:
			f0 = f
		case func(complex128) (complex128, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Complex128s()
		switch {
		case as && incr && f0 != nil:
			at[0] += f0(at[0])
		case as && incr && f0 == nil:
			var tmp complex128
			if tmp, err = f1(at[0]); err != nil {
				return
			}
			at[0] += tmp
		case as && !incr && f0 != nil:
			at[0] = f0(at[0])
		case as && !incr && f0 == nil:
			at[0], err = f1(at[0])
		case !as && incr && f0 != nil:
			MapIncrC128(f0, at)
		case !as && incr && f0 == nil:
			err = MapIncrErrC128(f1, at)
		case !as && !incr && f0 == nil:
			err = MapErrC128(f1, at)
		default:
			MapC128(f0, at)
		}
	case String:
		var f0 func(string) string
		var f1 func(string) (string, error)

		switch f := fn.(type) {
		case func(string) string:
			f0 = f
		case func(string) (string, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.Strings()
		switch {
		case as && incr && f0 != nil:
			at[0] += f0(at[0])
		case as && incr && f0 == nil:
			var tmp string
			if tmp, err = f1(at[0]); err != nil {
				return
			}
			at[0] += tmp
		case as && !incr && f0 != nil:
			at[0] = f0(at[0])
		case as && !incr && f0 == nil:
			at[0], err = f1(at[0])
		case !as && incr && f0 != nil:
			MapIncrStr(f0, at)
		case !as && incr && f0 == nil:
			err = MapIncrErrStr(f1, at)
		case !as && !incr && f0 == nil:
			err = MapErrStr(f1, at)
		default:
			MapStr(f0, at)
		}
	case UnsafePointer:
		var f0 func(unsafe.Pointer) unsafe.Pointer
		var f1 func(unsafe.Pointer) (unsafe.Pointer, error)

		switch f := fn.(type) {
		case func(unsafe.Pointer) unsafe.Pointer:
			f0 = f
		case func(unsafe.Pointer) (unsafe.Pointer, error):
			f1 = f
		default:
			return errors.Errorf("Cannot map fn of %T to array", fn)
		}

		at := a.UnsafePointers()
		if incr {
			return errors.Errorf("Cannot perform increment on t of %v", t)
		}
		switch {
		case as && f0 != nil:
			at[0] = f0(at[0])
		case as && f0 == nil:
			at[0], err = f1(at[0])
		case !as && f0 == nil:
			err = MapErrUnsafePointer(f1, at)
		default:
			MapUnsafePointer(f0, at)
		}
	default:
		return errors.Errorf("Cannot map t of %v", t)

	}

	return
}
