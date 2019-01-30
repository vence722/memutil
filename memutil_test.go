package memutil

import (
	"testing"
)

type TestType struct {
	TestField  int
	TestString string
}

func TestNativeNewAndDelete(t *testing.T) {
	var list []*TestType
	for i := 0; i < 100; i++ {
		var test *TestType
		err := NativeNew(&test)
		if err != nil {
			t.Error(err)
		}
		test.TestField = 3
		test.TestString = "sdfsdf"
		t.Log(test.TestField, test.TestString)
		list = append(list, test)
	}
	for i := 0; i < 100; i++ {
		NativeDelete(list[i])
	}
}

func TestNativeBuffer(t *testing.T) {
	buffer, err := NativeAllocateBuffer(1024)
	if err != nil {
		t.Error(err)
	}
	t.Log(len(buffer))
	buffer[1023] = 30
	t.Log(buffer[1023])
	t.Log("ptr:", buffer.Ptr())
	t.Log("string:", buffer.String())
	NativeFreeBuffer(buffer)
	t.Log(buffer[1023])
}
