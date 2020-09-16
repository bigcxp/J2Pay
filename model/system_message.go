package model

import (
	"j2pay-server/model/response"
	"strconv"
	"time"
)

type SystemMessage struct {
	Id        int
	Title     string    `gorm:"default:'';comment:'标题';"`
	BeginTime time.Time `gorm:"comment:'开始时间';type:timestamp;";json:"beginTime"`
	EndTime   time.Time `gorm:"comment:'结束时间';type:timestamp;";json:"endTime"`
	Status    int       `gorm:"default:1;comment:'是否作废 0：否，1：是';"`
}

// 获取所有系统消息
func (s *SystemMessage) GetAll(page, pageSize int, where ...interface{}) (response.SystemMessagePage, error) {
	all := response.SystemMessagePage{
		Total:       s.GetCount(where...),
		PerPage:     pageSize,
		CurrentPage: page,
		Data:        []response.SystemMessageList{},
	}
	offset := GetOffset(page, pageSize)
	err := Db.Table("system_message").
		Limit(pageSize).
		Offset(offset).
		Find(&all.Data, where...).Error
	if err != nil {
		return response.SystemMessagePage{}, err
	}
	return all, err
}

// 根据ID获取公告详情
func (s *SystemMessage) Detail(id ...int) (res response.SystemMessageList, err error) {
	searchId := s.Id
	if len(id) > 0 {
		searchId = id[0]
	}
	err = Db.Table("system_message").
		Where("id = ?", searchId).
		First(&res).
		Error
	return
}

//新增公告
func (s *SystemMessage) Create(users []int) error {
	tx := Db.Begin()
	if err := tx.Create(s).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&SystemMessage{
		Title:     s.Title,
		BeginTime: s.BeginTime,
		EndTime:   s.EndTime,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, v := range users {
		err := tx.Create(&SystemMessageUser{
			M:   s.Title,
			Mid: "message:" + strconv.Itoa(s.Id),
			Uid: "user:" + strconv.Itoa(v),
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// 编辑系统公告
func (s *SystemMessage) Edit(users []int) error {
	tx := Db.Begin()
	updateInfo := map[string]interface{}{
		"title":      s.Title,
		"begin_time": s.BeginTime,
		"end_time":   s.EndTime,
	}
	if s.Title != "" {
		updateInfo["title"] = s.Title
	}

	if err := tx.Model(&SystemMessage{Id: s.Id}).
		Updates(updateInfo).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(SystemMessageUser{}, "mid = ?", "message:"+strconv.Itoa(s.Id)).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, v := range users {
		err := tx.Create(&SystemMessageUser{
			M:   s.Title,
			Mid: "message:" + strconv.Itoa(s.Id),
			Uid: "user:" + strconv.Itoa(v),
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// 删除系统公告
func (s *SystemMessage) Del() error {
	tx := Db.Begin()
	if err := tx.Delete(s, "id = ?", s.Id).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(SystemMessageUser{}, "m = ? and mid = ?", s.Title, "mid:"+strconv.Itoa(s.Id)).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 获取所有系统公告数量
func (s *SystemMessage) GetCount(where ...interface{}) (count int) {
	if len(where) == 0 {
		Db.Model(&s).Count(&count)
		return
	}
	Db.Model(&s).Where(where[0], where[1:]...).Count(&count)
	return
}

func GetAllMessage() (mapping map[int]response.AdminSystemMessage) {
	var systemMessages []response.AdminSystemMessage
	mapping = make(map[int]response.AdminSystemMessage)
	Db.Table("system_message").Select("id,title,begin_time,end_time").Order("id desc").Find(&systemMessages)
	for _, systemMessage := range systemMessages {
		mapping[systemMessage.Id] = systemMessage
	}
	return
}
