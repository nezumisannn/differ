package connector

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Variables struct {
	Variable_name string
	Value         string
}

func Mysql(dbms string, user string, password string, protocol string, dbname string) *gorm.DB {
  
	connect := user+":"+password+"@"+protocol+"/"+dbname
	db,err := gorm.Open(dbms, connect)
  
	if err != nil {
		panic(err.Error())
	}
	return db
}

func ShowVariables(db *gorm.DB) [][]string {
	
	result := make([][]string, 0)
	rows, err := db.Raw("show variables").Rows()

    if err != nil {
        panic(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		svs := &Variables{}
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