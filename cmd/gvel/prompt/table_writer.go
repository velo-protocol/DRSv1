package prompt

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func TableWriter(headers []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.AppendBulk(data)
	table.Render()
}
