package builtin

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
	_components["IntentCheck"] = NewIntentCheck
	_components["CtxReset"] = NewCtxReset
	_components["TryGetProducts"] = NewTryGetProducts
	_components["MyInput"] = NewMyInput
	_components["QuerySave"] = NewQuerySave
	_components["ContextManager"] = NewContextManager
	_components["Final"] = NewFinal
	_components["Order"] = NewOrder
	_components["Praise"] = NewPraise
	_components["Unimplemented"] = NewUnimplemented
	_components["Critical"] = NewCritical
	_components["Abuse"] = NewAbuse
	_components["Greet"] = NewGreet
	_components["StockQuery"] = NewStockQuery
	_components["PriceQuery"] = NewPriceQuery

	return _components
}

func ComponentInfos() map[string]Component {
	var result = make(map[string]Component)
	result["dom/GetElement"] = Component{
		Name:        "dom/GetElement",
		Description: "get dom element object for demo",
		Icon:        "external-link",
		Constructor: NewGetElement,
	}

	result["Split"] = Component{
		Name:        "Split",
		Description: "Split ports",
		Icon:        "cut",
		Constructor: NewSplit,
	}

	result["Output"] = Component{
		Name:        "Output",
		Description: "Print string to terminal",
		Icon:        "sign-out",
		Constructor: NewOutput,
	}

	result["ReadFile"] = Component{
		Name:        "ReadFile",
		Description: "Read File data",
		Icon:        "file-o",
		Constructor: NewReadFile,
	}

	result["ReadLine"] = Component{
		Name:        "ReadLine",
		Description: "Read *File Stream from stdin",
		Icon:        "file-text-o",
		Constructor: NewReadLine,
	}

	result["LuisAnalyze"] = Component{
		Name:        "LuisAnalyze",
		Description: "Call LUIS Service to analyze query string",
		Icon:        "microchip",
		Constructor: NewLuisAnalyze,
	}

	result["Stringifier"] = Component{
		Name:        "Stringifier",
		Description: "stringify Luis Query response data",
		Icon:        "ellipsis-h",
		Constructor: NewStringifier,
	}

	result["IntentCheck"] = Component{
		Name:        "IntentCheck",
		Description: "Check LuisQuery result with Intent and Score",
		Icon:        "search-plus",
		Constructor: NewIntentCheck,
	}

	result["CtxReset"] = Component{
		Name:        "CtxReset",
		Description: "Reset Context to initial status",
		Icon:        "refresh",
		Constructor: NewCtxReset,
	}

	result["TryGetProducts"] = Component{
		Name:        "TryGetProducts",
		Description: "Try to query products data from Entities of Luis Query Result",
		Icon:        "product-hunt",
		Constructor: NewTryGetProducts,
	}

	result["MyInput"] = Component{
		Name:        "MyInput",
		Description: "Input with stdin stream",
		Icon:        "keyboard-o",
		Constructor: NewMyInput,
	}

	result["QuerySave"] = Component{
		Name:        "QuerySave",
		Description: "Merge Context with Value",
		Icon:        "save",
		Constructor: NewQuerySave,
	}

	result["ContextManager"] = Component{
		Name:        "ContextManager",
		Description: "Context Manager must be use in Context Component",
		Icon:        "server",
		Constructor: NewQuerySave,
	}

	result["Final"] = Component{
		Name:        "Final",
		Description: "Context Final Component",
		Icon:        "stop-circle",
		Constructor: NewFinal,
	}

	result["Order"] = Component{
		Name:        "Order",
		Description: "New Order Component",
		Icon:        "",
		Constructor: NewOrder,
	}

	result["Praise"] = Component{
		Name:        "Praise",
		Description: "Praise Component",
		Icon:        "",
		Constructor: NewPraise,
	}

	result["Unimplemented"] = Component{
		Name:        "Unimplemented",
		Description: "Unimplemented Component",
		Icon:        "",
		Constructor: NewUnimplemented,
	}

	result["Critical"] = Component{
		Name:        "Critical",
		Description: "Critical Component",
		Icon:        "",
		Constructor: NewCritical,
	}

	result["Abuse"] = Component{
		Name:        "Abuse",
		Description: "Abuse Component",
		Icon:        "",
		Constructor: NewAbuse,
	}

	result["Greet"] = Component{
		Name:        "Greet",
		Description: "Greet Component",
		Icon:        "",
		Constructor: NewGreet,
	}

	result["StockQuery"] = Component{
		Name:        "StockQuery",
		Description: "StockQuery Component",
		Icon:        "",
		Constructor: NewStockQuery,
	}

	result["PriceQuery"] = Component{
		Name:        "PriceQuery",
		Description: "PriceQuery Component",
		Icon:        "",
		Constructor: NewPriceQuery,
	}

	return result
}
