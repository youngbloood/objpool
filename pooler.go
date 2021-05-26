package objpool

type Pooler interface {
	Set(v interface{}, new func() interface{}) (objName string)
	Put(v interface{})
	Get(objName string, isSetZero bool) interface{}
}
