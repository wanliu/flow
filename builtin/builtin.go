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
	_components["MyInput"] = NewMyInput
	_components["QuerySave"] = NewQuerySave
	_components["Final"] = NewFinal
	_components["Order"] = GetOrder
	_components["NewOrder"] = NewNewOrder
	_components["PatchOrder"] = NewPatchOrder
	_components["Unimplemented"] = NewUnimplemented
	_components["Critical"] = NewCritical
	_components["Abuse"] = NewAbuse
	_components["ApiAi"] = NewApiAi
	_components["OperationNotice"] = NewOperationNotice
	_components["OrderPrinter"] = NewOrderPrinter
	_components["OrderCancel"] = NewOrderCancel
	_components["Voice"] = NewVoice
	_components["OrderAddress"] = NewOrderAddress
	_components["OrderCustomer"] = NewOrderCustomer
	_components["Confirm"] = NewConfirm
	_components["OrderTimeout"] = NewOrderTimeout
	_components["OrderTouch"] = NewOrderTouch
	_components["OrderDelete"] = NewOrderDelete
	_components["OrderNumber"] = NewOrderNumber
	_components["TextPreprocesor"] = NewTextPreprocesor
	_components["CustomerOrders"] = NewCustomerOrders
	_components["OrderItemDelete"] = NewOrderItemDelete
	_components["OrderGiftDelete"] = NewOrderGiftDelete

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
		Constructor: GetOrder,
	}

	result["NewOrder"] = Component{
		Name:        "NewOrder",
		Description: "New NewOrder Component",
		Icon:        "",
		Constructor: NewNewOrder,
	}

	result["PatchOrder"] = Component{
		Name:        "PatchOrder",
		Description: "New PatchOrder Component",
		Icon:        "",
		Constructor: NewPatchOrder,
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

	result["ApiAi"] = Component{
		Name:        "ApiAi",
		Description: "ApiAi Component",
		Icon:        "",
		Constructor: NewApiAi,
	}

	result["OperationNotice"] = Component{
		Name:        "OperationNotice",
		Description: "OperationNotice Component",
		Icon:        "",
		Constructor: NewOperationNotice,
	}

	result["OrderPrinter"] = Component{
		Name:        "OrderPrinter",
		Description: "OrderPrinter Component",
		Icon:        "",
		Constructor: NewOrderPrinter,
	}

	result["OrderCancel"] = Component{
		Name:        "OrderCancel",
		Description: "OrderCancel Component",
		Icon:        "",
		Constructor: NewOrderCancel,
	}

	result["Voice"] = Component{
		Name:        "Voice",
		Description: "Voice Component",
		Icon:        "",
		Constructor: NewVoice,
	}

	result["OrderAddress"] = Component{
		Name:        "OrderAddress",
		Description: "OrderAddress Component",
		Icon:        "",
		Constructor: NewOrderAddress,
	}

	result["OrderCustomer"] = Component{
		Name:        "OrderCustomer",
		Description: "OrderCustomer Component",
		Icon:        "",
		Constructor: NewOrderCustomer,
	}

	result["Confirm"] = Component{
		Name:        "Confirm",
		Description: "Confirm Component",
		Icon:        "",
		Constructor: NewConfirm,
	}

	result["OrderTimeout"] = Component{
		Name:        "OrderTimeout",
		Description: "OrderTimeout Component",
		Icon:        "",
		Constructor: NewOrderTimeout,
	}

	result["OrderTouch"] = Component{
		Name:        "OrderTouch",
		Description: "OrderTouch Component",
		Icon:        "",
		Constructor: NewOrderTouch,
	}

	result["OrderDelete"] = Component{
		Name:        "OrderDelete",
		Description: "OrderDelete Component",
		Icon:        "",
		Constructor: NewOrderDelete,
	}

	result["OrderNumber"] = Component{
		Name:        "OrderNumber",
		Description: "OrderNumber Component",
		Icon:        "",
		Constructor: NewOrderNumber,
	}

	result["TextPreprocesor"] = Component{
		Name:        "TextPreprocesor",
		Description: "TextPreprocesor Component",
		Icon:        "",
		Constructor: NewTextPreprocesor,
	}

	result["CustomerOrders"] = Component{
		Name:        "CustomerOrders",
		Description: "CustomerOrders Component",
		Icon:        "",
		Constructor: NewCustomerOrders,
	}

	result["OrderItemDelete"] = Component{
		Name:        "OrderItemDelete",
		Description: "OrderItemDelete Component",
		Icon:        "",
		Constructor: NewOrderItemDelete,
	}

	result["OrderGiftDelete"] = Component{
		Name:        "OrderGiftDelete",
		Description: "OrderGiftDelete Component",
		Icon:        "",
		Constructor: NewOrderGiftDelete,
	}

	return result
}
