package resolves

import (
	"bytes"
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
