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
	"github.com/nezumisannn/differ/config"
	"github.com/nezumisannn/differ/connector"
	"github.com/nezumisannn/differ/writer"
)

type options struct {
	output string
}

var o = &options{}

func NewCmdRun() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "run",
		Short: "A MySQL show variables diff command",
		Run: func(cmd *cobra.Command, args []string) {
			cfgFile := cfgFile
			Run(cfgFile)
		},
	}

	cmd.Flags().StringVarP(&o.output, "output", "o", "table", "result output option")
	return cmd
}

func Run(cfgFile string) {
	
	Varidete()

	if err := config.LoadFile(cfgFile); err != nil {
		fmt.Println(err)
	}

	Db01Dbms := "mysql"
	Db01User := config.Cfg.Differ.Database01.User
	Db01Password := config.Cfg.Differ.Database01.Password
	Db01Protocol := config.Cfg.Differ.Database01.Protocol
	Db01Dbname := "mysql"

	Db01Connect := connector.Mysql(Db01Dbms,Db01User,Db01Password,Db01Protocol,Db01Dbname)
	SvsDb01 := connector.ShowVariables(Db01Connect)

	Db02Dbms := "mysql"
	Db02User := config.Cfg.Differ.Database02.User
	Db02Password := config.Cfg.Differ.Database02.Password
	Db02Protocol := config.Cfg.Differ.Database02.Protocol
	Db02Dbname := "mysql"

	Db02Connect := connector.Mysql(Db02Dbms,Db02User,Db02Password,Db02Protocol,Db02Dbname)
	SvsDb02 := connector.ShowVariables(Db02Connect)

	SearchResult := Search(SvsDb01,SvsDb02)

	switch o.output {
		case "csv":
			writer.CSV(SearchResult)
		case "table":
			writer.Table(SearchResult)
	}
}

func Search(keywords [][]string, targets [][]string) [][]string {
	
	result := make([][]string, 0)
	exist := 0
	status := "Does not exist."

	for _,  keyword := range keywords {

		SearchDiff := make([]string, 0)
		KeywordName := keyword[2]
		KeywordValue := keyword[3]

		for _, target := range targets {

			TargetName := target[2]
			TargetValue := target[3]

			if KeywordName != TargetName {
				continue
			}
			
			exist = 1

			if KeywordValue != TargetValue {
				status = "Different."
				SearchDiff = append(SearchDiff, KeywordName)
				SearchDiff = append(SearchDiff, KeywordValue)
				SearchDiff = append(SearchDiff, TargetValue)
				SearchDiff = append(SearchDiff, status)
				result = append(result, SearchDiff)
			}
		}

		if exist == 0 {
			message := "Undefined"
			SearchDiff = append(SearchDiff, KeywordName)
			SearchDiff = append(SearchDiff, KeywordValue)
			SearchDiff = append(SearchDiff, message)
			SearchDiff = append(SearchDiff, status)
			result = append(result, SearchDiff)
		}
	}
	return result
}

func Varidete() {
	if o.output != "csv" && o.output != "table" {
		fmt.Println("Invalid output option. You can specify csv or table.")
		os.Exit(1)
	}
}