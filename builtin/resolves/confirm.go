package resolves

import (
	. "github.com/wanliu/flow/builtin/luis"
)

type Confirm struct {
	HintTxt string
	Cfm     bool
}

func (r Confirm) Hint() string {
	return r.HintTxt
}

func (r *Confirm) Solve(luis ResultParams) (bool, string, string) {
	// r.Address = "some where"
	if luis.TopScoringIntent.Intent == "取消" {
		r.Cfm = false

		return true, "已经取消", ""
	} else if luis.TopScoringIntent.Intent == "确认" {
		r.Cfm = true

		return true, "确认成功", ""
	} else {
		return false, "", "无效的输入:\"" + luis.Query + "\"\n" + r.Hint()
	}
}
