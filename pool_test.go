package objpool_test

import (
	"runtime"
	"testing"

	"github.com/youngbloood/objpool"
)

type Struct struct {
	Int int
}

func newStruct() interface{} {
	return &Struct{
		Int: 1000,
	}
}

func newInt() interface{} {
	return 2000
}

func TestObjPoolWithoutGC(t *testing.T) {
	pool := objpool.New()
	structName := pool.Set(&Struct{}, newStruct)
	pool.Set(8, newInt)

	stc := &Struct{
		Int: 20,
	}

	pool.Put(stc)
	structVal, ok := pool.Get(structName, false).(*Struct)
	if ok {
		t.Logf("structVal.ptr.before = %p\n", stc)
		t.Logf("structVal = %+v\n", structVal)
		t.Logf("structVal.ptr.after = %p\n", structVal)
	}

	intVal, ok := pool.Get("int", false).(int)
	if ok {
		t.Logf("intVal = %d\n", intVal)
		t.Logf("intVal.ptr = %p\n", &intVal)
	}
}

func TestObjPoolWithGC(t *testing.T) {
	pool := objpool.New()
	structName := pool.Set(&Struct{}, newStruct)
	pool.Set(8, newInt)

	stc := &Struct{
		Int: 20,
	}
	pool.Put(stc)
	runtime.GC()
	structVal, ok := pool.Get(structName, false).(*Struct)
	if ok {
		t.Logf("structVal.ptr.before = %p\n", stc)
		t.Logf("structVal = %+v\n", structVal)
		t.Logf("structVal.ptr.after = %p\n", structVal)
	}

	intVal, ok := pool.Get("int", false).(int)
	if ok {
		t.Logf("intVal = %d\n", intVal)
		t.Logf("intVal.ptr = %p\n", &intVal)
	}
}

func BenchmarkObjPoolGet(b *testing.B) {
	pool := objpool.New()
	structName := pool.Set(&Struct{}, newStruct)
	pool.Put(&Struct{
		Int: 20,
	})

	for i := 0; i < b.N; i++ {
		_ = pool.Get(structName, true)
	}
}

func BenchmarkObjPoolSet(b *testing.B) {
	pool := objpool.New()
	pool.Set(&Struct{}, newStruct)
	stc := &Struct{
		Int: 20,
	}
	pool.Put(stc)

	for i := 0; i < b.N; i++ {
		pool.Put(stc)
	}
}

func BenchmarkObjPoolGetAndSet(b *testing.B) {
	pool := objpool.New()
	structName := pool.Set(&Struct{}, newStruct)
	stc := &Struct{
		Int: 20,
	}
	pool.Put(stc)

	for i := 0; i < b.N; i += 2 {
		pool.Put(stc)
	}
	for i := 1; i < b.N; i += 2 {
		_ = pool.Get(structName, true)
	}
}
