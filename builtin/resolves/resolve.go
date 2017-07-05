package resolves

import (
	. "github.com/wanliu/flow/builtin/luis"
)

type Resolve interface {
	Hint() string
	Solve(ResultParams) (bool, string, string) // 是否全部完成，完成提示，下一步动作提醒
}
