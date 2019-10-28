package console

import (
	"bytes"
	"github.com/olekukonko/tablewriter"
)

func WriteTable(headers []string, data [][]string) {
	var buffer bytes.Buffer
	table := tablewriter.NewWriter(&buffer)
	table.SetHeader(headers)
	table.AppendBulk(data)
	table.Render()

	TableLogger.Println(buffer.String())
}
