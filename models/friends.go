package models

import (
	"fmt"
	"gobbs/common"
	"strconv"
)

// 定义好友的struct
type Friends struct {
	Id         int64  `json:"id"`
	Uid        int64  `json:"uid"`
	PullUid    int64  `json:"pull_uid"`
	Remark     string `json:"remark"`
	Status     int    `json:"status"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
	Count      int64  `json:""`
}

// 检验用户是否是好友关系了
func IsFriends(uid, pull_uid int64) int64 {
	sql := "SELECT id FROM friends where ((uid = ? and pull_uid = ?) OR (uid = ? and pull_uid = ?)) and status = 1 LIMIT 1"
	res := common.MysqlDb.QueryRow(sql, uid, pull_uid, pull_uid, uid)
	fmt.Println("uid", uid, "pull_uid", pull_uid)
	p := Friends{}
	res.Scan(&p.Id)
	if p.Id > 0 {
		return 1
	}
	return 0
}

// 是否申请过了
func IsExists(uid, pull_uid int64) int64 {
	sql := "SELECT id FROM friends where ((uid = ? and pull_uid = ?) OR (uid = ? and pull_uid = ?)) and status != 2 LIMIT 1"
	res := common.MysqlDb.QueryRow(sql, uid, pull_uid, pull_uid, uid)
	p := Friends{}
	res.Scan(&p.Id)
	return p.Id
}

// 获取申请添加好友的数量
func GetFriendsApplyNumber(pull_uid int64) int64 {
	sql := "SELECT count(id) FROM friends WHERE pull_uid = ? AND status = 0"
	res := common.MysqlDb.QueryRow(sql, pull_uid)
	p := Friends{}
	res.Scan(&p.Count)
	return p.Count
}

// 获取我的好友数量
func GetFriendsNumber(pull_uid int64) int64 {
	sql := "SELECT count(id) FROM friends WHERE (pull_uid = ? OR uid = ?) AND status = 1"
	res := common.MysqlDb.QueryRow(sql, pull_uid, pull_uid)
	p := Friends{}
	res.Scan(&p.Count)
	return p.Count
}

// 发起一个 添加朋友的请求
func AddFriends(uid, pull_uid int64, remark string) int64 {
	lastId, err := common.ORM.Insert(&Friends{Uid: uid, PullUid: pull_uid, Remark: remark, Status: 0, CreateTime: common.GetUnix()})
	if lastId > 0 {
		return lastId
	} else {
		fmt.Println(&Friends{Uid: uid, PullUid: pull_uid, Remark: remark, Status: 0, CreateTime: common.GetUnix()})
		fmt.Println("add Friends error", err)
		return 0
	}
}

// 获取我的请求列表
func GetMyApplyList(condition map[interface{}]interface{}) map[int]map[string]interface{} {
	list := common.ORM.From(&Friends{}).Order("create_time").Where(condition).Select("*")
	if len(list) > 0 {
		for _, v := range list {
			if v["CreateTime"] != "0" {
				times, _ := strconv.ParseInt(v["CreateTime"].(string), 10, 64)
				v["FormatTime"] = common.DateFormat(times)
			}
			// 查询NickName
			userCondition := map[interface{}]interface{}{
				"id": v["Uid"],
			}
			info := common.ORM.From(&Gousers{}).Where(userCondition).Find("name,avatar")
			v["NickName"] = info["Name"]
			v["Avatar"] = info["Avatar"]
		}
	}
	return list
}

// 同意或拒绝请求
func Accept(id, status int64) bool {
	sql := "UPDATE friends SET status = ?,update_time = ? WHERE id = ?"
	res, _ := common.MysqlDb.Exec(sql, status, common.GetUnix(), id)
	affect, _ := res.RowsAffected()
	if affect > 0 {
		return true
	}
	return false
}

// 获取我的好友
func GetMyFriendsList(uid int64) map[int]map[string]interface{} {
	sql := fmt.Sprintf("SELECT id,uid,pull_uid FROM friends where (uid = %d OR pull_uid = %d) and status = 1", uid, uid)
	list := common.ORM.SetSql(sql).QuerySql()
	for _, value := range list {
		var info map[string]interface{}
		listUid, _ := strconv.ParseInt(value["Uid"].(string), 10, 64)
		listPullUid, _ := strconv.ParseInt(value["PullUid"].(string), 10, 64)
		// 接收人是自己，则读取uid的信息
		if listPullUid == uid {
			info = GetUserBase(listUid)
			value["SendUid"] = listUid
		} else if listUid == uid {
			info = GetUserBase(listPullUid)
			value["SendUid"] = listPullUid
		}
		value["Name"] = info["Name"]
		value["Avatar"] = info["Avatar"]
	}
	return list
}

// 删除好友
func DeleteFriend(id string, uid int64) bool {
	sql := "UPDATE friends set status = 88 WHERE id = ? AND (uid = ? OR pull_uid = ?) and status = 1"
	res, _ := common.MysqlDb.Exec(sql, id, uid, uid)
	affect, _ := res.RowsAffected()
	if affect > 0 {
		return true
	}
	return false
}
