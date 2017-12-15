package resolves

import (
	"bytes"
	"strings"
	"text/template"
)

var header string = `订单正在处理, 已经添加{{ CnNum (len .Products.Products)}}种产品 {{if gt (len .Gifts.Products) 0 }},{{CnNum (len .Gifts.Products)}}种赠品
{{end}}`

var body string = `{{range .Products.Products }}{{.Product}} {{.Quantity}}{{.Unit}}
{{end}}
{{if gt (len .Gifts.Products) 0}}
申请的赠品:
{{range .Gifts.Products}}
{{.Product}} {{.Quantity}}{{.Unit}}
{{end}}
{{end}}
时间:{{.Time.Format "2006年01月02日"}}
{{if ne .Note ""}}
备注：{{.Note}}
{{end}}`

func RenderOrderBody(r OrderResolve) string {
	tmpl, err := template.New("orderBody").Parse(body)
	if err != nil {
		return err.Error()
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, r)
	if err != nil {
		return err.Error()
	}

	result := tpl.String()
	result = strings.Replace(result, "\n\n", "\n", -1)
	result = strings.Replace(result, "\n\n", "\n", -1)
	return result
}

func RenderOrderHeader(r OrderResolve) string {
	tmpl, err := template.New("orderHeader").Funcs(template.FuncMap{
		"CnNum": CnNum,
	}).Parse(header)

	if err != nil {
		return err.Error()
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, r)
	if err != nil {
		return err.Error()
	}

	result := tpl.String()
	result = strings.Replace(result, "\n\n", "\n", -1)
	return result
}
