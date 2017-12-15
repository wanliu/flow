package builtin

import (
	"reflect"
	"sync"
	"time"
)

type ProcessHandle func() error

type MultiField struct {
	cond    *sync.Cond
	once    *sync.Once
	Fields  []string
	Process ProcessHandle
	ms      *MultiFieldSet
}

type MultiFieldSet struct {
	Fields []string
	Values []MultiFieldValue
}

type MultiFieldValue struct {
	Name  string
	Value interface{}
}

func (mf *MultiField) Run() {
	if mf.ms == nil {
		mf.ms = &MultiFieldSet{
			Fields: mf.Fields,
			Values: make([]MultiFieldValue, 0),
		}
	}

	if mf.once == nil {
		mf.once = new(sync.Once)
	}

	mf.once.Do(func() {
		mf.cond = sync.NewCond(&sync.Mutex{})
		mf.cond.L.Lock()

		for !mf.ms.Valid() {
			mf.cond.Wait()
		}

		if mf.Process != nil {
			mf.Process()
		}

		mf.cond.L.Unlock()
		mf.Reset()
	})

}

func (mf *MultiField) SetValue(name string, val interface{}) {
	go mf.Run()
	time.Sleep(1 * time.Millisecond)
	mf.ms.SetValue(name, val)
	mf.cond.Signal()
}

func (mf *MultiField) Value(name string) interface{} {
	if mf.ms != nil {
		return mf.ms.Value(name)
	}
	return nil
}

func (mf *MultiField) Reset() {
	mf.ms = &MultiFieldSet{
		Fields: mf.Fields,
		Values: make([]MultiFieldValue, 0),
	}

	mf.once = nil
}

func (ms *MultiFieldSet) Valid() bool {
	for _, field := range ms.Fields {
		var match bool
		for _, val := range ms.Values {
			if val.Name == field {
				match = true
				break
			}
		}

		if !match {
			return false
		}
	}
	return true
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).CanSet() {
				z = z && isZero(v.Field(i))
			}
		}
		return z
	case reflect.Ptr:
		return isZero(reflect.Indirect(v))
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	result := v.Interface() == z.Interface()

	return result
}

func (ms *MultiFieldSet) SetValue(name string, val interface{}) {
	ms.Values = append(ms.Values, MultiFieldValue{Name: name, Value: val})
}

func (ms *MultiFieldSet) Value(name string) interface{} {
	for _, val := range ms.Values {
		if val.Name == name {
			return val.Value
		}
	}

	return nil
}
