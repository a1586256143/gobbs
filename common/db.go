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

func (self *Model) Insert(table interface{}) (int64 , error) {
	self.From(table)
	t := reflect.TypeOf(table).Elem()
	v := reflect.ValueOf(table).Elem()
	var fields []string
	var tmpValues []string
	values := make([]interface{} , 0)
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		field := t.Field(i)
		fieldName := field.Tag.Get("json")
		if fieldName != "" {
			fieldValue := v.Field(i)
			fields = append(fields, fieldName)
			var convValue interface{}
			switch fieldValue.Type().Kind().String(){
			case "string" :
				convValue = fieldValue.String()
			case "int" , "int64" :
				convValue = fieldValue.Int()
			case "float" , "float32" , "float64" :
				convValue = fieldValue.Float()
			}
			values = append(values , &convValue)
			tmpValues = append(tmpValues , "?")
		}
	}
	joinFields := strings.Join(fields , ",")
	joinValues := strings.Join(tmpValues , ",")
	s := "INSERT INTO " + self.TableName + "(" + joinFields + ") VALUES (" + joinValues + ")"
	rst , err := MysqlDb.Exec(s, values ...)
	if err == nil {
		rowsAffected , _ := rst.RowsAffected()
		lastInsertID,_ := rst.LastInsertId()
		if rowsAffected >= 1 {
			return lastInsertID , nil
		}
	}
	return 0 , nil
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

// 执行SQL语句
func (self *Model) execute() map[int]map[string]interface{}{
	rows, err := MysqlDb.Query(self.sql , self.bind... )
	fmt.Println("execute" , err)
	defer func() {
		self.clear()
		rows.Close()
	}()
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

// 获取SQL语句
func (self *Model) GetSql() string {
	return self.sql
}

// 情理
func (self *Model) clear()  {
	self.where , self.order = "" , ""
	self.bind = make([]interface{} , 0)
}