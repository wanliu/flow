package flow

import (
	"fmt"
	"plugin"

	goflow "github.com/trustmaster/goflow"
	"github.com/wanliu/components"
)

type Package struct {
	Name          string
	Version       string
	Path          string
	Components    []string
	componentList map[string]func() interface{}
}

func LoadPackage(filename string) (*Package, error) {
	var pk Package

	p, err := plugin.Open(filename)
	if err != nil {
		return nil, err
	}

	v, err := p.Lookup("Info")
	if err != nil {
		return nil, err
	}

	pkg, ok := v.(components.Package)
	if !ok {
		return nil, fmt.Errorf("invalid Info Struct")
	}

	pk.Name = pkg.Name
	pk.Version = pkg.Version

	v, err = p.Lookup("ComponentList")
	if err != nil {
		return nil, err
	}

	componentList, ok := v.(func() map[string]func() interface{})
	if !ok {
		return nil, fmt.Errorf("invalid ComponentList func")
	}

	pk.componentList = componentList()

	for name, _ := range pk.componentList {
		// goflow.Register(name, constructor)
		pk.Components = append(pk.Components, name)
	}

	return &pk, nil
}

func (pk *Package) RegisterComponents() error {
	for name, constructor := range pk.componentList {
		goflow.Register(name, constructor)
		// pk.Components = append(pk.Components, name)
	}
	return nil
}
