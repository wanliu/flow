package context

type Tabler interface {
	Render()
}

type Table struct {
	Headers []string
	Rows    [][]string
	Footers []string
}

func NewTable() *Table {
	return &Table{
		Headers: make([]string, 0),
		Rows:    make([][]string, 0),
		Footers: make([]string, 0),
	}
}

func (t *Table) SetHeader(header []string) {
	t.Headers = header
}

func (t *Table) SetFooter(footer []string) {
	t.Footers = footer
}

func (t *Table) AppendBulk(rows [][]string) {
	for _, row := range rows {
		t.Append(row)
	}
}

func (t *Table) Append(row []string) {
	t.Rows = append(t.Rows, row)
}
