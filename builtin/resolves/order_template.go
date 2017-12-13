package resolves

import (
	"bytes"
	"text/template"

	"github.com/wanliu/brain_data/database"
)

var solvedOrderTemplates = []string{
	`订单已经生成, 共{{len(.OrderItems)}}种产品
	{{range .OrderItems}}
		{{.ProductName}} {{.Quantity}}{{.Unit}}
	{{end}}
	{{if ge len(.GiftsItems) 0}}
		申请的赠品:
	 	{{range .GiftsItems}}
	 		{{.ProductName}} {{.Quantity}}{{.Unit}}
	 	{{end}}
	{{end}}
	时间:{{.DeliveryTime.Format "2006年01月02日"}}
	客户:{{.Customer}}
	订单已经生成，订单号为：{{.No}}
	订单入口: http://jiejie.wanliu.biz/order/QueryDetail/{{.GlobelId}}`,
}

func RenderSolvedOrder(order database.Order) string {
	tmpl, _ := template.New("solvedOrder").Parse(solvedOrderTemplates[0])
	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, order)
	if err != nil {
		return err.Error()
	}

	return tpl.String()
}
