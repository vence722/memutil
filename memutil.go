// memutil v0.2
// @description A util library for allocating native memory
// @authors     Vence Lam(vence722@gmail.com)
package memutil

//#include<stdlib.h>
import "C"

import (
	"github.com/pkg/errors"
	"reflect"
	"unsafe"
)

/**
 * Allocate native memory for a specific size which is derived by the pointer parameter.
 * Input parameter `ptr` must be a REFERENCE to a valid pointer which can be point to nil.
 * The input pointer will point to a valid object which is located OUTSIDE Golang heap.
 * When finish using this pointer, you MUST call memutil.NativeDelete() function to
 * release the memory manually.
 *
 * Example usage:
 *     type MyStruct struct {
 *         name string
 *     }
 *
 *     var m *MyStruct // declare an empty pointer
 *     memutil.NativeNew(&m) // should pass the REFERENCE of the pointer
 *     m.name = "Vence" // do something with m as usual
 *
 *     ......
 *
 *     memutil.NativeDelete(&m) // destroy the struct after using it
 */
func NativeNew(ptr interface{}) error {
	ppType := reflect.ValueOf(ptr).Type()
	pType := ppType.Elem()
	sType := pType.Elem()
	if ppType.Kind() != reflect.Ptr || pType.Kind() != reflect.Ptr {
		return errors.New("input parameter should be a pointer to pointer")
	}
	p := C.malloc(C.size_t(sType.Size()))
	if p == nil {
		return errors.New("failed to call C.malloc()")
	}
	reflect.ValueOf(ptr).Elem().Set(reflect.NewAt(sType, p))
	return nil
}

/**
 * Release the memory which allocated by the memutil.NativeNew() function.
 * Input parameter `ptr` must be a REFERENCE to a valid pointer which can be point to nil.
 */
func NativeDelete(ptr interface{}) {
	ppVal := reflect.ValueOf(ptr)
	pVal := ppVal.Elem()
	C.free(unsafe.Pointer(pVal.Pointer()))
}

type NativeBuffer []byte

func (this *NativeBuffer) init(addr unsafe.Pointer, size int) {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(this))
	sh.Data = uintptr(addr)
	sh.Len = size
	sh.Cap = size
}

func (this *NativeBuffer) Ptr() uintptr {
	ptrH := (*reflect.SliceHeader)(unsafe.Pointer(this))
	return ptrH.Data
}

func (this *NativeBuffer) Len() int {
	return len(*this)
}

func (this *NativeBuffer) String() string {
	strH := &reflect.StringHeader{
		Data: this.Ptr(),
		Len:  len(*this),
	}
	return *(*string)(unsafe.Pointer(strH))
}

func (this *NativeBuffer) Write(data []byte) int {
	written := 0
	for i := range data {
		if i >= len(*this) {
			break
		}
		(*this)[i] = data[i]
		written++
	}
	return written
}

func NativeAllocateBuffer(size int) (NativeBuffer, error) {
	buffer := NativeBuffer{}
	addr := C.malloc(C.size_t(size))
	if addr == nil {
		return nil, errors.New("failed to call C.malloc()")
	}
	buffer.init(addr, size)
	return buffer, nil
}

func NativeFreeBuffer(buffer NativeBuffer) {
	C.free(unsafe.Pointer(buffer.Ptr()))
}
