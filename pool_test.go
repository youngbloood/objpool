package objpool_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/youngbloood/objpool"
)

type Struct struct {
	Int int
}

func newStruct() interface{} {
	return &Struct{
		Int: 10,
	}
}

func newInt() interface{} {
	return 10
}

func TestObjpool(t *testing.T) {
	pool := objpool.New(0)
	pool.Set(&Struct{}, newStruct)
	pool.Set(8, newInt)

	stc := &Struct{
		Int: 20,
	}
	pool.Put(stc)

	fmt.Printf("ptr = %p\n", stc)

	runtime.GC()
	structVal, ok := pool.Get("Struct", false).(*Struct)
	if ok {
		fmt.Println("structVal = ", structVal)
		fmt.Printf("structVal.ptr = %p\n", structVal)
	}

	intVal, ok := pool.Get("int", false).(int)
	if ok {
		fmt.Println("intVal = ", intVal)
		fmt.Printf("intVal.ptr = %p\n", &intVal)
	}

}

func TestSyncPool(t *testing.T) {
	pool := sync.Pool{
		New: newStruct,
	}
	stc := &Struct{
		Int: 20,
	}
	fmt.Printf("ptr = %p\n", stc)

	pool.Put(stc)
	// runtime.GC()
	structVal, ok := pool.Get().(*Struct)
	if ok {
		fmt.Println("structVal = ", structVal)
		fmt.Printf("structVal.ptr = %p\n", structVal)
	}

}
