package resolves

import (
	"strings"

	. "github.com/wanliu/flow/builtin/luis"
)

type AddressResolve struct {
	// Address string
	Parent *OpenOrderResolve
}

func (ar AddressResolve) Hint() string {
	return "请告诉我送货地址"
}

func (pr AddressResolve) Solve(luis ResultParams) (bool, string, string) {
	// pr.Address = "some where"
	if luis.TopScoringIntent.Intent == "地址" {
		address := strings.Trim(luis.Entities[0].Entity, " ")
		pr.Parent.Address = address

		return true, "已经定好了送货地址:" + address, "err"
	} else {
		entity, exist := FetchEntity("地址", luis.Entities)

		if exist {
			pr.Parent.Address = strings.Trim(entity.Entity, " ")
			return true, "已经定好了送货地址:" + pr.Parent.Address, "err"
		} else {
			return false, "", "无效的输入:\"" + luis.Query + "\"\n" + pr.Hint()
		}
	}
}
