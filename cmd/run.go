/*
Copyright © 2020 Yuki.Teraoka <teraoka@beyondjapan.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"encoding/csv"
	"time"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/olekukonko/tablewriter"
)

type Options struct {
	output string
	config string
}

type Config struct {
	Database01 []Database01 `mapstructure:"database01"`
	Database02 []Database02 `mapstructure:"database02"`
}

type Database01 struct {
	Dbms     string `mapstructure:"dbms"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Protocol string `mapstructure:"protocol"`
	Dbname   string `mapstructure:"dbname"`
}

type Database02 struct {
	Dbms     string `mapstructure:"dbms"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Protocol string `mapstructure:"protocol"`
	Dbname   string `mapstructure:"dbname"`
}

type Svs struct {
	Variable_name string
	Value         string
}

var config Config
var o = &Options{}

func NewCmdRun() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "run",
		Short: "A MySQL show variables diff command",
		Run: func(cmd *cobra.Command, args []string) {
			Run(o)
		},
	}

	cmd.Flags().StringVarP(&o.output, "output", "o", "table", "result output option")
	return cmd
}

func Unmarshal(file string) (err error) {
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}
	return nil
}

func Run(o *Options) {
	
	output := o.output

	if output != "csv" && output != "table" {
		fmt.Println("Invalid output option. You can specify csv or table.")
		os.Exit(1)
	}

	if err := Unmarshal(o.config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	svsdb01 := make([][]string, 0)
	svsdb02 := make([][]string, 0)
	
	for _, database01 := range config.Database01 {
		dbms := database01.Dbms
		user := database01.User
		password := database01.Password
		protocol := database01.Protocol
		dbname := database01.Dbname

		db := Connect(dbms,user,password,protocol,dbname)
		svsdb01 = GetRows(db)
		defer db.Close()
	}

	for _, database02 := range config.Database02 {
		dbms := database02.Dbms
		user := database02.User
		password := database02.Password
		protocol := database02.Protocol
		dbname := database02.Dbname

		db := Connect(dbms,user,password,protocol,dbname)
		svsdb02 = GetRows(db)
		defer db.Close()
	}
	
	search_result := Search(svsdb01,svsdb02)

	switch output {
		case "csv":
			WriteCSV(search_result)
		case "table":
			WriteTable(search_result)
	}
}

func Connect(dbms string, user string, password string, protocol string, dbname string) *gorm.DB {
  
	connect := user+":"+password+"@"+protocol+"/"+dbname
	db,err := gorm.Open(dbms, connect)
  
	if err != nil {
		panic(err.Error())
	}
	return db
}

func GetRows(db *gorm.DB) [][]string {
	
	result := make([][]string, 0)
	rows, err := db.Raw("show variables").Rows()

    if err != nil {
        panic(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		svs := &Svs{}
		row := make([]string, 2)
		
		err := rows.Scan(
			&svs.Variable_name,
			&svs.Value,
		)

        if err != nil {
            panic(err.Error())
		}
		
		row = append(row, svs.Variable_name, svs.Value)
		result = append(result, row)
	}
	return result
}

func Search(keywords [][]string, targets [][]string) [][]string {
	
	result := make([][]string, 0)
	exist := 0
	status := "Does not exist."

	for _,  keyword := range keywords {

		search_diff := make([]string, 0)
		keyword_name := keyword[2]
		keyword_value := keyword[3]

		for _, target := range targets {

			target_name := target[2]
			target_value := target[3]

			if keyword_name != target_name {
				continue
			}
			
			exist = 1

			if keyword_value != target_value {
				status = "Different."
				search_diff = append(search_diff, keyword_name)
				search_diff = append(search_diff, keyword_value)
				search_diff = append(search_diff, target_value)
				search_diff = append(search_diff, status)
				result = append(result, search_diff)
			}
		}

		if exist == 0 {
			message := ""
			search_diff = append(search_diff, keyword_name)
			search_diff = append(search_diff, keyword_value)
			search_diff = append(search_diff, message)
			search_diff = append(search_diff, status)
			result = append(result, search_diff)
		}
	}
	return result
}

func WriteTable(search_result [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Parameter", "Value(database01)", "Value(database02)", "Status"})
	table.SetRowLine(true)
	table.AppendBulk(search_result)
	table.Render()
}

func WriteCSV(search_result [][]string) {
	layout := "2006-01-02-15-04-05"
	time := time.Now()
	path := "differ_" + time.Format(layout) + ".csv"

	file, err := os.Create(path)
    if err != nil {
        panic(err.Error())
    }
	defer file.Close()
	
    writer := csv.NewWriter(file)
    if err := writer.WriteAll(search_result); err != nil {
        panic(err.Error())
	}
	
    fmt.Println("CSV File Created.")
}