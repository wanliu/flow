// 开单产品选择
package resolves

import (
	"log"
	"time"

	. "github.com/wanliu/flow/builtin/luis"
)

type OrderTimeResolve struct {
	// Time   time.Time
	Parent *OpenOrderResolve
}

func (ar OrderTimeResolve) Hint() string {
	return "请告诉我送货时间"
}

func (pr OrderTimeResolve) Solve(luis ResultParams) (bool, string, string) {
	if luis.TopScoringIntent.Intent == "时间" {
		dTime := time.Now()

		for _, e := range luis.Entities {
			if e.Type == "builtin.datetime.date" {
				luisTime, err := time.Parse("2006-01-02", e.Resolution.Date)

				if err != nil {
					log.Printf("::::::::ERROR: %v", err)
				} else {
					dTime = luisTime
				}
			}
		}

		// dTime := strings.Trim(luis.Entities[0].Entity, " ")

		pr.Parent.Time = dTime
		return true, "已经定好了送货时间:" + dTime.String(), "err"
	} else {
		return false, "", "无效的输入\n" + pr.Hint()
	}
}
