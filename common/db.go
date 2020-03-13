package common

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

var ORM Model

// DB接口
type DB interface {
	Select(condition map[interface{}]interface{} , field string) string
}

// Model类
type Model struct {
	Table interface{}
	sql   string
	order string
	bind  []interface{}
	TableName string
	where string
}

// 设置表名
func(self *Model)From(table interface{}) *Model  {
	self.Table = table
	self.TableName = strings.ToLower(reflect.TypeOf(table).Elem().Name())
	return self
}

// 解析Where条件
func(self *Model) ParseWhere(condition map[interface{}]interface{}) string {
	var buffer bytes.Buffer
	buffer.Grow(len(condition))
	var i int
	var value []interface{}
	for k , v := range condition{
		buffer.WriteString(k.(string))
		buffer.WriteString("=?")
		value = append(value , v)
		//switch v.(type){
		//case string :
		//	buffer.WriteString(v.(string))
		//case float64 :
		//	buffer.WriteString(strconv.FormatFloat(v.(float64) ,'f' , 2 , 64))
		//case int :
		//	buffer.WriteString(strconv.Itoa(v.(int)))
		//default :
		//	fmt.Println("unkown" , k , v)
		//}
		if i != len(condition) - 1 {
			buffer.WriteString(" AND ")
		}
		i ++
	}
	self.bind = value
	return buffer.String()

}

func (self *Model) Where(condition map[interface{}]interface{}) *Model  {
	self.where = self.ParseWhere(condition)
	return self
}


// Model实现DB接口
func(self *Model) Select(field string) map[int]map[string]interface{}{
	execSql := "SELECT " + field + " FROM " + self.TableName
	if self.where != "" {
		execSql += " WHERE " + self.where
	}
	if self.order != "" {
		execSql += " ORDER BY " + self.order
	}
	self.sql = execSql
	return self.execute()
}

// 排序
func(self *Model) Order(order string) *Model {
	self.order = order
	split := strings.Split(order , " ")
	if len(split) == 1 {
		split = append(split, "DESC")
	}
	self.order = strings.Join(split , " ")
	return self
}

func (self *Model) execute() map[int]map[string]interface{}{
	rows, _ := MysqlDb.Query(self.sql , self.bind... )
	if rows != nil {
		cols , _ := rows.Columns()
		values := make([]sql.RawBytes, len(cols))
		scans := make([]interface{}, len(cols))

		// 将接口指针转换为指针类型的接口
		for i := range values{
			scans[i] = &values[i]
		}
		results := make(map[int]map[string]interface{})
		i := 0
		for rows.Next() {
			if err := rows.Scan(scans...); err != nil{
				fmt.Println("error , " , err)
			}
			row := make(map[string]interface{})

			for j , value := range values { //注意：此处用vals
				key := AuthTuoFeng(cols[j])
				row[key] = string(value)
			}
			results[i] = row
			i ++
		}
		return results

	}
	return nil
}