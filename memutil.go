// memutil v0.1
// @description A util library for allocating native memory
// @authors     Vence Lin(vence722@gmail.com)
package memutil

//#include<stdlib.h>
import "C"

import (
	"reflect"
	"unsafe"
)

func NativeNew(proto interface{}) interface{} {
	p := C.malloc(C.size_t(reflect.TypeOf(proto).Size()))
	if p == nil {
		panic("failed to call C.malloc()")
	}
	return reflect.NewAt(reflect.TypeOf(proto), p).Interface()
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

func (this *NativeBuffer) String() string {
	strH := &reflect.StringHeader{
		Data: this.Ptr(),
		Len:  len(*this),
	}
	return *(*string)(unsafe.Pointer(strH))
}

func NativeAllocateBuffer(size int) NativeBuffer {
	buffer := NativeBuffer{}
	addr := C.malloc(C.size_t(size))
	buffer.init(addr, size)
	return buffer
}

func NativeFreeBuffer(buffer NativeBuffer) {
	C.free(unsafe.Pointer(buffer.Ptr()))
}
