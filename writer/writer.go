package writer

import (
	"fmt"
	"os"
	"encoding/csv"
	"time"

	"github.com/olekukonko/tablewriter"
)

func Table(SearchResult [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Parameter", "Value(database01)", "Value(database02)", "Status"})
	table.SetRowLine(true)
	table.AppendBulk(SearchResult)
	table.Render()
}

func CSV(SearchResult [][]string) {
	layout := "2006-01-02-15-04-05"
	time := time.Now()
	path := "differ_" + time.Format(layout) + ".csv"

	file, err := os.Create(path)
    if err != nil {
        panic(err.Error())
    }
	defer file.Close()
	
    writer := csv.NewWriter(file)
    if err := writer.WriteAll(SearchResult); err != nil {
        panic(err.Error())
	}
	
    fmt.Println("CSV File Created.")
}