// memutil v0.1
// @description A util library for allocating native memory
// @authors     Vence Lin(vence722@gmail.com)
package memutil

//#include<stdlib.h>
import "C"

import (
	"github.com/pkg/errors"
	"reflect"
	"unsafe"
)

func NativeNew(ptr interface{}) error {
	ppType := reflect.TypeOf(ptr)
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

func NativeDelete(ptr interface{}) {
	C.free(unsafe.Pointer(reflect.ValueOf(ptr).Pointer()))
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
