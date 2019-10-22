package console

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func WriteTable(headers []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(data)
	table.Render()
}
