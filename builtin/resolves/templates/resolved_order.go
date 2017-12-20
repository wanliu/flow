package templates

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"

	"github.com/wanliu/brain_data/database"
)

var solvedOrderTemplates = []string{
	`订单已经生成, 共{{CnNum (len .OrderItems)}}种产品
{{range .OrderItems}}{{.ProductName}} {{.Quantity}}{{.Unit}}{{end}}
{{if gt (len .GiftItems) 0}}
申请的赠品:
{{range .GiftItems}}
{{.ProductName}} {{.Quantity}}{{.Unit}}
{{end}}
{{end}}
时间:{{.DeliveryTime.Format "2006年01月02日"}}
客户:{{.Customer.Name}}
订单号为:{{.No}}
订单入口: http://jiejie.wanliu.biz/order/QueryDetail/{{.GlobelId}}`,
}

func RenderSolvedOrder(order database.Order) string {
	cus := order.GetCustomer()
	if cus != nil {
		order.Customer = *cus
	}

	tmpl, err := template.New("solvedOrder").Funcs(template.FuncMap{
		"CnNum": CnNum,
	}).Parse(solvedOrderTemplates[0])

	if err != nil {
		return err.Error()
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, order)
	if err != nil {
		return err.Error()
	}

	result := tpl.String()
	return strings.Replace(result, "\n\n", "\n", -1)
}

func CnNum(num int) string {
	switch num {
	case 1:
		return "一"
	case 2:
		return "两"
	case 3:
		return "三"
	case 4:
		return "四"
	case 5:
		return "五"
	case 6:
		return "六"
	case 7:
		return "七"
	case 8:
		return "八"
	case 9:
		return "九"
	case 10:
		return "十"
	default:
		return strconv.Itoa(num)
	}
}
