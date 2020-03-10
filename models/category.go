package models
import (
	"../common"
)
type Category struct{
	Id int 
	Title string
}

func GetCategory() (categorys []Category){
	rows , _ := common.MysqlDb.Query("select id,title from category")
	if rows != nil {
		for rows.Next() {
			category := Category{}
			rows.Scan(&category.Id , &category.Title)
			categorys = append(categorys , category)
		}
		return
	}
	return
}