package templates

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/wanliu/brain_data/database"
)

var cusOrdersHeaderTpl = `
{{if eq .Duration ""}}
	{{if ne .ProductName ""}}
查询客户\"{{.CustomerName}}\"{{.Duration}}包含\"{{.ProductName}}\"的订单
	{{else}}
查询客户\"{{.CustomerName}}\"{{.Duration}}的订单
	{{end}}
{{else if .QuertyTime.IsZero}}
	{{if ne r.ProductName ""}}
查询客户\"{{.CustomerName}}\"{{.QuertyTime.Format "2006年1月2日"}}包含\"{{.ProductName}}\"的订单
	{{else}}
查询客户\"{{.CustomerName}}\"{{.QuertyTime.Format "2006年1月2日"}}的订单
	{{end}}
{{else}}
	{{if eq r.ProductName ""}}
查询客户\"{{.CustomerName}}\"最近包含\"{{.ProductName}}\"的订单
	{{else}}
查询客户\"{{.CustomerName}}\"最近的订单
	{{end}}
{{end}}
{{if eq r.Total 0}}
	没有订单可以显示
{{else}}
	{{if gt .Count 0}}
		{{if gt .Count .Total}}
共{{.Total}}个订单，以下为第{{.Prefetched+1}}到第{{.Fetched}}个：
		{{else}}
共{{.Total}}个订单，显示最近的{{.Count}}个订单，以下为第{{.Prefetched+1}}到第{{.Fetched}}个：
		{{end}}
	{{else}}
共{{.Total}}个订单，以下为第{{.Prefetched+1}}到第{{.Fetched}}个：
	{{end}}
{{end}}`

var cusOrdersBodyTpl = `
{{if and (ne . nil) (gt (len .) 0)}}
	{{range .}}
------------------------
订单号：{{.No}}
总金额：{{Sprintf "%.2f" .Amount}}
送货时间：{{.DeliveryTime.Format "2006年01月02日"}}
		{{if ne .Note ""}}
备注：{{.Note}}
		{{end}}
商品:
		{{range .OrderItems}}
{{.ProductName}} {{.Quantity}}{{.Unit}}
		{{end}}
		{{if gt (len .GiftItems) 0}}
赠品:
			{{range .GiftItems}}
{{.ProductName}} {{.Quantity}}{{.Unit}}
			{{end}}
		{{end}}
	{{end}}
	------------------------
	{{if .Done}}
		result = result + "输入\"继续\"，或者\"下一页\"，查看剩下的订单\n"
	{{else}}
		{{if and (gt .Count 0) (le r.Count r.Total)}}
			{{.Count}}个订单已经全部显示完毕！
		{{else}}
			{{.Total}}个订单已经全部显示完毕！
		{{end}}
	{{end}}
{{end}}`

func RenderCustomerOrders(r interface{}, orders *[]database.Order) string {
	tmpl, err := template.New("cusOrdersHeader").Parse(cusOrdersHeaderTpl)
	if err != nil {
		return err.Error()
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, r)
	if err != nil {
		return err.Error()
	}

	result := tpl.String()
	header := strings.Replace(result, "\n\n", "\n", -1)

	tmpl, err = template.New("cusOrdersBody").Parse(cusOrdersBodyTpl)
	if err != nil {
		return err.Error()
	}

	err = tmpl.Execute(&tpl, orders)
	if err != nil {
		return err.Error()
	}

	result = tpl.String()
	body := strings.Replace(result, "\n\n", "\n", -1)

	return header + body
}
