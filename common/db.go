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
	Select(field string)  map[int]map[string]interface{}
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
func(m *Model)From(table interface{}) *Model  {
	m.Table = table
	m.TableName = strings.ToLower(reflect.TypeOf(table).Elem().Name())
	return m
}

// 解析Where条件
func(m *Model) ParseWhere(condition map[interface{}]interface{}) string {
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
	m.bind = value
	return buffer.String()

}

// where条件
func (m *Model) Where(condition map[interface{}]interface{}) *Model  {
	m.where = m.ParseWhere(condition)
	return m
}


// Model实现DB接口，查询方法
func(m *Model) Select(field string) map[int]map[string]interface{}{
	execSql := "SELECT " + field + " FROM " + m.TableName
	if m.where != "" {
		execSql += " WHERE " + m.where
	}
	if m.order != "" {
		execSql += " ORDER BY " + m.order
	}
	m.sql = execSql
	return m.execute()
}

// Model实现DB接口，查询方法
func(m *Model) Find(field string) map[string]interface{}{
	execSql := "SELECT " + field + " FROM " + m.TableName
	if m.where != "" {
		execSql += " WHERE " + m.where
	}
	m.sql = execSql
	data := m.execute()
	return data[0]
}

// 统计条数
func (m *Model) Count(field string) string {
	data := m.Find("count(" + field + ") as count_field")
	return data["CountField"].(string)
}

/**
 * 插入操作
 * common.ORM.Insert(Article{Id:1,Title:111})
 * 返回插入的ID
 */
func (m *Model) Insert(table interface{}) (int64 , error) {
	m.From(table)
	joinFields , joinValues , values := GetModelInfo(table)
	s := "INSERT INTO " + m.TableName + "(" + joinFields + ") VALUES (" + joinValues + ")"
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
func(m *Model) Order(order string) *Model {
	m.order = order
	split := strings.Split(order , " ")
	if len(split) == 1 {
		split = append(split, "DESC")
	}
	m.order = strings.Join(split , " ")
	return m
}

// 执行SQL语句
func (m *Model) execute() map[int]map[string]interface{}{
	rows, err := MysqlDb.Query(m.sql , m.bind... )
	defer func() {
		m.clear()
		rows.Close()
	}()
	if rows != nil && err == nil {
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

			for j , value := range values { //注意：此处用values
				key := AutoTuoFeng(cols[j])
				row[key] = string(value)
			}
			results[i] = row
			i ++
		}
		return results

	}else {
		fmt.Println("execute error : " , err)
	}
	return nil
}

// 获取SQL语句
func (m *Model) GetSql() string {
	return m.sql
}

// 情理
func (m *Model) clear()  {
	m.where , m.order = "" , ""
	m.bind = make([]interface{} , 0)
}