package objpool

import (
	"reflect"
	"sync"

	"github.com/youngbloood/zero"
)

type pool struct {
	obj map[string]*sync.Pool
}

func New(size int) *pool {
	if size < 0 {
		size = 0
	}
	return &pool{
		obj: make(map[string]*sync.Pool, size),
	}
}

func (p *pool) Set(v interface{}, new func() interface{}) (objName string) {
	rt := reflect.TypeOf(v)
	objName = rt.String()
	p.obj[objName] = &sync.Pool{
		New: new,
	}
	return objName
}

func (p *pool) Put(v interface{}) {
	rt := reflect.TypeOf(v)
	sp := p.obj[rt.String()]
	if sp == nil {
		return
	}
	sp.Put(v)
}

// objName=${packageName}.${typeName}
// eg:packageName=pool_test,typeName=Struct;objName="pool_test.Struct" or "*pool_test.Struct"
func (p *pool) Get(objName string, isSetZero bool) interface{} {
	sp := p.obj[objName]
	if sp == nil {
		return nil
	}
	obj := sp.Get()
	if isSetZero {
		zero.Reset(obj, nil)
	}
	return obj
}
