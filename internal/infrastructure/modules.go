package infrastructure

import (
	"fmt"
	"reflect"
)

var modules = make(map[string]Module)

type Module interface {
	Init(mono *Monolith) error
}

func Register(mono *Monolith, mod Module) error {
	name := reflect.TypeOf(mod).Elem().PkgPath()
	modules[name] = mod
	return mod.Init(mono)
}

func Resolve[T Module](mono *Monolith) (T, error) {
	var t T
	name := reflect.TypeOf(t).Elem().PkgPath()
	mod, ok := modules[name]
	if !ok {
		return t, fmt.Errorf("module %s not registered", name)
	}
	return mod.(T), nil
}
