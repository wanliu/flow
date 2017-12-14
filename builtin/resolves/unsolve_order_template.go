package resolves

import (
	"bytes"
	"strings"
	"text/template"
)

var body string = `{{range .Products.Products }}
{{.Product}} {{.Quantity}}{{.Unit}}
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

func RenderOrderBody(order OrderResolve) string {
	tmpl, err := template.New("orderBody").Parse(body)
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
