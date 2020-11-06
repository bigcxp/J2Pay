package model

//app 锁

type AppLock struct {
	ID         int64  `db:"id" json:"id"`
	K          string `gorm:"unique;comment:'上锁键值'";json:"k"`               // 上锁键值
	V          int64  `gorm:"comment:'是否锁定 1:锁定，0：解锁'";json:"v"`        // 是否锁定(1:锁定，0：释放)
	CreateTime int64  `gorm:"default:0;comment:'创建时间戳'";json:"create_time"` // 上锁时间戳
}

//根据k查询
func SQLGetAppLockColByK(k string) (*AppLock, error) {
	var row *AppLock
	err := Getdb().Model(&AppLock{}).First(&row, k).Error
	return row, err
}

// 创建锁
func SQLCreateAppLockUpdate(row *AppLock) (int64, error) {
	tx := Getdb().Begin()
	if err := Getdb().Model(&AppLock{}).Create(row).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return 0, nil

}

//  更新锁
func SQLUpdateTAppLockByK(row *AppLock) (int64, error) {
	// 执行事务处理
	tx := Getdb().Begin()
	// 1.更新角色表
	if err := tx.Model(&AppLock{K: row.K}).
		Updates(map[string]interface{}{
			"v":           row.V,
			"create_time": row.CreateTime,
		}).Error;
		err != nil {
		tx.Rollback()
		return 0, err
	}
	return 0, nil
}
