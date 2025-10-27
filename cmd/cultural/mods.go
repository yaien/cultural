package main

import (
	"fmt"
	"reflect"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs"
	"github.com/yaien/cultural/internal/modules/landing"
)

var modules = []infrastructure.Module{
	&configs.Module{},
	&landing.Module{},
}

func register(mono *infrastructure.Monolith) error {
	for _, mod := range modules {
		if err := infrastructure.Register(mono, mod); err != nil {
			name := reflect.TypeOf(mod).Elem().PkgPath()
			return fmt.Errorf("failed to register module %s: %w", name, err)
		}
	}
	return nil
}
