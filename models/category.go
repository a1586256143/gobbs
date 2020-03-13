package models

import (
	"../common"
)
type Category struct{
	Id int 
	Title string
}

func GetCategory() map[int]map[string]interface{} {
	rows := common.ORM.From(&Category{}).Select("id,title")
	return rows
}