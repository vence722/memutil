package memutil

import (
	"fmt"
	"os/exec"
	"testing"
)

type TestType struct {
	TestInt  int
	TestString string
}

type TestTypeSimple struct {
	TestByte byte
	TestInt int
}

func printMemUsage() {
	cmd := exec.Command("tasklist", "/FI", "MEMUSAGE gt 200000")
	output, _ := cmd.Output()
	fmt.Println(string(output))
}

func TestNativeNewAndDelete(t *testing.T) {
	var list []*TestType

	// Allocate a bunch of TestType struct instances
	for i := 0; i < 1000000; i++ {
		var test *TestType
		err := NativeNew(&test)
		if err != nil {
			t.Error(err)
		}
		test.TestInt = 722
		test.TestString = "Vence is smart!"
		list = append(list, test)
	}

	fmt.Println("After allocate native memory:")
	printMemUsage()

	for i := 0; i < 1000000; i++ {
		NativeDelete(&list[i])
	}

	fmt.Println("After release allocated native memory:")
	printMemUsage()
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

func TestPrintPtrMemory(t *testing.T) {
	a := &TestTypeSimple{
		TestByte: 200,
		TestInt: 10000,
	}
	mem := PrintPtrMemory(a)
	t.Log("mem", fmt.Sprintf("%v", mem))
}
