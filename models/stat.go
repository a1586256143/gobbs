package models

import (
	// "fmt"
	"gobbs/common"
	"strconv"
)

type Stat struct {
	Id    int64 `json:"id"`
	Uid   int64 `json:"uid"`
	Day   int64 `json:"day"`
	Views int64 `json:"views"`
}

// 更新用户的浏览数
func UpdateViews(uid string) int64 {
	day := common.GetToday()
	userid, _ := strconv.ParseInt(uid, 10, 64)
	today, _ := strconv.ParseInt(day, 10, 64)
	condition := map[interface{}]interface{}{
		"uid": userid,
		"day": today,
	}
	find := common.ORM.From(&Stat{}).Where(condition).Find("id")
	if len(find) == 0 {
		lastId, _ := common.ORM.Insert(&Stat{Uid: userid, Day: today, Views: 1})
		if lastId == 0 {
			return 0
		}
	} else {
		sql := "UPDATE stat SET views = views + 1 WHERE uid = ? AND day = ?"
		_, err := common.MysqlDb.Exec(sql, userid, today)
		if err != nil {
			return 0
		}
	}
	sql := "UPDATE gousers SET views = views + 1 WHERE id = ?"
	_, err := common.MysqlDb.Exec(sql, userid)
	if err != nil {
		return 0
	}
	return 1
}

// 获取今天的浏览数
func GetTodayNumber(uid string) interface{} {
	userid, _ := strconv.ParseInt(uid, 10, 64)
	today, _ := strconv.ParseInt(common.GetToday(), 10, 64)
	condition := map[interface{}]interface{}{
		"uid": userid,
		"day": today,
	}
	find := common.ORM.From(&Stat{}).Where(condition).Find("views")
	if len(find) == 0 {
		return "0"
	} else {
		return find["Views"]
	}
}

// 获取本周的浏览数
func GetWeekNumber(uid string) interface{} {
	userid, _ := strconv.ParseInt(uid, 10, 64)
	startDay, endDay := common.GetWeek()
	sql := "SELECT views FROM stat WHERE uid = ? and day between ? and ? LIMIT 1"
	res := common.MysqlDb.QueryRow(sql, userid, startDay, endDay)
	p := Stat{}
	res.Scan(&p.Views)
	return p.Views
}
