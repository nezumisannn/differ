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
var tabledata = make([][]string, 0)

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
	
	search(svsdb01,svsdb02)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Parameter", "Value(database01)", "Value(database02)", "Status"})
	table.SetRowLine(true)
	table.AppendBulk(tabledata)
	table.Render()
}

func Connect(dbms string, user string, password string, protocol string, dbname string) *gorm.DB {
  
	connect := user+":"+password+"@"+protocol+"/"+dbname
	db,err := gorm.Open(dbms, connect)
  
	if err != nil {
		panic(err.Error())
	}
	return db
}

func GetRows(db *gorm.DB) (result [][]string) {
	slicerows := make([][]string, 0)
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
		slicerows = append(slicerows, row)
	}

	return slicerows
}

func search(slice [][]string, targets [][]string) {
	for _,  rows := range slice {
		hit := 0
		status := "Does not exist."
		name01 := rows[2]
		value01 := rows[3]
		for _, target := range targets {
			name02 := target[2]
			value02 := target[3]

			if name01 != name02 {
				continue
			}

			hit = 1
			
			if value01 != value02 {
				status = "Different."
				diff := make([]string, 0)
				diff = append(diff, name01)
				diff = append(diff, value01)
				diff = append(diff, value02)
				diff = append(diff, status)
				tabledata = append(tabledata, diff)
			}
		}
		if hit == 0 {
			message := "NONE"
			diff := make([]string, 0)
			diff = append(diff, name01)
			diff = append(diff, value01)
			diff = append(diff, message)
			diff = append(diff, status)
			tabledata = append(tabledata, diff)
		}
	}
}