package memutil

import (
	"testing"
)

type TestType struct {
	TestField  int
	TestString string
}

var TestTypeProto TestType

func TestNativeNewAndDelete(t *testing.T) {
	var list []*TestType
	for i := 0; i < 100; i++ {
		testType := NativeNew(TestTypeProto).(*TestType)
		testType.TestField = 3
		testType.TestString = "sdfsdf"
		t.Log(testType.TestField, testType.TestString)
		list = append(list, testType)
	}
	for i := 0; i < 100; i++ {
		NativeDelete(list[i])
	}
}

func TestNativeBuffer(t *testing.T) {
	buffer := NativeAllocateBuffer(1024)
	t.Log(len(buffer))
	buffer[1023] = 30
	t.Log(buffer[1023])
	t.Log("ptr:", buffer.Ptr())
	t.Log("string:", buffer.String())
	NativeFreeBuffer(buffer)
	t.Log(buffer[1023])
}
