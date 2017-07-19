package resolves

import (
	"strings"

	. "github.com/wanliu/flow/builtin/luis"
)

type AddressResolve struct {
	// Address string
	Parent *OpenOrderResolve
	// Confirm *Confirm
	Confirm string
}

func (ar AddressResolve) Hint() string {
	return "请告诉我送货地址"
}

func (pr *AddressResolve) Solve(luis ResultParams) (bool, string, string) {

	if luis.TopScoringIntent.Intent == "地址" {
		var address string

		input, exist := FetchEntity("地址", luis.Entities)

		if exist {
			address = input.Entity
		} else {
			address = luis.Query
		}

		address = strings.Trim(address, " ")
		pr.Parent.Customer = address

		return true, "已经定好了送货地址:" + address, ""
	} else if luis.TopScoringIntent.Intent == "确认" && pr.Confirm != "" {
		pr.Parent.Customer, pr.Confirm = pr.Confirm, ""

		return true, "已经定好了送货地址:" + pr.Parent.Customer, ""
	} else if luis.TopScoringIntent.Intent == "取消" && pr.Confirm != "" {
		pr.Confirm = ""

		return false, "", pr.Hint()
	} else {
		entity, exist := FetchEntity("地址", luis.Entities)

		if exist {
			pr.Parent.Customer = strings.Trim(entity.Entity, " ")
			return true, "已经定好了送货地址:" + pr.Parent.Customer, ""
		} else {
			address := strings.Trim(luis.Query, " ")
			pr.Confirm = address
			return false, "", "收到您的回复：\"" + address + "\"\n 是否将\"" + address + "\"做为收货地址？"
		}
	}
}
