package builitn

import "github.com/wanliu/components"

var _components = make(map[string]func() interface{})

var Info = components.Package{
	Name:        "wanliu-components",
	Description: "wanliu flow builtin components",
	Version:     "0.0.1",
}

func ComponentList() map[string]func() interface{} {
	_components["dom/GetElement"] = NewGetElement
	_components["Split"] = NewSplit
	_components["Output"] = NewOutput
	_components["ReadFile"] = NewReadFile
	_components["ReadLine"] = NewReadLine
	_components["LuisAnalyze"] = NewLuisAnalyze
	_components["Stringifier"] = NewStringifier
	_components["StringifierB"] = NewStringifierB
	_components["TextReader"] = NewTextReader

	return _components
}
