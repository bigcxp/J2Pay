package model

import "strconv"

type SystemMessageUser struct {
	M   string
	Mid string
	Uid string
}

// 获取所有系统消息和用户对应表
// [系统消息ID] => [用户ID_1, 用户ID_2, 用户ID_3]
func GetMessageUserMapping() map[int][]int {
	var all []SystemMessageUser
	Db.Select([]string{"substring(mid, 9) as mid", "substring(uid, 6) as uid"}).
		Where("mid like ?",  "message:%").
		Find(&all)
	hash := make(map[int][]int)
	for _, v := range all {
		messageId, _ := strconv.Atoi(v.Mid)
		userId, _ := strconv.Atoi(v.Uid)
		_, ok := hash[messageId]
		if ok {
			hash[messageId] = append(hash[messageId], userId)
		} else {
			hash[messageId] = []int{userId}
		}
	}
	return hash
}