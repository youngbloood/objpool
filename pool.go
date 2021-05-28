package objpool

import (
	"reflect"
	"sync"

	"github.com/youngbloood/zero"
)

type pool struct {
	mux sync.RWMutex
	obj *sync.Map
}

var _ Pooler = (*pool)(nil)

func New() *pool {
	return &pool{
		obj: &sync.Map{},
	}
}

func (p *pool) Set(v interface{}, new func() interface{}) (objName string) {
	rt := reflect.TypeOf(v)
	objName = rt.String()
	p.mux.RLock()
	defer p.mux.RUnlock()
	p.obj.Store(objName, &sync.Pool{
		New: new,
	})
	return objName
}

func (p *pool) Put(v interface{}) {
	rt := reflect.TypeOf(v)
	pool := p.getPool(rt.String())
	if pool == nil {
		return
	}
	pool.Put(v)
}

// objName=${packageName}.${typeName}
// eg:packageName=pool_test,typeName=Struct;objName="pool_test.Struct" or "*pool_test.Struct"
func (p *pool) Get(objName string, isSetZero bool) interface{} {
	pool := p.getPool(objName)
	if pool == nil {
		return nil
	}
	obj := pool.Get()
	if isSetZero {
		zero.Reset(obj)
	}
	return obj
}

func (p *pool) Reset() {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.obj = &sync.Map{}
}

func (p *pool) getPool(objName string) *sync.Pool {
	p.mux.RLock()
	defer p.mux.RUnlock()
	pv, exist := p.obj.Load(objName)
	if !exist {
		return nil
	}
	pool, ok := pv.(*sync.Pool)
	if !ok {
		return nil
	}
	return pool
}
