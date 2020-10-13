package model

type SystemMessageUser struct {
	AdminUserId     int
	SystemMessageId int
}

// 获取所有系统消息和用户对应表
// [系统消息ID] => [用户ID_1, 用户ID_2, 用户ID_3]
func GetMessageUserMapping() map[int][]int {
	var all []SystemMessageUser
	Db.Select([]string{"system_message_id", "admin_user_id"}).
		Where("system_message_id !=0 ").
		Find(&all)

	//select * from message
	//left join system_messgage
	//
	hash := make(map[int][]int)
	for _, v := range all {
		messageId := v.SystemMessageId
		userId := v.AdminUserId
		_, ok := hash[messageId]
		if ok {
			hash[messageId] = append(hash[messageId], userId)
		} else {
			hash[messageId] = []int{userId}
		}
	}
	return hash
}

// 获取所有系统消息和用户对应表
// [用户ID] => [公告ID_1, 公告ID_2, 公告ID_3]
func GetUserMessageMapping() map[int][]int {
	var all []SystemMessageUser
	Db.Select([]string{"system_message_id", "admin_user_id"}).
		Where("admin_user_id != 0 ").
		Find(&all)
	hash := make(map[int][]int)
	for _, v := range all {
		messageId := v.SystemMessageId
		userId := v.AdminUserId
		_, ok := hash[userId]
		if ok {
			hash[userId] = append(hash[userId], messageId)
		} else {
			hash[userId] = []int{messageId}
		}
	}
	return hash
}
