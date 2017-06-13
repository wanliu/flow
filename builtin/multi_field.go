package builtin

import (
	"reflect"
	"sync"
)

type ProcessHandle func() error

type MultiField struct {
	sync.Once
	cond    sync.Cond
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
	mf.Do(func() {
		mf.ms = &MultiFieldSet{
			Fields: mf.Fields,
			Values: make([]MultiFieldValue, 0),
		}

		if mf.cond.L == nil {
			mf.cond.L = &sync.Mutex{}
		}
		mf.cond.L.Lock()

		for !mf.ms.Valid() {
			mf.cond.Wait()
		}
		if mf.Process != nil {
			mf.Process()
		}
		// mf.Reset()
		mf.ms = nil
		mf.cond.L.Unlock()

	})
}

func (mf *MultiField) SetValue(name string, val interface{}) {
	mf.Run()
	mf.ms.SetValue(name, val)
}

func (mf *MultiField) Value(name string) interface{} {
	if mf.ms != nil {
		return mf.Value(name)
	}
	return nil
}

func (ms *MultiFieldSet) Valid() bool {
	for _, field := range ms.Fields {
		var match bool
		for _, val := range ms.Values {
			v := reflect.ValueOf(val.Value)
			if val.Name == field && isZero(v) {
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
			z = z && isZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
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
