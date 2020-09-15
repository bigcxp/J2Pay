package model

import (
	"j2pay-server/model/response"
	"time"
)

type SystemMessage struct {
	Id           int
	Title        string    `gorm:"default:'';comment:'标题';"`
	OrganizeName string    `gorm:"default:'';comment:'组织名称';"`
	BeginTime    time.Time `gorm:"comment:'开始时间';type:timestamp;";json:"beginTime"`
	EndTime      time.Time `gorm:"comment:'结束时间';type:timestamp;";json:"endTime"`
	Status       int       `gorm:"default:0;comment:'是否作废 0：否，1：是';"`
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

//新增公告
func (s *SystemMessage) Create() error {
	tx := Db.Begin()
	if err := tx.Create(s).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&SystemMessage{
		Title:        s.Title,
		BeginTime:    s.BeginTime,
		EndTime:      s.EndTime,
		OrganizeName: s.OrganizeName,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 编辑系统公告
func (s *SystemMessage) Edit() error {
	tx := Db.Begin()
	updateInfo := map[string]interface{}{
		"title":      s.Title,
		"begin_time": s.BeginTime,
		"end_time":   s.EndTime,
	}
	if s.Title != "" {
		updateInfo["title"] = s.Title
		updateInfo["begin_time"] = s.BeginTime
		updateInfo["end_time"] = s.EndTime
	}
	if err := tx.Model(&SystemMessage{Id: s.Id}).
		Updates(updateInfo).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&SystemMessage{
		Title:        s.Title,
		BeginTime:    s.BeginTime,
		EndTime:      s.EndTime,
		OrganizeName: s.OrganizeName,
	}).Error; err != nil {
		tx.Rollback()
		return err
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
	tx.Commit()
	return nil
}

// 获取所有后台用户数量
func (s *SystemMessage) GetCount(where ...interface{}) (count int) {
	if len(where) == 0 {
		Db.Model(&s).Count(&count)
		return
	}
	Db.Model(&s).Where(where[0], where[1:]...).Count(&count)
	return
}
